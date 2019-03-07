package repository

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/dgraph-io/badger"
)

type Riddle struct {
	Question     string
	Crypted      []byte
	IgnoreCase   bool
	IgnoreSpaces bool
	Expiry       time.Time
}

func (r *Riddle) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Riddle) IsExpired() bool {
	return time.Now().After(r.Expiry)
}

func DecodeRecord(b []byte) (*Riddle, error) {
	var r Riddle
	buf := bytes.NewBuffer(b)

	if err := gob.NewDecoder(buf).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

type Repository interface {
	Put(r *Riddle) (string, error)
	Get(id string) (*Riddle, error)
	Delete(id string) error
	DeleteExpired() error
}

func New(db *badger.DB) Repository {
	return &badgerRepository{
		db: db,
	}
}

type badgerRepository struct {
	db *badger.DB
}

func (repo *badgerRepository) Put(r *Riddle) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := uid.String()[:5]

	b, err := r.Encode()
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

func (repo *badgerRepository) Get(id string) (*Riddle, error) {
	var r *Riddle
	if err := repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		b, err := item.Value()
		if err != nil {
			return err
		}

		r, err = DecodeRecord(b)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if r.IsExpired() {
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
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			b, err := item.Value()
			if err != nil {
				return err
			}

			r, err := DecodeRecord(b)
			if err != nil {
				return err
			}

			if r.IsExpired() {
				txn.Delete(k)
			}
		}
		return nil
	})
}
