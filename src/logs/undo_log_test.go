package logs_test

import (
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"logs"
	"os"
	"path"
	"testing"
	"time"
	"transaction"
	"util"
)

func TestCreateSomeLog(t *testing.T) {

	isExisted, err := file.CheckExisted(logs.LOG_PATH)
	if isExisted == false && err == nil {
		// path not exists
		os.Mkdir(logs.LOG_PATH, os.ModePerm)
	} else {
		util.Check(err)
	}

	for i := 1; i < 10; i++ {
		ts := time.Now().Add(time.Hour * time.Duration(i))
		filename := fmt.Sprintf("%d.log", ts.Unix()/3600*3600)
		log.Println(path.Join(logs.LOG_PATH, filename))
		f, err := os.Create(path.Join(logs.LOG_PATH, filename))
		util.Check(err)
		f.Close()
	}
}

func TestCleanLogs(t *testing.T) {
	isExisted, err := file.CheckExisted(logs.LOG_PATH)
	if isExisted == true && err == nil {
		os.RemoveAll(logs.LOG_PATH)
	} else {
		util.Check(err)
	}

}

func TestNewLog(t *testing.T) {
	l := logs.NewLog()
	assert.FileExistsf(t, path.Join(logs.LOG_PATH, l.Logfile), "file must exist")
}

func TestLog_Write(t *testing.T) {
	l := logs.NewLog()

	req := &transaction.Request{
		RequestType: transaction.REQUEST_START,
		Transaction: &transaction.Transaction{
			Trans: []transaction.Transfer{
				{1, 2, 1},
				{3, 4, 1},
				{4, 5, 1}},
		},
	}
	l.Write(req)
	req.RequestType = transaction.REQUEST_PUT
	l.Write(req)
	req.RequestType = transaction.REQUEST_COMMIT
	l.Write(req)

}
