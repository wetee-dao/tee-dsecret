package model

import (
	"github.com/cometbft/cometbft/abci/types"
)

func GetValidators() (validators []types.ValidatorUpdate, err error) {
	list, err := GetAbciMessageList[*types.ValidatorUpdate]("", "val")
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		validators = append(validators, *item)
	}

	return
}
