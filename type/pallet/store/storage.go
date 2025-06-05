package store

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for NextAppId id={{false [6]}}
//
//	获取下一个应用 id
func MakeNextAppIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Store", "NextAppId")
}

var NextAppIdResultDefaultBytes, _ = hex.DecodeString("01000000000000000000000000000000")

func GetNextAppId(state state.State, bhash types.Hash) (ret types.U128, err error) {
	key, err := MakeNextAppIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextAppIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextAppIdLatest(state state.State) (ret types.U128, err error) {
	key, err := MakeNextAppIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextAppIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for AccountOfApp
//
//	用户对应集群的信息
func MakeAccountOfAppStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Store", "AccountOfApp", byteArgs...)
}
func GetAccountOfApp(state state.State, bhash types.Hash, u1280 types.U128) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAccountOfAppStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAccountOfAppLatest(state state.State, u1280 types.U128) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAccountOfAppStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for AccountApps
//
//	获取用户应用列表
func MakeAccountAppsStorageKey(tupleOfByteArray32U1280 [32]byte, tupleOfByteArray32U1281 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfByteArray32U1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfByteArray32U1281)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Store", "AccountApps", byteArgs...)
}

var AccountAppsResultDefaultBytes, _ = hex.DecodeString("")

func GetAccountApps(state state.State, bhash types.Hash, tupleOfByteArray32U1280 [32]byte, tupleOfByteArray32U1281 types.U128) (ret struct{}, err error) {
	key, err := MakeAccountAppsStorageKey(tupleOfByteArray32U1280, tupleOfByteArray32U1281)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountAppsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetAccountAppsLatest(state state.State, tupleOfByteArray32U1280 [32]byte, tupleOfByteArray32U1281 types.U128) (ret struct{}, err error) {
	key, err := MakeAccountAppsStorageKey(tupleOfByteArray32U1280, tupleOfByteArray32U1281)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountAppsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Apps
//
//	应用信息
func MakeAppsStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Store", "Apps", byteArgs...)
}
func GetApps(state state.State, bhash types.Hash, u1280 types.U128) (ret types1.AppTemplate, isSome bool, err error) {
	key, err := MakeAppsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAppsLatest(state state.State, u1280 types.U128) (ret types1.AppTemplate, isSome bool, err error) {
	key, err := MakeAppsStorageKey(u1280)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for VersionLists
//
//	应用版本
func MakeVersionListsStorageKey(tupleOfU128Uint160 types.U128, tupleOfU128Uint161 uint16) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfU128Uint160)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfU128Uint161)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Store", "VersionLists", byteArgs...)
}
func GetVersionLists(state state.State, bhash types.Hash, tupleOfU128Uint160 types.U128, tupleOfU128Uint161 uint16) (ret types1.TupleOfImageSliceUint32, isSome bool, err error) {
	key, err := MakeVersionListsStorageKey(tupleOfU128Uint160, tupleOfU128Uint161)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetVersionListsLatest(state state.State, tupleOfU128Uint160 types.U128, tupleOfU128Uint161 uint16) (ret types1.TupleOfImageSliceUint32, isSome bool, err error) {
	key, err := MakeVersionListsStorageKey(tupleOfU128Uint160, tupleOfU128Uint161)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for AppStakings
//
//	应用点数
func MakeAppStakingsStorageKey(u1280 types.U128) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(u1280)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Store", "AppStakings", byteArgs...)
}

var AppStakingsResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000")

func GetAppStakings(state state.State, bhash types.Hash, u1280 types.U128) (ret types.U128, err error) {
	key, err := MakeAppStakingsStorageKey(u1280)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AppStakingsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetAppStakingsLatest(state state.State, u1280 types.U128) (ret types.U128, err error) {
	key, err := MakeAppStakingsStorageKey(u1280)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AppStakingsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
