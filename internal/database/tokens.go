package database

import (
	"time"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) CreateRefreshToken(userID int, value string, expiresAt time.Time) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStructure.Users[userID]; !ok {
		return ErrNotExist
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Value:     value,
		ExpiresAt: expiresAt,
	}
	dbStructure.RefreshTokens[userID] = refreshToken

	if err = db.writeDB(dbStructure); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRefreshTokenByValue(value string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	for _, rt := range dbStructure.RefreshTokens {
		if rt.Value == value {
			if rt.ExpiresAt.After(time.Now().UTC()) {
				return rt, nil
			}
		}
	}

	return RefreshToken{}, ErrNotExist
}

func (db *DB) DeleteRefreshTokenByValue(value string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	for _, rt := range dbStructure.RefreshTokens {
		if rt.Value == value {
			delete(dbStructure.RefreshTokens, rt.UserID)

			if err = db.writeDB(dbStructure); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrNotExist
}
