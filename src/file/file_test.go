package file_test

import (
	"config"
	"file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	cfg := config.NewConfig()
	file.DeleteFile(cfg.TestDBFile)
}

func TestAppendToFile(t *testing.T) {
	cfg := config.NewConfig()

	file.AppendToFile(cfg.TestDBFile, "1\n")
	file.AppendToFile(cfg.TestDBFile, "2\n")
	file.AppendToFile(cfg.TestDBFile, "3")
	file.AppendToFile(cfg.TestDBFile, "4")
	file.AppendToFile(cfg.TestDBFile, "5")
}

func TestReadFile(t *testing.T) {
	cfg := config.NewConfig()

	bytes := file.ReadFile(cfg.TestDBFile)
	assert.Equal(t, []byte("1\n2\n3\n4\n5\n"), bytes)
}

func TestReplaceFileLine(t *testing.T) {
	cfg := config.NewConfig()

	file.ReplaceFileLine(cfg.TestDBFile, "1", "10")
	assert.Equal(t, []byte("10\n2\n3\n4\n5\n"), file.ReadFile(cfg.TestDBFile))
}

func TestMakeDir(t *testing.T) {
}
