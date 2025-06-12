package sidechain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cryptoencoding "github.com/cometbft/cometbft/crypto/encoding"
	"wetee.app/dsecret/internal/model"
)

func isBanTx(tx []byte) bool {
	return strings.Contains(string(tx), "username")
}

func (app *SideChain) getValidators() ([]types.ValidatorUpdate, error) {
	var err error
	validators, err := model.GetValidators()
	if err != nil {
		return nil, err
	}
	return validators, nil
}

func (app *SideChain) updateValidator(v types.ValidatorUpdate) error {
	pubKey, err := cryptoencoding.PubKeyFromTypeAndBytes(v.PubKeyType, v.PubKeyBytes)
	if err != nil {
		return fmt.Errorf("can't decode public key: %w", err)
	}
	key := "val" + string(pubKey.Bytes())

	// add or update validator
	value := bytes.NewBuffer(make([]byte, 0))
	if err := types.WriteMessage(&v, value); err != nil {
		return err
	}
	if err = model.Set(key, value.Bytes()); err != nil {
		return err
	}
	app.valAddrToPubKeyMap[string(pubKey.Address())] = pubKey
	return nil
}

func hasCurseWord(word string, curseWords string) bool {
	// Define your list of curse words here
	// For example:
	return strings.Contains(curseWords, word)
}

const (
	CodeTypeOK              uint32 = 0
	CodeTypeEncodingError   uint32 = 1
	CodeTypeInvalidTxFormat uint32 = 2
	CodeTypeBanned          uint32 = 3
)

func UpdateOrSetUser(uname string, toBan bool, txn *model.Txn) error {
	var u *model.User
	u, err := model.FindUserByName(uname)
	if errors.Is(err, pebble.ErrNotFound) {
		u = new(model.User)
		u.Name = uname
		u.PubKey = ed25519.GenPrivKey().PubKey().Bytes()
		u.Banned = toBan
	} else {
		if err == nil {
			u.Banned = toBan
		} else {
			err = fmt.Errorf("not able to process user")
			return err
		}
	}

	userBytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println("Error marshalling user")
		return err
	}

	return txn.Set([]byte(uname), userBytes)
}

func DeduplicateCurseWords(inWords string) string {
	curseWordMap := make(map[string]struct{})
	for _, word := range strings.Split(inWords, "|") {
		curseWordMap[word] = struct{}{}
	}
	deduplicatedWords := ""
	for word := range curseWordMap {
		if deduplicatedWords == "" {
			deduplicatedWords = word
		} else {
			deduplicatedWords = deduplicatedWords + "|" + word
		}
	}
	return deduplicatedWords
}
