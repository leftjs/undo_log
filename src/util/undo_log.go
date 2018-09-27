package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"transaction"
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

const (
	START_TPL  = "<START T%d>"
	PUT_TPL    = "<T%d, %s, %d>"
	COMMIT_TPL = "<COMMIT T%d>"
)

type Log struct {
	mu sync.RWMutex

	CurrentTransactionID int    // 当前事务 id
	logfile              string // 当前 log 写入的文件

}

func NewUndoLog() *Log {
	existed, err := CheckExisted(LOG_PATH)
	var tId int

	if existed == true {
		// 读文件
		tId = 0
	} else if err == nil {
		// 文件夹不存在，需要创建
		//makeLogDir()
		tId = 1
	} else {
		Check(err)
	}

	return &Log{
		CurrentTransactionID: tId,
	}

}

func WriteLog() {
	d := []byte("hello\ngo\n")
	f, err := os.Create("1.log")
	Check(err)

	defer f.Close()

	d = []byte("haha\ngo\n")

	f.Write(d)
}

// 从 log_path 的 undo logs 中 scan 出最后一个 transaction id
func getLastTransactionId() int {
	// log path 不存在则返回 id = 0
	if existed, _ := CheckExisted(LOG_PATH); existed == false {
		return 0
	}
	// 挑选最新的 log file
	files, err := ioutil.ReadDir(LOG_PATH)
	Check(err)

	if len(files) == 0 {
		return 0
	}
	sort.SliceStable(files, func(i, j int) bool {
		return extractUnixFromFileName(files[i].Name()) < extractUnixFromFileName(files[j].Name())
	})

	//last := files[len(files) -1]
	return 1

}

func (l *Log) seekLog(filename string) {

	l.mu.RLock()
	defer l.mu.RUnlock()

	f, err := os.Open(filename)
	Check(err)
	defer f.Close()

	//r := bufio.NewReader(f)
}

func (l *Log) WriteALog(request transaction.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.Open(l.logfile)
	defer f.Close()
	Check(err)

}

func extractUnixFromFileName(filename string) int {
	i, _ := strconv.Atoi(strings.Split(filename, ".")[0])
	return i
}

func writeStart(tId int) string {
	return fmt.Sprintf("<START %d>", tId)
}
func writeUpdate(tId int, x string, v int) string {
	return fmt.Sprintf("<%d,%s,%d>", tId, x, v)
}
func writeCommit(tId int) string {
	return fmt.Sprintf("<COMMIT %d>", tId)
}
