package database

import "errors"

var ErrPermission = errors.New("permission denied")

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		AuthorID: authorID,
		Body:     body,
	}
	dbStructure.Chirps[id] = chirp

	if err = db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) GetChirps(authorID int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if authorID != -1 {
			if chirp.AuthorID == authorID {
				chirps = append(chirps, chirp)
			}
		} else {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

func (db *DB) DeleteChirpByID(chirpId, authorID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[chirpId]
	if !ok {
		return errors.New("chirp not found")
	}

	if chirp.AuthorID == authorID {
		delete(dbStructure.Chirps, chirpId)

		if err = db.writeDB(dbStructure); err != nil {
			return err
		}

		return nil
	}

	return ErrPermission
}
