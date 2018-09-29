package logs

import (
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
const LOG_PATH = "../../log/"

type Log struct {
	mu sync.RWMutex

	Logfile string // 当前 log 写入的文件

}

func NewLog() *Log {
	l := &Log{}
	initializeLastLogfile(l)
	return l
}

/**
从 log_path 的 undo logs 中 scan 出最后一个 logfile
若不存在则新建
*/
func initializeLastLogfile(l *Log) {

	if existed, _ := file.CheckExisted(LOG_PATH); existed == false {
		file.MakeDir(LOG_PATH)
	}
	files, err := ioutil.ReadDir(LOG_PATH)
	util.Check(err)
	if len(files) == 0 {
		name := file.CreateFile(path.Join(LOG_PATH, fmt.Sprintf("%d.log", time.Now().Unix())))
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
获取下一个 transaction id
*/
func (l *Log) GetNextTransactionId() int {

	data := file.ReadFile(path.Join(LOG_PATH, l.Logfile))
	ls := strings.Split(strings.Trim(string(data), "\n"), "\n")
	if len(ls) == 1 && ls[0] == strings.Trim(string(data), "\n") {
		// 空
		return 1
	}

	startIdC := make(chan int)
	var wg sync.WaitGroup

	for i := len(ls) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			results := regexp.MustCompile(`^<START T(\d+)>$`).FindStringSubmatch(ls[ii])
			if len(results) > 1 {
				id, _ := strconv.Atoi(results[1])
				startIdC <- id
			}
		}(i)
	}

	// max transactionId
	max := 0

	go func() {
		wg.Wait()
		close(startIdC)
	}()

	for id := range startIdC {
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
func (l *Log) Undo(tId int) (bool, error) {

	data := file.ReadFile(path.Join(LOG_PATH, l.Logfile))

	lStr := strings.Trim(string(data), "\n")
	ls := strings.Split(lStr, "\n")

	if len(ls) == 1 && ls[0] == lStr {
		// undo log is null
		return false, nil
	}

	needUndo, err := l.CheckUndoLog(tId, ls)
	if err != nil {
		// has been undid or committed
		return false, err
	}
	if !needUndo {
		// undo log is null
		return false, nil
	}

	userDB := db.NewUserDB()
	for i := len(ls) - 1; i >= 0; i-- {
		fromId, fromCash, toId, toCash := extractTransferFromPutLog(tId, ls[i])

		if fromId != -1 && toId != -1 {
			// TODO: Error handle?
			userDB.UpdateCash(fromId, fromCash)
			userDB.UpdateCash(toId, toCash)
		}

	}

	return true, nil
}

func extractTransferFromPutLog(tId int, put string) (int, int, int, int) {
	results := regexp.MustCompile(`^<T(\d+)\D(\d+)\D(\d+)\D(\d+)\D(\d+)>$`).FindStringSubmatch(put)
	if len(results) == 6 {
		// matched
		id, _ := strconv.Atoi(results[1])
		if id == tId {
			// starting undo
			fromId, _ := strconv.Atoi(results[2])
			fromCash, _ := strconv.Atoi(results[3])
			toId, _ := strconv.Atoi(results[4])
			toCash, _ := strconv.Atoi(results[5])
			return fromId, fromCash, toId, toCash
		}
	}
	return -1, -1, -1, -1
}

/**
检查是否需要做 undo 操作
*/
func (l *Log) CheckUndoLog(tId int, ls []string) (bool, error) {

	errC := make(chan error)
	var wg sync.WaitGroup
	for _, l := range ls {
		wg.Add(1)
		go func(ll string) {
			defer wg.Done()
			err := CheckDone(ll, tId)
			if err != nil {
				errC <- err
			}
		}(l)
	}

	go func() {
		wg.Wait()
		close(errC)
	}()

	for err := range errC {
		return false, err
	}

	return true, nil
}

/**
触发一次 undo log的 gc 请求
*/
func (l *Log) GCUndoLog() (bool, error) {

}

// some log file write function
func (l *Log) WriteStart(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<START T%d>", tId))
}
func (l *Log) WritePut(tId, fromId, fromOriginalCash, toId, toOriginalCash int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<T%d,%d,%d,%d,%d>", tId, fromId, fromOriginalCash, toId, toOriginalCash))
}
func (l *Log) WriteCommit(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<COMMIT T%d>", tId))
}
func (l *Log) writeUndo(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<UNDO T%d>", tId))
}
