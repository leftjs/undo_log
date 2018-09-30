package main

import (
	"config"
	"db"
	"file"
	"log"
	"logs"
	"path"
	"sync"
	"transaction"
)

func main() {

	//s := system.NewSystem()

	users := []*db.User{
		db.NewUser("Tom", 10),
		db.NewUser("Jerry", 10),
		db.NewUser("Spike", 10),
	}

	trans := []*transaction.Transaction{
		{
			Trans: []db.Transfer{{1, 2, 10}},
		},
		{
			Trans: []db.Transfer{{2, 3, 5}},
		},
		{
			Trans: []db.Transfer{{3, 1, 20}},
		},
		{
			Trans: []db.Transfer{{2, 1, 10}},
		},
	}

	l := logs.NewLog()
	userDB := db.NewUserDB(l)

	// 1. add some user
	addUser(userDB, users)

	log.Println("------before transaction------")
	// 2. list user before transaction
	listUser(userDB)
	log.Println("------before transaction------")
	log.Println("")

	log.Println("do parallel transactions, if error occurs will rollback!!!!")
	// 3. do transaction parallel
	doTransaction(l, userDB, trans)
	log.Println("")
	log.Println("------after transaction------")
	// 4. list user after transaction
	listUser(userDB)
	log.Println("------after transaction------")

}

/**
初始化运行环境,清空 userdb & logfile
*/
func init() {
	cfg := config.NewConfig()
	l := logs.NewLog()
	file.DeleteFile(cfg.UserDBFile)
	file.DeleteFile(path.Join(cfg.LogPath, l.Logfile))
}

func addUser(userDB *db.UserDB, users []*db.User) {
	for _, user := range users {
		userDB.AddUser(user)
	}
}

func listUser(userDB *db.UserDB) {
	users := userDB.Users
	for i := 1; i <= len(users); i++ {
		user := users[i]
		log.Printf("[ID: %d] - %s has %d money", user.ID, user.Name, user.Cash)
	}
}

func doTransaction(l *logs.Log, userDB *db.UserDB, trans []*transaction.Transaction) {
	var wg sync.WaitGroup
	for _, tran := range trans {
		wg.Add(1)
		go func(t *transaction.Transaction) {
			req := transaction.NewRequest(l, userDB)
			req.Send(t.Trans)
			wg.Done()
		}(tran)
	}
	wg.Wait()
}
