package authorship

import (
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Make a storage key for Author id={{false []}}
//
//	Author of current block.
func MakeAuthorStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Authorship", "Author")
}
func GetAuthor(state state.State, bhash types.Hash) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAuthorStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAuthorLatest(state state.State) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAuthorStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}
