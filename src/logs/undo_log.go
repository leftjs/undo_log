package logs

import (
	"config"
	ds "datastructure"
	"db"
	"file"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"util"
)

/**
维持一个时间窗口

*/
/**

按如下约定记录日志：

1. 事务开始时，记录START T
2. 事务修改时，记录（T，x，v），说明事务T操作对象x，x的值为v
3. 事务结束时，记录COMMIT T
*/

/**
使用undo log时事务执行顺序

1. 记录START T
2. 记录需要修改的记录的旧值（要求持久化）
3. 根据事务的需要更新数据库（要求持久化）
4. 记录COMMIT T

使用undo log进行宕机回滚
*/

// undo log path

const (
	REQUEST_START RequestType = iota
	REQUEST_PUT
	REQUEST_COMMIT
	REQUEST_UNDO
)

type RequestType int // transaction request type

type Log struct {
	mu sync.RWMutex

	Logfile  string // 当前 log 写入的文件
	UndoLogs map[int]*ds.LinkedList
	Config   *config.Config
}

// log object in memory
type Undo struct {
	ID     int
	Type   RequestType
	States []State
}

// preserved state in log object
type State struct {
	UserId int
	Cash   int
}

func generateUndoFromString(str string) *Undo {
	results := regexp.MustCompile(`^<(COMMIT|UNDO|START)?\s?T(\d+)(,(\d+),(\d+),(\d+),(\d+))?>$`).FindStringSubmatch(str)
	if len(results) != 8 {
		return nil
	}
	id, _ := strconv.Atoi(results[2])

	var undo *Undo
	if results[1] != "" {
		// commit or undo or state request
		undo = &Undo{ID: id}
		if results[1] == "COMMIT" {
			undo.Type = REQUEST_COMMIT
		} else if results[1] == "UNDO" {
			undo.Type = REQUEST_UNDO
		} else {
			undo.Type = REQUEST_START
		}
	} else {
		fromId, fromCash, toId, toCash := extractTransferFromPutLog(str)
		undo = &Undo{
			id,
			REQUEST_PUT,
			[]State{{fromId, fromCash}, {toId, toCash}},
		}
	}

	return undo
}

func NewLog() *Log {
	l := &Log{}
	l.Config = config.NewConfig()
	l.initializeLastLogfile()
	l.UndoLogs = make(map[int]*ds.LinkedList)
	l.buildUndoLogs()
	return l
}

/**
从 log_path 的 undo logs 中 scan 出最后一个 logfile
若不存在则新建
*/
func (l *Log) initializeLastLogfile() {
	if existed, _ := file.CheckExisted(l.Config.LogPath); existed == false {
		file.MakeDir(l.Config.LogPath)
	}
	files, err := ioutil.ReadDir(l.Config.LogPath)
	util.Check(err)
	if len(files) == 0 {
		name := file.CreateFile(path.Join(l.Config.LogPath, fmt.Sprintf("%d.log", time.Now().Unix())))
		if &name != nil {
			l.Logfile = name
		}
		return
	}

	sort.SliceStable(files, func(i, j int) bool {
		return extractUnixFromFileName(files[i].Name()) < extractUnixFromFileName(files[j].Name())
	})
	last := files[len(files)-1]
	l.Logfile = last.Name()
}

/**
从日志文件构建内存中的undolog
*/
func (l *Log) buildUndoLogs() {
	content := strings.Trim(string(file.ReadFile(path.Join(l.Config.LogPath, l.Logfile))), "\n")
	if len(content) == 0 {
		return
	}

	logs := strings.Split(content, "\n")

	for _, entry := range logs {
		undo := generateUndoFromString(entry)
		l.appendLogToMemory(undo.ID, undo)
	}
}

/**
获取下一个 transaction id
*/
func (l *Log) GetNextTransactionId() int {
	max := 0
	for id, _ := range l.UndoLogs {
		if id > max {
			max = id
		}
	}
	return max + 1
}

/**
从文件名字中提取 timestamp
*/
func extractUnixFromFileName(filename string) int {
	i, _ := strconv.Atoi(strings.Split(filename, ".")[0])
	return i
}

func CheckDone(s string, tId int) error {

	results := regexp.MustCompile(`^<(COMMIT|UNDO) T(\d+)>$`).FindStringSubmatch(s)
	if len(results) == 3 {
		id, _ := strconv.Atoi(results[2])
		if id == tId {
			// 该事务已经结束
			return errors.New("transaction has been committed or undid")
		}
	}
	return nil
}

/**
undo 必须按照顺序，为保证能够恢复到初始状态，不推荐异步
*/
func (l *Log) Undo(tId int) {

	userDB := db.NewUserDB()

	undo := l.UndoLogs[tId]
	if undo == nil || !(undo.Tail.Data.(*Undo).Type == REQUEST_PUT) {
		return
	}

	cur := undo.Tail

RollBack:
	for {
		for _, state := range cur.Data.(*Undo).States {
			userDB.UpdateCash(state.UserId, state.Cash)
		}
		if cur.HasPrev() && cur.Prev.Data.(*Undo).Type != REQUEST_START {
			cur = cur.Prev
		} else {
			break RollBack
		}
	}

	l.writeUndo(tId)
}

func extractTransferFromPutLog(put string) (int, int, int, int) {
	results := regexp.MustCompile(`^<T(\d+)\D(\d+)\D(\d+)\D(\d+)\D(\d+)>$`).FindStringSubmatch(put)
	fromId, _ := strconv.Atoi(results[2])
	fromCash, _ := strconv.Atoi(results[3])
	toId, _ := strconv.Atoi(results[4])
	toCash, _ := strconv.Atoi(results[5])
	return fromId, fromCash, toId, toCash
}

/**
触发一次 undo log的 gc 请求
*/
//func (l *Log) GCUndoLog() (bool, error) {
//
//}

/**
append log to memory
*/
func (l *Log) appendLogToMemory(tId int, undo *Undo) {
	log := l.UndoLogs[tId]
	if log == nil {
		l.UndoLogs[tId] = ds.NewLinkedList()
		log = l.UndoLogs[tId]
	}
	log.Append(undo)
}

// some log file write function
func (l *Log) WriteStart(tId int) {
	s := fmt.Sprintf("<START T%d>", tId)
	l.appendLogToMemory(tId, generateUndoFromString(s))
	file.AppendToFile(path.Join(l.Config.LogPath, l.Logfile), s)
}
func (l *Log) WritePut(tId, fromId, fromOriginalCash, toId, toOriginalCash int) {
	s := fmt.Sprintf("<T%d,%d,%d,%d,%d>", tId, fromId, fromOriginalCash, toId, toOriginalCash)
	l.appendLogToMemory(tId, generateUndoFromString(s))
	file.AppendToFile(path.Join(l.Config.LogPath, l.Logfile), s)
}
func (l *Log) WriteCommit(tId int) {
	s := fmt.Sprintf("<COMMIT T%d>", tId)
	l.appendLogToMemory(tId, generateUndoFromString(s))
	file.AppendToFile(path.Join(l.Config.LogPath, l.Logfile), s)
}
func (l *Log) writeUndo(tId int) {
	s := fmt.Sprintf("<UNDO T%d>", tId)
	l.appendLogToMemory(tId, generateUndoFromString(s))
	file.AppendToFile(path.Join(l.Config.LogPath, l.Logfile), s)
}
