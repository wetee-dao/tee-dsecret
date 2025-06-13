package sidechain

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

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
	var state AppState

	stateBytes, err := model.GetKey("", stateKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return state, nil
	}
	if len(stateBytes) == 0 {
		return state, nil
	}
	err = json.Unmarshal(stateBytes, &state)
	fmt.Println("ST:", state)

	if err != nil {
		return state, err
	}

	return state, nil
}

func saveState(state *AppState) error {
	return model.SetJson("", stateKey, state)
}
