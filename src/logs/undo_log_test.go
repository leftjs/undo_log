package logs_test

import (
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"logs"
	"os"
	"path"
	"regexp"
	"testing"
	"time"
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
	assert.NotNil(t, l.UndoLogs)
	assert.FileExistsf(t, path.Join(logs.LOG_PATH, l.Logfile), "file must exist")
}

func Test_Test(t *testing.T) {
	results := regexp.MustCompile(`^<(COMMIT|UNDO|START)?\s?T(\d+)(,(\d+),(\d+),(\d+),(\d+))?>$`).FindStringSubmatch("<START T1>")
	log.Println(len(results))
	log.Println(results[1])
}
