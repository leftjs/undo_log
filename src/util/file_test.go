package util_test

import (
	"db"
	"github.com/stretchr/testify/assert"
	"testing"
	"util"
)

func TestDeleteFile(t *testing.T) {
	util.DeleteFile(db.USER_DB_FILE)
}

func TestAppendToFile(t *testing.T) {
	util.AppendToFile(db.USER_DB_FILE, "1\n")
	util.AppendToFile(db.USER_DB_FILE, "2\n")
	util.AppendToFile(db.USER_DB_FILE, "3")
	util.AppendToFile(db.USER_DB_FILE, "4")
	util.AppendToFile(db.USER_DB_FILE, "5")
}

func TestReadFile(t *testing.T) {
	bytes := util.ReadFile(db.USER_DB_FILE)
	assert.Equal(t, []byte("1\n2\n3\n4\n5\n"), bytes)
}
