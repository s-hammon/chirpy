package database

import (
	"errors"
	"fmt"
)

var ErrExists = errors.New("already exists")

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email string, pwd string) (User, error) {
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

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) UpdateUser(id int, email string, pwd string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.Password = pwd
	dbStructure.Users[id] = user

	if err = db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}
