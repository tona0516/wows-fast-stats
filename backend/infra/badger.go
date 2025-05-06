package infra

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

func read[T any](db *badger.DB, key string) (T, error) {
	var result T

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return failure.Wrap(err)
		}

		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			return gob.NewDecoder(buf).Decode(&result)
		})

		return failure.Wrap(err)
	})

	return result, err
}

func write[T any](db *badger.DB, key string, target T) error {
	err := db.Update(func(txn *badger.Txn) error {
		buf := bytes.NewBuffer(nil)

		if err := gob.NewEncoder(buf).Encode(&target); err != nil {
			return failure.Wrap(err)
		}

		entry := badger.NewEntry([]byte(key), buf.Bytes())

		return failure.Wrap(txn.SetEntry(entry))
	})

	return err
}

func remove(db *badger.DB, key string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))

		return failure.Wrap(err)
	})
}

func isErrKeyNotFound(err error) bool {
	return err != nil && errors.Is(err, badger.ErrKeyNotFound)
}
