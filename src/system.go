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

// writeUndoLog writes undo log to file
func (s *System) writeUndoLog(t *Transcation) error {
	// TODO: implement writeUndoLog

	return nil
}

// gcUndoLog the old undo log
func (s *System) gcUndoLog() {
	// TODO: implement gcUndoLog
}

// UndoTranscation roll back some transcations
func (s *System) UndoTranscation(fromID int) error {
	// TODO: implement UndoTranscation
	// undo transcation from fromID to the last transcation

	return nil
}
