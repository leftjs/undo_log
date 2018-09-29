package config

import (
	"path"
	"runtime"
)

type Config struct {
	LogPath    string
	UserDBFile string
	TestDBFile string
}

func NewConfig() *Config {
	dbFile := "../../data/users.db"
	logPath := "../../log/"
	testDb := "../../test/test.db"

	_, localFile, _, _ := runtime.Caller(0)

	return &Config{
		LogPath:    path.Join(path.Dir(localFile), logPath),
		UserDBFile: path.Join(path.Dir(localFile), dbFile),
		TestDBFile: path.Join(path.Dir(localFile), testDb),
	}
}
