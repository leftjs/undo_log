package main

//
//import (
//	"log"
//	"system"
//)
//
//func main() {
//
//	s := system.NewSystem()
//
//	users := []*system.User{
//		{
//			ID:   1,
//			Name: "Tom",
//			Cash: 10,
//		},
//		{
//			ID:   2,
//			Name: "Jerry",
//			Cash: 10,
//		},
//		{
//			ID:   3,
//			Name: "Spike",
//			Cash: 10,
//		},
//	}
//	transactions := []*system.Transaction{
//		{
//			TransactionID: 1,
//			FromID:        1,
//			ToID:          2,
//			Cash:          10,
//		},
//		{
//			TransactionID: 2,
//			FromID:        2,
//			ToID:          3,
//			Cash:          5,
//		},
//		{
//			TransactionID: 3,
//			FromID:        3,
//			ToID:          1,
//			Cash:          20,
//		},
//		{
//			TransactionID: 4,
//			FromID:        2,
//			ToID:          1,
//			Cash:          10,
//		},
//	}
//
//	for _, user := range users {
//		if err := s.AddUser(user); err != nil {
//			log.Printf("add user failed %v", err)
//		}
//	}
//
//	// TODO: do transaction parallel
//	for _, transaction := range transactions {
//		if err := s.DoTransaction(transaction); err != nil {
//			log.Printf("do transcation failed %v", err)
//		}
//	}
//
//	for _, user := range s.Users {
//		log.Printf("after transcation, %s has %d money", user.Name, user.Cash)
//	}
//
//	if err := s.UndoTransaction(2); err != nil {
//		log.Printf("undo transcation failed %v", err)
//	}
//
//	for _, user := range s.Users {
//		log.Printf("after undo transcation, %s has %d money", user.Name, user.Cash)
//	}
//}
