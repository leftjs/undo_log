package transaction

import (
	"db"
	"errors"
	"file"
	"log"
	"logs"
	"path"
	"strings"
	"sync"
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

const (
	REQUEST_START RequestType = iota
	REQUEST_PUT
	REQUEST_COMMIT
	REQUEST_UNDO
)

type RequestType int

type Request struct {
	L           *logs.Log
	UserDB      *db.UserDB
	RequestType RequestType
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

	// 1. begin a transaction
	r.RequestType = REQUEST_START
	r.Write()

	// 2. send some transfer
	r.RequestType = REQUEST_PUT
	status, err := r.Write()

	// 3. send a terminate state
	if status == true {
		// need to commit
		r.RequestType = REQUEST_COMMIT

	} else {
		// need to undo according to log file
		if err != nil {
			log.Println(err.Error())
		}
		r.RequestType = REQUEST_UNDO
	}

	r.Write()
}

/**
transaction request 检查
*/
func (r *Request) checkAndFixTransactionRequest() error {
	if r.RequestType == REQUEST_START {
		r.Transaction.ID = r.L.GetNextTransactionId()
		return nil
	}

	data := file.ReadFile(path.Join(logs.LOG_PATH, r.L.Logfile))
	ls := strings.Split(strings.Trim(string(data), "\n"), "\n")
	if len(ls) == 1 && ls[0] == strings.Trim(string(data), "\n") {
		// 空
		return nil
	}

	t := r.Transaction

	errC := make(chan error)
	var wg sync.WaitGroup

	for i := len(ls) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(ii int) {
			err := logs.CheckDone(ls[ii], t.ID)
			if err != nil {
				errC <- err
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(errC)
	}()

	for e := range errC {
		if e != nil {
			return e
		}
	}

	return nil

}

/**
写日志请求
*/
func (r *Request) Write() (bool, error) {

	// 检查并修正请求
	r.checkAndFixTransactionRequest()

	t := r.Transaction
	switch r.RequestType {
	case REQUEST_START:
		r.L.WriteStart(t.ID)
	case REQUEST_PUT:

		for _, transfer := range t.Trans {
			status, err := r.doOneTransfer(transfer)
			if status == false {
				return status, err
			}
		}

	case REQUEST_COMMIT:
		r.L.WriteCommit(t.ID)
	case REQUEST_UNDO:
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
