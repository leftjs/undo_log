package main

import (
	"errors"
	"sync"
)

// User saves user's information
type User struct {
	ID   int
	Name string
	Cash int
}

// Transcation record a transcation.
type Transcation struct {
	TranscationID int
	FromID        int
	ToID          int
	Cash          int
}

// System keeps the user and transcation information
type System struct {
	sync.RWMutex

	Users map[int]*User

	Transcations []*Transcation

	// TODO: add some variables about undo log
}

// NewSystem returns a System
func NewSystem() *System {
	return &System{
		Users:        make(map[int]*User),
		Transcations: make([]*Transcation, 0, 10),
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

// DoTransaction applys a transaction
func (s *System) DoTransaction(t *Transcation) error {
	// TODO: implement DoTransaction
	// if after this transcation, user's cash is less than zero,
	// rollback this transcation according to undo log.

	return nil
}

// WriteUndoLog writes undo log to file
func (s *System) WriteUndoLog(t *Transcation) error {
	// TODO: implement DoTransaction

	return nil
}

// GCUndoLog the old undo log
func (s *System) GCUndoLog() {
	// TODO: implement DoTransaction
}