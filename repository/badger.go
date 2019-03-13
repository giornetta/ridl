package repository

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/giornetta/ridl"

	"github.com/gofrs/uuid"

	"github.com/dgraph-io/badger"
)

// encode correctly encodes a riddle into a bytes sequence
func encode(r *ridl.Riddle) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decode converts a bytes sequence into a Riddle struct
func decode(b []byte) (*ridl.Riddle, error) {
	var r ridl.Riddle
	buf := bytes.NewBuffer(b)

	if err := gob.NewDecoder(buf).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

// New returns a concrete implementation of the Repository interface
func NewBadger(db *badger.DB) ridl.Repository {
	return &badgerRepository{
		db: db,
	}
}

type badgerRepository struct {
	db *badger.DB
}

func (repo *badgerRepository) Put(r *ridl.Riddle) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := uid.String()[:5]

	b, err := encode(r)
	if err != nil {
		return "", err
	}

	if err := repo.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(id), b)
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (repo *badgerRepository) Get(id string) (*ridl.Riddle, error) {
	var r *ridl.Riddle
	if err := repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		b, err := item.Value()
		if err != nil {
			return err
		}

		r, err = decode(b)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// If the riddle is expired and it hasn't already been deleted
	// delete it and return an error.
	if r.IsExpired() {
		_ = repo.Delete(id)
		return nil, errors.New("expired")
	}

	return r, nil
}

func (repo *badgerRepository) Delete(id string) error {
	return repo.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}

func (repo *badgerRepository) DeleteExpired() error {
	return repo.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// Range through all the values
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			b, err := item.Value()
			if err != nil {
				continue
			}

			r, err := decode(b)
			if err != nil {
				continue
			}

			// Delete record if it has expired
			if r.IsExpired() {
				_ = txn.Delete(k)
			}
		}
		return nil
	})
}
