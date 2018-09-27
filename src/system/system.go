package system

import (
	"errors"
	"sync"
	"transaction"
)

// User saves user's information
type User struct {
	ID   int
	Name string
	Cash int
}

// System keeps the user and transaction information
type System struct {
	sync.RWMutex

	Users map[int]*User

	Transactions []*transaction.Transaction

	// TODO: add some variables about undo log
}

// NewSystem returns a System
func NewSystem() *System {
	return &System{
		Users:        make(map[int]*User),
		Transactions: make([]*transaction.Transaction, 0, 10),
	}
}

// AddUser adds a new user to the system
func (s *System) AddUser(u *User) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Users[u.ID]; ok {
		return errors.New("user id is already exists")
	}

	s.Users[u.ID] = u

	return nil
}

// DoTransaction apply a transaction
func (s *System) DoTransaction(t *transaction.Transaction) error {
	// TODO: implement DoTransaction
	// if after this transaction, user's cash is less than zero,
	// rollback this transaction according to undo log.

	return nil
}

// writeUndoLog writes undo log to file
func (s *System) writeUndoLog(t *transaction.Transaction) error {
	// TODO: implement writeUndoLog

	return nil
}

// gcUndoLog the old undo log
func (s *System) gcUndoLog() {
	// TODO: implement gcUndoLog
}

// UndoTransaction roll back some transactions
func (s *System) UndoTransaction(fromID int) error {
	// TODO: implement UndoTransaction
	// undo Transaction from fromID to the last transaction

	return nil
}
