package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.validate()
	return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		fmt.Printf("error creating database: %v", err)
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		fmt.Printf("error writing to database: %v", err)
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) UpdateUser(id int, email string, pwd []byte) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		fmt.Printf("error creating database: %v", err)
		return User{}, err
	}

	updatedUser := User{
		ID:       id,
		Email:    email,
		Password: pwd,
	}

	dbStructure.Users[id] = updatedUser

	if err = db.writeDB(dbStructure); err != nil {
		fmt.Printf("error writing to database: %v", err)
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) CreateUser(email string, pwd []byte) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		fmt.Printf("error creating database: %v", err)
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if email == user.Email {
			return User{}, errors.New("user with that email already exists")
		}
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: pwd,
	}
	dbStructure.Users[id] = user

	if err = db.writeDB(dbStructure); err != nil {
		fmt.Printf("error writing to database: %v", err)
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	if err = json.Unmarshal(dat, &dbStructure); err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	if err = os.WriteFile(db.path, dat, 0600); err != nil {
		return err
	}

	return nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}

	return db.writeDB(dbStructure)
}

func (db *DB) validate() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}

	return err
}
