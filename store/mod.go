package store

import (
	"os"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/estore"
)

const dbPath = "./db"
const sealedKeyFile = "./db/key"

var DB *estore.DB

func InitDB() error {
	encryptionKey := []byte{13, 72, 146, 87, 232, 212, 174, 12, 78, 40, 239, 24, 124, 79, 203, 205}

	// Check if the database exists
	if _, err := os.Stat(dbPath + "/CURRENT"); os.IsNotExist(err) {
		if err := os.Mkdir(dbPath, 0o700); err != nil {
			return err
		}

		sealedKey, err := ecrypto.SealWithUniqueKey(encryptionKey, nil)
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
