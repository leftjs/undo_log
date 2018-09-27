package db

const USER_DB_FILE = "../../data/users.db"

type User struct {
	ID   int
	Name string
	Cash int
}

type DB struct {
	users map[int]*User
}

func (db *DB) RegisterUser(u User) {
}
