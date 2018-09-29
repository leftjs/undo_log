package main

import (
	"db"
	"sync"
)

func main() {

	//s := system.NewSystem()

	userDB := db.NewUserDB()

	users := []*db.User{
		db.NewUser("Tom", 10),
		db.NewUser("Jerry", 10),
		db.NewUser("Spike", 10),
	}

	//transactions := []*system.Transaction{
	//	{
	//		TransactionID: 1,
	//		FromID:        1,
	//		ToID:          2,
	//		Cash:          10,
	//	},
	//	{
	//		TransactionID: 2,
	//		FromID:        2,
	//		ToID:          3,
	//		Cash:          5,
	//	},
	//	{
	//		TransactionID: 3,
	//		FromID:        3,
	//		ToID:          1,
	//		Cash:          20,
	//	},
	//	{
	//		TransactionID: 4,
	//		FromID:        2,
	//		ToID:          1,
	//		Cash:          10,
	//	},
	//}

	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go func(u *db.User) {
			userDB.AddUser(u)
			wg.Done()
		}(user)
	}
	wg.Wait()

	//// TODO: do transaction parallel
	//for _, transaction := range transactions {
	//	if err := s.DoTransaction(transaction); err != nil {
	//		log.Printf("do transcation failed %v", err)
	//	}
	//}
	//
	//for _, user := range s.users {
	//	log.Printf("after transcation, %s has %d money", user.Name, user.Cash)
	//}
	//
	//if err := s.UndoTransaction(2); err != nil {
	//	log.Printf("undo transcation failed %v", err)
	//}
	//
	//for _, user := range s.users {
	//	log.Printf("after undo transcation, %s has %d money", user.Name, user.Cash)
	//}
}
