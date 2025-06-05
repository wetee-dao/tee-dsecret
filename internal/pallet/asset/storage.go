package asset

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for ChainID id={{false [4]}}
func MakeChainIDStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Asset", "ChainID")
}

var ChainIDResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetChainID(state state.State, bhash types.Hash) (ret uint32, err error) {
	key, err := MakeChainIDStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ChainIDResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetChainIDLatest(state state.State) (ret uint32, err error) {
	key, err := MakeChainIDStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ChainIDResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for AssetInfos
func MakeAssetInfosStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "AssetInfos", byteArgs...)
}
func GetAssetInfos(state state.State, bhash types.Hash, uint640 uint64) (ret types1.AssetInfo, isSome bool, err error) {
	key, err := MakeAssetInfosStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAssetInfosLatest(state state.State, uint640 uint64) (ret types1.AssetInfo, isSome bool, err error) {
	key, err := MakeAssetInfosStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for AssetSymbols
func MakeAssetSymbolsStorageKey(byteSlice0 []byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteSlice0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "AssetSymbols", byteArgs...)
}
func GetAssetSymbols(state state.State, bhash types.Hash, byteSlice0 []byte) (ret uint64, isSome bool, err error) {
	key, err := MakeAssetSymbolsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAssetSymbolsLatest(state state.State, byteSlice0 []byte) (ret uint64, isSome bool, err error) {
	key, err := MakeAssetSymbolsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for ParaAssetMaps
func MakeParaAssetMapsStorageKey(tupleOfUint32ByteSlice0 uint32, tupleOfUint32ByteSlice1 []byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint32ByteSlice0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint32ByteSlice1)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "ParaAssetMaps", byteArgs...)
}
func GetParaAssetMaps(state state.State, bhash types.Hash, tupleOfUint32ByteSlice0 uint32, tupleOfUint32ByteSlice1 []byte) (ret uint64, isSome bool, err error) {
	key, err := MakeParaAssetMapsStorageKey(tupleOfUint32ByteSlice0, tupleOfUint32ByteSlice1)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetParaAssetMapsLatest(state state.State, tupleOfUint32ByteSlice0 uint32, tupleOfUint32ByteSlice1 []byte) (ret uint64, isSome bool, err error) {
	key, err := MakeParaAssetMapsStorageKey(tupleOfUint32ByteSlice0, tupleOfUint32ByteSlice1)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for AssetParaIds
func MakeAssetParaIdsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "AssetParaIds", byteArgs...)
}
func GetAssetParaIds(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeAssetParaIdsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAssetParaIdsLatest(state state.State, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeAssetParaIdsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for LocationToSymbols
func MakeLocationToSymbolsStorageKey(byteSlice0 []byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteSlice0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "LocationToSymbols", byteArgs...)
}
func GetLocationToSymbols(state state.State, bhash types.Hash, byteSlice0 []byte) (ret []byte, isSome bool, err error) {
	key, err := MakeLocationToSymbolsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetLocationToSymbolsLatest(state state.State, byteSlice0 []byte) (ret []byte, isSome bool, err error) {
	key, err := MakeLocationToSymbolsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for SymbolToLocations
func MakeSymbolToLocationsStorageKey(byteSlice0 []byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteSlice0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Asset", "SymbolToLocations", byteArgs...)
}
func GetSymbolToLocations(state state.State, bhash types.Hash, byteSlice0 []byte) (ret []byte, isSome bool, err error) {
	key, err := MakeSymbolToLocationsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetSymbolToLocationsLatest(state state.State, byteSlice0 []byte) (ret []byte, isSome bool, err error) {
	key, err := MakeSymbolToLocationsStorageKey(byteSlice0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}
