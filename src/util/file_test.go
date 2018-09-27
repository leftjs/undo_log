package util_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"util"
)

const FILE_TEST = "../../test/test.db"

func TestDeleteFile(t *testing.T) {
	util.DeleteFile(FILE_TEST)
}

func TestAppendToFile(t *testing.T) {
	util.AppendToFile(FILE_TEST, "1\n")
	util.AppendToFile(FILE_TEST, "2\n")
	util.AppendToFile(FILE_TEST, "3")
	util.AppendToFile(FILE_TEST, "4")
	util.AppendToFile(FILE_TEST, "5")
}

func TestReadFile(t *testing.T) {
	bytes := util.ReadFile(FILE_TEST)
	assert.Equal(t, []byte("1\n2\n3\n4\n5\n"), bytes)
}

func TestReplaceFileLine(t *testing.T) {
	util.ReplaceFileLine(FILE_TEST, "1", "10")
	assert.Equal(t, []byte("10\n2\n3\n4\n5\n"), util.ReadFile(FILE_TEST))
}
