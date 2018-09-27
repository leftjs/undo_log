package util_test

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"
	"util"
)

func TestCreateSomeLog(t *testing.T) {

	isExisted, err := util.CheckExisted(util.LOG_PATH)
	if isExisted == false && err == nil {
		// path not exists
		os.Mkdir(util.LOG_PATH, os.ModePerm)
	} else {
		util.Check(err)
	}

	for i := 1; i < 10; i++ {
		ts := time.Now().Add(time.Hour * time.Duration(i))
		filename := fmt.Sprintf("%d.log", ts.Unix()/3600*3600)
		log.Println(path.Join(util.LOG_PATH, filename))
		f, err := os.Create(path.Join(util.LOG_PATH, filename))
		util.Check(err)
		f.Close()
	}
}

func TestCleanLogs(t *testing.T) {
	isExisted, err := util.CheckExisted(util.LOG_PATH)
	if isExisted == true && err == nil {
		os.RemoveAll(util.LOG_PATH)
	} else {
		util.Check(err)
	}
}

func TestNewUndoLog(t *testing.T) {
}

func TestWriteLog(t *testing.T) {
	util.WriteLog()
}
