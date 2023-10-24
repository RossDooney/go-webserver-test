package database

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash []byte
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	checkEmail := db.CheckEmail(email, dbStructure)
	if len(checkEmail.Email) != 0 {
		return User{}, ErrUsersAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
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

func (db *DB) CheckUserLogin(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if err != nil {
		return User{}, err
	}
	user := db.CheckEmail(email, dbStructure)
	if len(user.Email) == 0 {
		return User{}, ErrIncorrectLogin
	} else {
		err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
		if err != nil {
			return User{}, ErrIncorrectLogin
		}
		return user, nil
	}
}

func (db *DB) CheckEmail(email string, dbStructure DBStructure) User {
	for _, user := range dbStructure.Users {
		if email == user.Email {
			return user
		}
	}
	return User{}
}
