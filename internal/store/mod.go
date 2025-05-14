package store

import (
	"flag"

	"github.com/cockroachdb/pebble/v2"
	"github.com/edgelesssys/ego/ecrypto"
)

const dbPath = "./db"

var DB *pebble.DB

func InitDB(password string) error {
	var err error
	DB, err = pebble.Open(dbPath, &pebble.Options{})

	return err
}

func SetKey(namespace, key string, value []byte) error {
	val, err := SealWithProductKey(value, nil)
	if err != nil {
		return err
	}

	return DB.Set([]byte(namespace+"__"+key), val, nil)
}

func GetKey(namespace, key string) ([]byte, error) {
	value, _, err := DB.Get([]byte(namespace + "__" + key))
	if err != nil {
		return nil, err
	}

	return ecrypto.Unseal(value, nil)
}

func SealWithProductKey(val []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.SealWithProductKey(val, additionalData)
	}
	return val, nil
}

func Unseal(ciphertext []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.Unseal(ciphertext, additionalData)
	}
	return ciphertext, nil
}
