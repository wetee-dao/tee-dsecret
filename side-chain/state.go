package sidechain

import (
	"encoding/binary"
	"errors"

	"github.com/cockroachdb/pebble"
	"wetee.app/dsecret/internal/model"
)

type AppState struct {
	Size   int64 `json:"size"`
	Height int64 `json:"height"`
}

var stateKey = "appstate"

func (s AppState) Hash() []byte {
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, s.Size)
	return appHash
}

func loadState() (AppState, error) {
	state, err := model.GetJson[AppState]("", stateKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return AppState{}, nil
	}

	return *state, nil
}

func saveState(state *AppState) error {
	return model.SetJson("", stateKey, state)
}
