package transaction

import (
	"db"
	"log"
	"logs"
)

/**
一次事务中可以有多个转账
*/
type Transaction struct {
	ID    int
	Trans []db.Transfer
}

type Request struct {
	L           *logs.Log
	UserDB      *db.UserDB
	RequestType logs.RequestType
	Transaction *Transaction
}

func NewRequest(l *logs.Log, userDB *db.UserDB) *Request {

	// new a log module to record undo logs
	return &Request{
		L:      l,
		UserDB: userDB,
	}
}

func (r *Request) Send(trans []db.Transfer) {
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
		log.Printf("transaction: T%d need to undo", r.Transaction.ID)
		r.RequestType = logs.REQUEST_UNDO
	}

	r.Write()
}

/**
transaction request 检查
*/
func (r *Request) checkAndFixTransactionRequest() bool {

	if r.RequestType == logs.REQUEST_START {
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
	s := r.checkAndFixTransactionRequest()
	if !s {
		return false, nil
	}

	t := r.Transaction
	switch r.RequestType {
	case logs.REQUEST_START:
		r.Transaction.ID = r.L.Write(logs.REQUEST_START, t.ID)
	case logs.REQUEST_PUT:
		status, err := r.UserDB.DoTransfers(t.ID, t.Trans)
		if status == false {
			return status, err
		}
	case logs.REQUEST_COMMIT:
		r.L.Write(logs.REQUEST_COMMIT, t.ID)
	case logs.REQUEST_UNDO:
		r.UserDB.Undo(t.ID)
	}

	return true, nil
}
