package task

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/tee-dsecret/chains/pallets/generated/types"
)

// Make a storage key for NextTeeId id={{false [12]}}
//
//	The id of the next app to be created.
//	获取下一个app id
func MakeNextTeeIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Task", "NextTeeId")
}

var NextTeeIdResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetNextTeeId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNextTeeIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextTeeIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextTeeIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNextTeeIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextTeeIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for TEETasks
//
//	Task
//	应用
func MakeTEETasksStorageKey(tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6410)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6411)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "TEETasks", byteArgs...)
}
func GetTEETasks(state state.State, bhash types.Hash, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret types1.TeeTask, isSome bool, err error) {
	key, err := MakeTEETasksStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTEETasksLatest(state state.State, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret types1.TeeTask, isSome bool, err error) {
	key, err := MakeTEETasksStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Prices
//
//	Price of resource
//	价格
func MakePricesStorageKey(byte0 byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byte0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "Prices", byteArgs...)
}
func GetPrices(state state.State, bhash types.Hash, byte0 byte) (ret types1.Price1, isSome bool, err error) {
	key, err := MakePricesStorageKey(byte0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetPricesLatest(state state.State, byte0 byte) (ret types1.Price1, isSome bool, err error) {
	key, err := MakePricesStorageKey(byte0)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for TaskIdAccounts
//
//	Task 对应账户
//	user's K8sCluster information
func MakeTaskIdAccountsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "TaskIdAccounts", byteArgs...)
}
func GetTaskIdAccounts(state state.State, bhash types.Hash, uint640 uint64) (ret [32]byte, isSome bool, err error) {
	key, err := MakeTaskIdAccountsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTaskIdAccountsLatest(state state.State, uint640 uint64) (ret [32]byte, isSome bool, err error) {
	key, err := MakeTaskIdAccountsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Envs
//
//	Task setting
//	Task设置
func MakeEnvsStorageKey(tupleOfUint64Uint160 uint64, tupleOfUint64Uint161 uint16) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64Uint160)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64Uint161)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "Envs", byteArgs...)
}
func GetEnvs(state state.State, bhash types.Hash, tupleOfUint64Uint160 uint64, tupleOfUint64Uint161 uint16) (ret types1.Env1, isSome bool, err error) {
	key, err := MakeEnvsStorageKey(tupleOfUint64Uint160, tupleOfUint64Uint161)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetEnvsLatest(state state.State, tupleOfUint64Uint160 uint64, tupleOfUint64Uint161 uint16) (ret types1.Env1, isSome bool, err error) {
	key, err := MakeEnvsStorageKey(tupleOfUint64Uint160, tupleOfUint64Uint161)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for SecretEnvs
//
//	Secret app setting
//	加密设置
func MakeSecretEnvsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "SecretEnvs", byteArgs...)
}
func GetSecretEnvs(state state.State, bhash types.Hash, uint640 uint64) (ret []byte, isSome bool, err error) {
	key, err := MakeSecretEnvsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetSecretEnvsLatest(state state.State, uint640 uint64) (ret []byte, isSome bool, err error) {
	key, err := MakeSecretEnvsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for CodeSignature
//
//	代码版本
func MakeCodeSignatureStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "CodeSignature", byteArgs...)
}

var CodeSignatureResultDefaultBytes, _ = hex.DecodeString("00")

func GetCodeSignature(state state.State, bhash types.Hash, uint640 uint64) (ret []byte, err error) {
	key, err := MakeCodeSignatureStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CodeSignatureResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetCodeSignatureLatest(state state.State, uint640 uint64) (ret []byte, err error) {
	key, err := MakeCodeSignatureStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CodeSignatureResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for CodeSigner
//
//	代码打包签名人
func MakeCodeSignerStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "CodeSigner", byteArgs...)
}

var CodeSignerResultDefaultBytes, _ = hex.DecodeString("00")

func GetCodeSigner(state state.State, bhash types.Hash, uint640 uint64) (ret []byte, err error) {
	key, err := MakeCodeSignerStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CodeSignerResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetCodeSignerLatest(state state.State, uint640 uint64) (ret []byte, err error) {
	key, err := MakeCodeSignerStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CodeSignerResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for TaskVersion
//
//	Task version
//	Task 版本
func MakeTaskVersionStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Task", "TaskVersion", byteArgs...)
}
func GetTaskVersion(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeTaskVersionStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTaskVersionLatest(state state.State, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeTaskVersionStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}
