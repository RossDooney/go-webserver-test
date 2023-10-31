package database

import "time"

type RevokeToken struct {
	ID         string    `json:"id"`
	RevokeTime time.Time `json:"time"`
}

func (db *DB) CreateRevokeToken(id string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	token := RevokeToken{
		ID:         id,
		RevokeTime: time.Now().UTC(),
	}
	dbStructure.RevokeToken[id] = token

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRevokeToken() (RevokeToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RevokeToken{}, err
	}

	token, ok := dbStructure.RevokeToken[""]
	if !ok {
		return RevokeToken{}, ErrNotExist
	}

	return token, nil
}
