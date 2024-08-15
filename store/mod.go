package store

import (
	"errors"
	"fmt"
	"os"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/estore"
	"golang.org/x/crypto/blake2b"
)

const dbPath = "./db"
const sealedKeyFile = "./db/key"

var DB *estore.DB

func InitDB(password string) error {
	var encryptionKey []byte

	// Check if the database exists
	if _, err := os.Stat(sealedKeyFile); os.IsNotExist(err) {
		if err := os.Mkdir(dbPath+"/addon", 0o700); err != nil {
			return err
		}

		if _, err := os.Create(sealedKeyFile); err != nil {
			return err
		}

		if password == "" {
			return errors.New("password is empty")
		}

		hash := blake2b.Sum256([]byte(password))
		encryptionKey = hash[:]

		sealedKey, err := ecrypto.SealWithProductKey(encryptionKey, nil)
		if err != nil {
			return err
		}

		if err := os.WriteFile(sealedKeyFile, sealedKey, 0o600); err != nil {
			return err
		}
	} else {
		sealedKey, err := os.ReadFile(sealedKeyFile)
		if err != nil {
			return err
		}

		encryptionKey, err = ecrypto.Unseal(sealedKey, nil)
		if err != nil {
			return err
		}
	}

	// Create an encrypted store
	opts := &estore.Options{
		EncryptionKey: encryptionKey,
	}

	// Open the store
	var err error
	DB, err = estore.Open(dbPath, opts)

	return err
}

func SetKey(namespace, key string, value []byte) error {
	return DB.Set([]byte(namespace+"__"+key), value, nil)
}

func GetKey(namespace, key string) ([]byte, error) {
	value, _, err := DB.Get([]byte(namespace + "__" + key))
	return value, err
}

func GetList(namespace string, page int, size int) ([][]byte, error) {
	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil // no upper-bound
	}

	prefixIterOptions := func(prefix []byte) *estore.IterOptions {
		return &estore.IterOptions{
			LowerBound: prefix,
			UpperBound: keyUpperBound(prefix),
		}
	}

	iter, err := DB.NewIter(prefixIterOptions([]byte(namespace)))
	if err != nil {
		return nil, err
	}

	value := make([][]byte, 0, page)
	for iter.First(); iter.Valid(); iter.Next() {
		fmt.Printf("%s\n", iter.Value())
		value = append(value, iter.Value())
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return value, nil
}
