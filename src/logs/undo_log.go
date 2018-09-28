package logs

import (
	"db"
	"file"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"transaction"
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
		log.Println(id)
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

/**
transaction request 检查
*/
func (l *Log) checkAndFixTransactionRequest(req *transaction.Request) error {
	if req.RequestType == transaction.REQUEST_START {
		req.Transaction.ID = l.GetNextTransactionId()
	}

	data, err := ioutil.ReadFile(path.Join(LOG_PATH, l.Logfile))
	util.Check(err)
	ls := strings.Split(string(data), "\n")
	if len(ls) == 1 && ls[0] == string(data) {
		// 空
		return nil
	}

	t := req.Transaction

	errC := make(chan error)
	var wg sync.WaitGroup

	for i := len(ls) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(ii int) {
			err := checkDone(ls[ii], t)
			if err != nil {
				errC <- err
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	select {
	case e := <-errC:
		return e
	default:
		return nil
	}
}

func checkDone(s string, t *transaction.Transaction) error {

	results := regexp.MustCompile(`^<(COMMIT|UNDO) T(\d+)>$`).FindStringSubmatch(s)
	if len(results) == 3 {
		id, _ := strconv.Atoi(results[2])
		if id == t.ID {
			// 该事务已经结束
			return errors.New("transaction has been submitted or cancelled")
		}
	}
	return nil
}

/**
写日志请求
*/
func (l *Log) Write(req *transaction.Request) (bool, error) {

	// 检查并修正请求
	l.checkAndFixTransactionRequest(req)

	t := req.Transaction
	userDB := db.NewUserDB()
	switch req.RequestType {
	case transaction.REQUEST_START:
		l.writeStart(t.ID)
	case transaction.REQUEST_PUT:
		for _, transfer := range t.Trans {
			var user *db.User
			var err error
			if user = userDB.GetUser(transfer.FromID); user == nil {
				return false, errors.New("user doesn't exists")
			}
			fromCash := user.Cash
			if user = userDB.GetUser(transfer.ToID); user == nil {
				return false, errors.New("user doesn't exists")
			}
			toCash := user.Cash
			l.writePut(t.ID, transfer.FromID, fromCash, transfer.ToID, toCash)

			if _, err = userDB.UpdateCash(transfer.FromID, fromCash-transfer.Cash); err != nil {
				return false, err
			}
			if _, err = userDB.UpdateCash(transfer.ToID, toCash+transfer.Cash); err != nil {
				return false, err
			}
		}
	case transaction.REQUEST_COMMIT:
		l.writeCommit(t.ID)
	case transaction.REQUEST_UNDO:
		l.writeUndo(t.ID)
	}

	return true, nil
}

// TODO
func (l *Log) Undo(req *transaction.Request) (bool, error) {
	return false, nil
}

func (l *Log) writeStart(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<START T%d>", tId))
}
func (l *Log) writePut(tId, fromId, fromOriginalCash, toId, toOriginalCash int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<T%d,%d,%d,%d,%d>", tId, fromId, fromOriginalCash, toId, toOriginalCash))
}
func (l *Log) writeCommit(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<COMMIT T%d>", tId))
}
func (l *Log) writeUndo(tId int) {
	file.AppendToFile(path.Join(LOG_PATH, l.Logfile), fmt.Sprintf("<UNDO T%d>", tId))
}
