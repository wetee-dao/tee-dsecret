package bridge

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for NextId id={{false [6]}}
func MakeNextIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Bridge", "NextId")
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

// Make a storage key for TEECalls
func MakeTEECallsStorageKey(tupleOfUint64U1280 uint64, tupleOfUint64U1281 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64U1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64U1281)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Bridge", "TEECalls", byteArgs...)
}
func GetTEECalls(state state.State, bhash types.Hash, tupleOfUint64U1280 uint64, tupleOfUint64U1281 types.U128) (ret types1.TEECall, isSome bool, err error) {
	key, err := MakeTEECallsStorageKey(tupleOfUint64U1280, tupleOfUint64U1281)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTEECallsLatest(state state.State, tupleOfUint64U1280 uint64, tupleOfUint64U1281 types.U128) (ret types1.TEECall, isSome bool, err error) {
	key, err := MakeTEECallsStorageKey(tupleOfUint64U1280, tupleOfUint64U1281)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for ApiMetas
//
//	App
//	应用
func MakeApiMetasStorageKey(workId0 types1.WorkId) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(workId0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Bridge", "ApiMetas", byteArgs...)
}
func GetApiMetas(state state.State, bhash types.Hash, workId0 types1.WorkId) (ret types1.ApiMeta, isSome bool, err error) {
	key, err := MakeApiMetasStorageKey(workId0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetApiMetasLatest(state state.State, workId0 types1.WorkId) (ret types1.ApiMeta, isSome bool, err error) {
	key, err := MakeApiMetasStorageKey(workId0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Results
func MakeResultsStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Bridge", "Results", byteArgs...)
}
func GetResults(state state.State, bhash types.Hash, u1280 types.U128) (ret types1.ContractResult, isSome bool, err error) {
	key, err := MakeResultsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetResultsLatest(state state.State, u1280 types.U128) (ret types1.ContractResult, isSome bool, err error) {
	key, err := MakeResultsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Enums
func MakeEnumsStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Bridge", "Enums", byteArgs...)
}
func GetEnums(state state.State, bhash types.Hash, u1280 types.U128) (ret types1.Test, isSome bool, err error) {
	key, err := MakeEnumsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetEnumsLatest(state state.State, u1280 types.U128) (ret types1.Test, isSome bool, err error) {
	key, err := MakeEnumsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Tuples
func MakeTuplesStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Bridge", "Tuples", byteArgs...)
}
func GetTuples(state state.State, bhash types.Hash, u1280 types.U128) (ret types1.TupleOfUint32Uint32, isSome bool, err error) {
	key, err := MakeTuplesStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTuplesLatest(state state.State, u1280 types.U128) (ret types1.TupleOfUint32Uint32, isSome bool, err error) {
	key, err := MakeTuplesStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}
