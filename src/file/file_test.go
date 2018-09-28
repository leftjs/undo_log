package file_test

import (
	"file"
	"github.com/stretchr/testify/assert"
	"testing"
)

const FILE_TEST = "../../test/test.db"

func TestDeleteFile(t *testing.T) {
	file.DeleteFile(FILE_TEST)
}

func TestAppendToFile(t *testing.T) {
	file.AppendToFile(FILE_TEST, "1\n")
	file.AppendToFile(FILE_TEST, "2\n")
	file.AppendToFile(FILE_TEST, "3")
	file.AppendToFile(FILE_TEST, "4")
	file.AppendToFile(FILE_TEST, "5")
}

func TestReadFile(t *testing.T) {
	bytes := file.ReadFile(FILE_TEST)
	assert.Equal(t, []byte("1\n2\n3\n4\n5\n"), bytes)
}

func TestReplaceFileLine(t *testing.T) {
	file.ReplaceFileLine(FILE_TEST, "1", "10")
	assert.Equal(t, []byte("10\n2\n3\n4\n5\n"), file.ReadFile(FILE_TEST))
}
