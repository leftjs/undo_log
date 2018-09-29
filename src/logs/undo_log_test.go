package logs_test

import (
	"config"
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"logs"
	"os"
	"path"
	"testing"
	"time"
	"util"
)

func TestCreateSomeLog(t *testing.T) {
	cfg := config.NewConfig()
	isExisted, err := file.CheckExisted(cfg.LogPath)
	if isExisted == false && err == nil {
		// path not exists
		os.Mkdir(cfg.LogPath, os.ModePerm)
	} else {
		util.Check(err)
	}

	for i := 1; i < 10; i++ {
		ts := time.Now().Add(time.Hour * time.Duration(i))
		filename := fmt.Sprintf("%d.log", ts.Unix()/3600*3600)
		f, err := os.Create(path.Join(cfg.LogPath, filename))
		util.Check(err)
		f.Close()
	}
}

func TestCleanLogs(t *testing.T) {
	cfg := config.NewConfig()
	isExisted, err := file.CheckExisted(cfg.LogPath)
	if isExisted == true && err == nil {
		os.RemoveAll(cfg.LogPath)
	} else {
		util.Check(err)
	}

}

func TestNewLog(t *testing.T) {
	l := logs.NewLog()
	assert.NotNil(t, l.UndoLogs)
	assert.FileExistsf(t, path.Join(l.Config.LogPath, l.Logfile), "file must exist")
}
