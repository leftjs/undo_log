package transaction_test

import (
	"db"
	"file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path"
	"sync"
	"testing"
	"transaction"
)

func Test_Init(t *testing.T) {
	userDB := db.NewUserDB()

	file.DeleteFile(userDB.Config.UserDBFile)

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

	r := transaction.NewRequest()
	file.DeleteFile(path.Join(r.L.Config.LogPath, r.L.Logfile))
}

func TestRequest_Send(t *testing.T) {
	r := transaction.NewRequest()

	trans := []transaction.Transfer{{1, 2, 1}, {3, 4, 1}}
	r.Send(trans)

	assert.Equal(t, 99, r.UserDB.GetUser(1).Cash)
	assert.Equal(t, 101, r.UserDB.GetUser(2).Cash)
	assert.Equal(t, 99, r.UserDB.GetUser(3).Cash)
	assert.Equal(t, 101, r.UserDB.GetUser(4).Cash)
}

func TestRequest_WithErrorRequest(t *testing.T) {

	r := transaction.NewRequest()
	trans := []transaction.Transfer{
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
