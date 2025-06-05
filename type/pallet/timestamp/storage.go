package timestamp

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for Now id={{false [12]}}
//
//	The current time for the current block.
func MakeNowStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Timestamp", "Now")
}

var NowResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetNow(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNowStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NowResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNowLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNowStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NowResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for DidUpdate id={{false [8]}}
//
//	Whether the timestamp has been updated in this block.
//
//	This value is updated to `true` upon successful submission of a timestamp by a node.
//	It is then checked at the end of each block execution in the `on_finalize` hook.
func MakeDidUpdateStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Timestamp", "DidUpdate")
}

var DidUpdateResultDefaultBytes, _ = hex.DecodeString("00")

func GetDidUpdate(state state.State, bhash types.Hash) (ret bool, err error) {
	key, err := MakeDidUpdateStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DidUpdateResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetDidUpdateLatest(state state.State) (ret bool, err error) {
	key, err := MakeDidUpdateStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DidUpdateResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
