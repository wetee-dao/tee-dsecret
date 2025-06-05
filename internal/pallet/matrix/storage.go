package matrix

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for Matrix
//
//	All Nodes that have been created.
//	所有节点
func MakeMatrixStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Matrix", "Matrix", byteArgs...)
}
func GetMatrix(state state.State, bhash types.Hash, u1280 types.U128) (ret types1.NodeInfo, isSome bool, err error) {
	key, err := MakeMatrixStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetMatrixLatest(state state.State, u1280 types.U128) (ret types1.NodeInfo, isSome bool, err error) {
	key, err := MakeMatrixStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for NextId id={{false [6]}}
//
//	The id of the next node to be created.
//	获取下一个节点id
func MakeNextIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Matrix", "NextId")
}

var NextIdResultDefaultBytes, _ = hex.DecodeString("01000000000000000000000000000000")

func GetNextId(state state.State, bhash types.Hash) (ret types.U128, err error) {
	key, err := MakeNextIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextIdLatest(state state.State) (ret types.U128, err error) {
	key, err := MakeNextIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
