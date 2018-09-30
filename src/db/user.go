package db

import (
	"config"
	"datastructure"
	"file"
	"fmt"
	"github.com/pkg/errors"
	"logs"
	"path"
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

/**
转账
*/
type Transfer struct {
	FromID int
	ToID   int
	Cash   int
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

	Users  map[int]*User
	Config *config.Config
	L      *logs.Log
}

func NewUserDB(l *logs.Log) *UserDB {
	cfg := config.NewConfig()

	db := &UserDB{
		Config: cfg,
		L:      l,
	}
	db.Users = make(map[int]*User)
	db.loadUsersFromDBFile()
	return db
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

	return db.Users[id]
}

/**
need to be called after locking
*/
func (db *UserDB) loadUsersFromDBFile() {

	data := string(file.ReadFile(db.Config.UserDBFile))
	users := strings.Split(data, "\n")
	if len(db.Users) == 0 {
		db.Users = make(map[int]*User)
	}
	for _, u := range users {
		user := NewUserFromString(u)
		if user == nil {
			continue
		}
		db.Users[user.ID] = user
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

	userIDs := make([]int, len(db.Users))
	for id, _ := range db.Users {
		userIDs = append(userIDs, id)
	}
	sort.Ints(userIDs)
	lastId := userIDs[len(userIDs)-1]
	lastId++
	u.ID = lastId

	// 先写文件
	file.AppendToFile(db.Config.UserDBFile, u.String())
	// 再写内存
	db.Users[u.ID] = u
}

/**
update cash
1. !(cash < 0)
2. id's user must exist
*/
func (db *UserDB) updateCash(id, cash int) {

	oldContent := db.Users[id].String()
	var newUser User
	newUser = *db.Users[id]
	newUser.Cash = cash
	newContent := newUser.String()

	file.ReplaceFileLine(db.Config.UserDBFile, oldContent, newContent)
	db.Users[id] = &newUser
}

func (db *UserDB) UpdateCash(id, cash int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.updateCash(id, cash)
}

func (db *UserDB) DoTransfers(tId int, trans []Transfer) (bool, error) {

	db.mu.Lock()
	defer db.mu.Unlock()

	if db.L.UndoLogs[tId].Head.Data.(*logs.Undo).Type != logs.REQUEST_START {
		return false, nil
	}

	for _, transfer := range trans {
		if transfer.Cash <= 0 {
			return false, errors.New("transfer cash must > 0 ")
		}
		fromUser := db.getUser(transfer.FromID)
		if fromUser == nil {
			return false, errors.New("user doesn't exists")
		}
		toUser := db.getUser(transfer.ToID)

		if toUser == nil {
			return false, errors.New("user doesn't exists")
		}
		db.L.Write(logs.REQUEST_PUT, tId, fromUser.ID, fromUser.Cash, toUser.ID, toUser.Cash)

		if cash := fromUser.Cash - transfer.Cash; cash < 0 {
			return false, errors.New("after transfer from user cash < 0")
		} else {
			db.updateCash(fromUser.ID, fromUser.Cash-transfer.Cash)
		}

		db.updateCash(toUser.ID, toUser.Cash+transfer.Cash)

	}

	return true, nil
}

func (db *UserDB) Undo(tId int) {

	undo := db.L.UndoLogs[tId]
	if undo == nil || !(undo.Tail.Data.(*logs.Undo).Type == logs.REQUEST_PUT || undo.Tail.Data.(*logs.Undo).Type == logs.REQUEST_START) {
		return
	}

	cur := undo.Tail

RollBack:
	for {
		for _, state := range cur.Data.(*logs.Undo).States {
			if &state == nil {
				continue
			}
			db.updateCash(state.UserId, state.Cash)
		}
		if cur.HasPrev() && cur.Prev.Data.(*logs.Undo).Type != logs.REQUEST_START {
			cur = cur.Prev
		} else {
			break RollBack
		}
	}

	db.L.Write(logs.REQUEST_UNDO, tId)
}

/**
rollback transaction after special id
TODO: do some user and log flush
*/
func (db *UserDB) RollbackAfter(tId int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 触发一个 undo 检查
	db.triggerUndo()

	undos := db.L.UndoLogs

	if tId > len(undos) {
		return
	}
	// 收集需要做 rollback 的事务链表
	needToRBUndoes := make(map[int]*datastructure.LinkedList)
	for id, undo := range undos {
		if id > tId {
			needToRBUndoes[id] = undo
		}
	}

	// 强行 rollback - -!!
	for i := len(needToRBUndoes) + tId; i > tId; i-- {
		cur := needToRBUndoes[i].Tail
	RollBack:
		for {
			for _, state := range cur.Data.(*logs.Undo).States {
				if &state == nil {
					continue
				}
				db.updateCash(state.UserId, state.Cash)
			}
			if cur.HasPrev() && cur.Prev.Data.(*logs.Undo).Type != logs.REQUEST_START {
				cur = cur.Prev
			} else {
				break RollBack
			}

		}
	}

}

/**
gc undo log
*/
func (db *UserDB) GCUndoLog() {

	db.mu.Lock()
	defer db.mu.Unlock()

	db.triggerUndo()

	// gc current log file
	file.DeleteFile(path.Join(db.Config.LogPath, db.L.Logfile))
	db.L.InitializeLastLogfile()
	db.L.BuildUndoLogs()
}

/**
trigger undo action
*/
func (db *UserDB) triggerUndo() {

	var wg sync.WaitGroup

	for _, log := range db.L.UndoLogs {
		wg.Add(1)
		undo := log.Tail.Data.(*logs.Undo)
		go func(tId int) {
			if undo.Type != logs.REQUEST_COMMIT && undo.Type != logs.REQUEST_UNDO {
				// need undo,  sync call
				db.Undo(tId)
			}
			wg.Done()
		}(undo.ID)
	}

	wg.Wait()
}
