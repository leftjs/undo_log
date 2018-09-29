package db

import (
	"config"
	"file"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	ID   int
	Name string
	Cash int
}

func NewUser(name string, cash int) *User {
	return &User{Name: name, Cash: cash}
}

func NewUserFromString(u string) *User {
	uS := strings.Split(u, ",")
	if len(uS) < 3 {
		return nil
	}
	id, _ := strconv.Atoi(uS[0])
	cash, _ := strconv.Atoi(uS[2])
	return &User{
		id,
		string(uS[1]),
		cash,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("%d,%s,%d", u.ID, u.Name, u.Cash)
}

type UserDB struct {
	mu sync.RWMutex

	users  map[int]*User
	Config *config.Config
}

func NewUserDB() *UserDB {
	cfg := config.NewConfig()
	return &UserDB{
		Config: cfg,
	}
}

/**
can't be called in locking state!!!
*/
func (db *UserDB) GetUser(id int) *User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getUser(id)
}

/**
need to be called after read locking
*/
func (db *UserDB) getUser(id int) *User {

	db.loadUsersFromDBFile()
	return db.users[id]
}

/**
need to be called after locking
*/
func (db *UserDB) loadUsersFromDBFile() {

	data := string(file.ReadFile(db.Config.UserDBFile))
	users := strings.Split(data, "\n")
	if len(db.users) == 0 {
		db.users = make(map[int]*User)
	}
	for _, u := range users {
		user := NewUserFromString(u)
		if user == nil {
			continue
		}
		db.users[user.ID] = user
	}

}

/**
add user
1. add to db file firstly
2. add to memory
*/
func (db *UserDB) AddUser(u *User) {

	db.mu.Lock()
	defer db.mu.Unlock()

	// sync with file
	db.loadUsersFromDBFile()

	userIDs := make([]int, len(db.users))
	for id, _ := range db.users {
		userIDs = append(userIDs, id)
	}
	sort.Ints(userIDs)
	lastId := userIDs[len(userIDs)-1]
	lastId++
	u.ID = lastId

	// 先写文件
	file.AppendToFile(db.Config.UserDBFile, u.String())
	// 再写内存
	if len(db.users) == 0 {
		db.users = make(map[int]*User)
		db.users[u.ID] = u
	}
}

/**
update cash
1. !(cash < 0)
2. id's user must exist
*/
func (db *UserDB) UpdateCash(id, cash int) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 不存在
	if db.getUser(id) == nil {
		return false, errors.New("user doesn't exists")
	}

	// 金额为负
	if cash < 0 {
		return false, errors.New("cash must larger than 0")
	}

	oldContent := db.users[id].String()
	var newUser User
	newUser = *db.users[id]
	newUser.Cash = cash
	newContent := newUser.String()

	file.ReplaceFileLine(db.Config.UserDBFile, oldContent, newContent)
	return true, nil
}
