package transaction

import (
	"db"
	"errors"
	"log"
	"logs"
)

/**
转账
*/
type Transfer struct {
	FromID int
	ToID   int
	Cash   int
}

/**
一次事务中可以有多个转账
*/
type Transaction struct {
	ID    int
	Trans []Transfer
}

type Request struct {
	L           *logs.Log
	UserDB      *db.UserDB
	RequestType logs.RequestType
	Transaction *Transaction
}

func NewRequest() *Request {

	// new a log module to record undo logs
	l := logs.NewLog()
	userDB := db.NewUserDB()
	return &Request{
		L:      l,
		UserDB: userDB,
	}
}

func (r *Request) Send(trans []Transfer) {
	t := &Transaction{
		Trans: trans,
	}

	r.Transaction = t

	// 1. start a transaction
	r.RequestType = logs.REQUEST_START
	r.Write()

	// 2. send some transfer
	r.RequestType = logs.REQUEST_PUT
	status, err := r.Write()

	// 3. send a terminate state
	if status == true {
		// need to commit
		r.RequestType = logs.REQUEST_COMMIT

	} else {
		// need to undo according to log file
		if err != nil {
			log.Println(err.Error())
		}
		r.RequestType = logs.REQUEST_UNDO
	}

	r.Write()
}

/**
transaction request 检查
*/
func (r *Request) checkAndFixTransactionRequest() bool {
	if r.RequestType == logs.REQUEST_START {
		r.Transaction.ID = r.L.GetNextTransactionId()
		return true
	}
	t := r.L.UndoLogs[r.Transaction.ID].Tail.Data.(*logs.Undo).Type
	if r.RequestType == logs.REQUEST_PUT && (t == logs.REQUEST_COMMIT || t == logs.REQUEST_UNDO) {
		return false
	}
	return true

}

/**
写日志请求
*/
func (r *Request) Write() (bool, error) {

	// 检查并修正请求
	status := r.checkAndFixTransactionRequest()
	if !status {
		return false, nil
	}

	t := r.Transaction
	switch r.RequestType {
	case logs.REQUEST_START:
		r.L.WriteStart(t.ID)
	case logs.REQUEST_PUT:
		for _, transfer := range t.Trans {
			status, err := r.doOneTransfer(transfer)
			if status == false {
				return status, err
			}
		}
	case logs.REQUEST_COMMIT:
		r.L.WriteCommit(t.ID)
	case logs.REQUEST_UNDO:
		r.L.Undo(t.ID)
	}

	return true, nil
}

func (r *Request) doOneTransfer(transfer Transfer) (bool, error) {
	var user *db.User
	var err error
	if transfer.Cash <= 0 {
		return false, errors.New("transfer cash must > 0 ")
	}
	if user = r.UserDB.GetUser(transfer.FromID); user == nil {
		return false, errors.New("user doesn't exists")
	}
	fromCash := user.Cash
	if user = r.UserDB.GetUser(transfer.ToID); user == nil {
		return false, errors.New("user doesn't exists")
	}
	toCash := user.Cash
	r.L.WritePut(r.Transaction.ID, transfer.FromID, fromCash, transfer.ToID, toCash)

	if _, err = r.UserDB.UpdateCash(transfer.FromID, fromCash-transfer.Cash); err != nil {
		return false, err
	}
	if _, err = r.UserDB.UpdateCash(transfer.ToID, toCash+transfer.Cash); err != nil {
		return false, err
	}

	return true, nil
}
