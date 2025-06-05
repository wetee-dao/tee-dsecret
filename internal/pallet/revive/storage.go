package revive

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for PristineCode
//
//	A mapping from a contract's code hash to its code.
func MakePristineCodeStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "PristineCode", byteArgs...)
}
func GetPristineCode(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []byte, isSome bool, err error) {
	key, err := MakePristineCodeStorageKey(byteArray320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetPristineCodeLatest(state state.State, byteArray320 [32]byte) (ret []byte, isSome bool, err error) {
	key, err := MakePristineCodeStorageKey(byteArray320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for CodeInfoOf
//
//	A mapping from a contract's code hash to its code info.
func MakeCodeInfoOfStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "CodeInfoOf", byteArgs...)
}
func GetCodeInfoOf(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret types1.CodeInfo, isSome bool, err error) {
	key, err := MakeCodeInfoOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetCodeInfoOfLatest(state state.State, byteArray320 [32]byte) (ret types1.CodeInfo, isSome bool, err error) {
	key, err := MakeCodeInfoOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for ContractInfoOf
//
//	The code associated with a given account.
func MakeContractInfoOfStorageKey(byteArray200 [20]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray200)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "ContractInfoOf", byteArgs...)
}
func GetContractInfoOf(state state.State, bhash types.Hash, byteArray200 [20]byte) (ret types1.ContractInfo, isSome bool, err error) {
	key, err := MakeContractInfoOfStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetContractInfoOfLatest(state state.State, byteArray200 [20]byte) (ret types1.ContractInfo, isSome bool, err error) {
	key, err := MakeContractInfoOfStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for ImmutableDataOf
//
//	The immutable data associated with a given account.
func MakeImmutableDataOfStorageKey(byteArray200 [20]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray200)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "ImmutableDataOf", byteArgs...)
}
func GetImmutableDataOf(state state.State, bhash types.Hash, byteArray200 [20]byte) (ret []byte, isSome bool, err error) {
	key, err := MakeImmutableDataOfStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetImmutableDataOfLatest(state state.State, byteArray200 [20]byte) (ret []byte, isSome bool, err error) {
	key, err := MakeImmutableDataOfStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for DeletionQueue
//
//	Evicted contracts that await child trie deletion.
//
//	Child trie deletion is a heavy operation depending on the amount of storage items
//	stored in said trie. Therefore this operation is performed lazily in `on_idle`.
func MakeDeletionQueueStorageKey(uint320 uint32) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "DeletionQueue", byteArgs...)
}
func GetDeletionQueue(state state.State, bhash types.Hash, uint320 uint32) (ret []byte, isSome bool, err error) {
	key, err := MakeDeletionQueueStorageKey(uint320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetDeletionQueueLatest(state state.State, uint320 uint32) (ret []byte, isSome bool, err error) {
	key, err := MakeDeletionQueueStorageKey(uint320)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for DeletionQueueCounter id={{false [234]}}
//
//	A pair of monotonic counters used to track the latest contract marked for deletion
//	and the latest deleted contract in queue.
func MakeDeletionQueueCounterStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Revive", "DeletionQueueCounter")
}

var DeletionQueueCounterResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetDeletionQueueCounter(state state.State, bhash types.Hash) (ret types1.DeletionQueueManager, err error) {
	key, err := MakeDeletionQueueCounterStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DeletionQueueCounterResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetDeletionQueueCounterLatest(state state.State) (ret types1.DeletionQueueManager, err error) {
	key, err := MakeDeletionQueueCounterStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DeletionQueueCounterResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for OriginalAccount
//
//	Map a Ethereum address to its original `AccountId32`.
//
//	When deriving a `H160` from an `AccountId32` we use a hash function. In order to
//	reconstruct the original account we need to store the reverse mapping here.
//	Register your `AccountId32` using [`Pallet::map_account`] in order to
//	use it with this pallet.
func MakeOriginalAccountStorageKey(byteArray200 [20]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray200)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Revive", "OriginalAccount", byteArgs...)
}
func GetOriginalAccount(state state.State, bhash types.Hash, byteArray200 [20]byte) (ret [32]byte, isSome bool, err error) {
	key, err := MakeOriginalAccountStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetOriginalAccountLatest(state state.State, byteArray200 [20]byte) (ret [32]byte, isSome bool, err error) {
	key, err := MakeOriginalAccountStorageKey(byteArray200)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}
