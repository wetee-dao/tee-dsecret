package sidechain

import (
	"bytes"
	"fmt"

	"github.com/cometbft/cometbft/abci/types"
	cryptoencoding "github.com/cometbft/cometbft/crypto/encoding"
	"wetee.app/dsecret/internal/model"
)

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

const (
	CodeTypeOK              uint32 = 0
	CodeTypeEncodingError   uint32 = 1
	CodeTypeInvalidTxFormat uint32 = 2
	CodeTypeBanned          uint32 = 3
)
