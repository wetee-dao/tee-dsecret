package model

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/dgraph-io/badger/v4"
)

type User struct {
	Name          string
	PubKey        ed25519.PubKey `badgerhold:"index"` // this is just a wrapper around bytes
	Moderator     bool
	Banned        bool
	NumMessages   int64
	Version       uint64
	SchemaVersion int
}

type PublicUser struct {
	// User SS58 address
	Address string `json:"address"`
	// User sign time
	Timestamp int64 `json:"timestamp"`
}

func CreateUser(user *User) error {
	_, err := GetKey("", user.Name)
	if err == nil {
		return errors.New("user already exists")
	}

	return SetJson("", user.Name, user)
}

func FindUserByName(name string) (*User, error) {
	user, err := GetJson[User]("", name)
	if err != nil {
		fmt.Println("Error in retrieving user: ", err)
		return nil, err
	}

	return user, nil
}

func UpdateOrSetUser(uname string, toBan bool, txn *badger.Txn) error {
	user, err := FindUserByName(uname)

	// If user is not in the db, then add it
	if errors.Is(err, pebble.ErrNotFound) {
		u := new(User)
		u.Name = uname
		u.PubKey = ed25519.GenPrivKey().PubKey().Bytes()
		u.Banned = toBan
		user = u
	} else {
		if err == nil {
			user.Banned = toBan
		} else {
			err = fmt.Errorf("not able to process user")
			return err
		}
	}

	return SetJson[User]("", user.Name, user)
}
