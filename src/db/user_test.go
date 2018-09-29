package db_test

import (
	"config"
	"db"
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

func TestRemoveUserDBFile(t *testing.T) {
	cfg := config.NewConfig()
	file.DeleteFile(cfg.UserDBFile)
}

func TestUserDB_AddUser(t *testing.T) {
	userDB := db.NewUserDB()
	userDB.AddUser(db.NewUser("leftjs", 100))
	assert.Equal(t, []byte("1,leftjs,100\n"), file.ReadFile(userDB.Config.UserDBFile))
}

func TestUserDB_UpdateCash(t *testing.T) {
	userDB := db.NewUserDB()
	userDB.UpdateCash(1, 20)
	assert.Equal(t, []byte("1,leftjs,20\n"), file.ReadFile(userDB.Config.UserDBFile))
}

const SIZE = 1000 // 并发写入数据点数

func TestUserDB_AddUser_Concurrent(t *testing.T) {
	userDB := db.NewUserDB()
	file.DeleteFile(userDB.Config.UserDBFile)

	var users []*db.User
	var wg sync.WaitGroup

	for i := 1; i < SIZE; i++ {
		users = append(users, db.NewUser(fmt.Sprintf("leftjs_%d", i), rand.Intn(100)))
		wg.Add(1)
	}
	for _, u := range users {
		go func(uu *db.User) {
			userDB.AddUser(uu)
			wg.Done()
		}(u)
	}

	wg.Wait()

	assert.Equal(t, SIZE, len(strings.Split(string(file.ReadFile(userDB.Config.UserDBFile)), "\n")))
}

func TestUserDB_UpdateCash_Concurrent(t *testing.T) {
	userDB := db.NewUserDB()
	var wg sync.WaitGroup
	wg.Add(SIZE)
	for i := 0; i < SIZE; i++ {
		go func(ii int) {
			userDB.UpdateCash(ii, 10)
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 1; i < SIZE; i++ {
		assert.Equal(t, 10, userDB.GetUser(i).Cash)
	}

}

func Test_POST(t *testing.T) {
	cfg := config.NewConfig()
	file.DeleteFile(cfg.UserDBFile)
}
