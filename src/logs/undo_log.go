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

	CurrentTransactionID int    // 当前事务 id
	logfile              string // 当前 log 写入的文件

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
			l.logfile = name
		}
		return
	}

	sort.SliceStable(files, func(i, j int) bool {
		return extractUnixFromFileName(files[i].Name()) < extractUnixFromFileName(files[j].Name())
	})
	last := files[len(files)-1]
	l.logfile = last.Name()
}

/**
获取下一个 transaction id
*/
func (l *Log) GetNextTransactionId() int {

	data := file.ReadFile(l.logfile)
	logs := strings.Split(string(data), "\n")

	var tIdPtr *int
	// scan last start id, then add 1
	for i := len(logs) - 1; i >= 0; i-- {
		results := regexp.MustCompile(`^<START T(\d)>$`).FindStringSubmatch(logs[i])
		if len(results) > 1 {
			id, err := strconv.Atoi(results[1])
			tIdPtr = &id
			util.Check(err)
			break
		}
	}
	if tIdPtr == nil {
		l.CurrentTransactionID = -1
	} else {
		l.CurrentTransactionID = *tIdPtr
	}

	return l.CurrentTransactionID + 1

}

/**
从文件名字中提取 timestamp
*/
func extractUnixFromFileName(filename string) int {
	i, _ := strconv.Atoi(strings.Split(filename, ".")[0])
	return i
}

/**
写日志请求
*/
func (l *Log) Write(req *transaction.Request) (bool, error) {
	if &req.Transaction.ID == nil {
		req.Transaction.ID = l.GetNextTransactionId()
	}
	t := req.Transaction
	userDB := db.NewUserDB()
	switch req.RequestType {
	case transaction.REQUEST_START:
		writeStart(t.ID)
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
			writePut(t.ID, transfer.FromID, fromCash, transfer.ToID, toCash)

			if _, err = userDB.UpdateCash(transfer.FromID, user.Cash-transfer.Cash); err != nil {
				return false, err
			}
			if _, err = userDB.UpdateCash(transfer.ToID, user.Cash+transfer.Cash); err != nil {
				return false, err
			}
		}
	case transaction.REQUEST_COMMIT:
		writeCommit(t.ID)
	case transaction.REQUEST_UNDO:
		writeUndo(t.ID)
	}

	return true, nil
}

// TODO
func (l *Log) Undo(req *transaction.Request) (bool, error) {
	return false, nil
}

func writeStart(tId int) string {
	return fmt.Sprintf("<START T%d>", tId)
}
func writePut(tId, fromId, fromOriginalCash, toId, toOriginalCash int) string {
	return fmt.Sprintf("<T%d,%d,%d,%d,%d>", tId, fromId, fromOriginalCash, toId, toOriginalCash)
}
func writeCommit(tId int) string {
	return fmt.Sprintf("<COMMIT T%d>", tId)
}
func writeUndo(tId int) string {
	return fmt.Sprintf("<UNDO T%d>", tId)
}
