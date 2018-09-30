package transaction_test

import (
	"config"
	"db"
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"logs"
	"path"
	"sync"
	"testing"
	"transaction"
)

func Test_Init(t *testing.T) {
	cfg := config.NewConfig()
	file.DeleteFile(cfg.UserDBFile)

	l := logs.NewLog()
	userDB := db.NewUserDB(l)

	var users []*db.User
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		users = append(users, db.NewUser(fmt.Sprintf("leftjs_%d", i), 100))
		wg.Add(1)
	}
	for _, u := range users {
		go func(uu *db.User) {
			userDB.AddUser(uu)
			wg.Done()
		}(u)
	}

	wg.Wait()

	r := transaction.NewRequest(l, userDB)
	file.DeleteFile(path.Join(r.L.Config.LogPath, r.L.Logfile))
}

func TestRequest_Send(t *testing.T) {
	l := logs.NewLog()
	userDB := db.NewUserDB(l)
	r := transaction.NewRequest(l, userDB)

	trans := []db.Transfer{{1, 2, 1}, {3, 4, 1}}
	r.Send(trans)
	trans = []db.Transfer{{1, 2, 1}, {3, 4, 1}}
	r.Send(trans)
	trans = []db.Transfer{{1, 2, 1}, {3, 4, 1}}
	r.Send(trans)

	assert.Equal(t, 97, r.UserDB.GetUser(1).Cash)
	assert.Equal(t, 103, r.UserDB.GetUser(2).Cash)
	assert.Equal(t, 97, r.UserDB.GetUser(3).Cash)
	assert.Equal(t, 103, r.UserDB.GetUser(4).Cash)
}

func TestRollbackAfter(t *testing.T) {
	l := logs.NewLog()
	userDB := db.NewUserDB(l)
	userDB.RollbackAfter(1)

	assert.Equal(t, 99, userDB.GetUser(1).Cash)
	assert.Equal(t, 101, userDB.GetUser(2).Cash)
	assert.Equal(t, 99, userDB.GetUser(3).Cash)
	assert.Equal(t, 101, userDB.GetUser(4).Cash)

}

func TestRequest_WithErrorRequest(t *testing.T) {
	l := logs.NewLog()

	userDB := db.NewUserDB(l)
	r := transaction.NewRequest(l, userDB)
	trans := []db.Transfer{
		{1, 2, 1},
		{3, 4, 1},
		{4, 5, -1},
		{7, 8, -1},
	}
	r.Send(trans)

	assert.Equal(t, 99, r.UserDB.GetUser(1).Cash)
	assert.Equal(t, 101, r.UserDB.GetUser(2).Cash)
	assert.Equal(t, 99, r.UserDB.GetUser(3).Cash)
	assert.Equal(t, 101, r.UserDB.GetUser(4).Cash)
	assert.Equal(t, 100, r.UserDB.GetUser(5).Cash)
	assert.Equal(t, 100, r.UserDB.GetUser(6).Cash)
	assert.Equal(t, 100, r.UserDB.GetUser(7).Cash)
	assert.Equal(t, 100, r.UserDB.GetUser(8).Cash)
}

func TestRequest_Send_Parallel(t *testing.T) {
	var wg sync.WaitGroup

	l := logs.NewLog()

	userDB := db.NewUserDB(l)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			r := transaction.NewRequest(l, userDB)

			trans := []db.Transfer{{1, 2, 1}, {3, 4, 1}}
			r.Send(trans)
			wg.Done()
		}()
	}

	wg.Wait()
}
