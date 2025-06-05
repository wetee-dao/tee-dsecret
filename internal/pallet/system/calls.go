package system

import types "github.com/wetee-dao/go-sdk/pallet/types"

// Make some on-chain remark.
//
// Can be executed by every `origin`.
func MakeRemarkCall(remark0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsRemark:        true,
			AsRemarkRemark0: remark0,
		},
	}
}

// Set the number of pages in the WebAssembly environment's heap.
func MakeSetHeapPagesCall(pages0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetHeapPages:       true,
			AsSetHeapPagesPages0: pages0,
		},
	}
}

// Set the new runtime code.
func MakeSetCodeCall(code0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetCode:      true,
			AsSetCodeCode0: code0,
		},
	}
}

// Set the new runtime code without doing any checks of the given `code`.
//
// Note that runtime upgrades will not run if this is called with a not-increasing spec
// version!
func MakeSetCodeWithoutChecksCall(code0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetCodeWithoutChecks:      true,
			AsSetCodeWithoutChecksCode0: code0,
		},
	}
}

// Set some items of storage.
func MakeSetStorageCall(items0 []types.TupleOfByteSliceByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetStorage:       true,
			AsSetStorageItems0: items0,
		},
	}
}

// Kill some items from storage.
func MakeKillStorageCall(keys0 [][]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsKillStorage:      true,
			AsKillStorageKeys0: keys0,
		},
	}
}

// Kill all storage items with a key that starts with the given prefix.
//
// **NOTE:** We rely on the Root origin to provide us the number of subkeys under
// the prefix we are removing to accurately calculate the weight of this function.
func MakeKillPrefixCall(prefix0 []byte, subkeys1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsKillPrefix:         true,
			AsKillPrefixPrefix0:  prefix0,
			AsKillPrefixSubkeys1: subkeys1,
		},
	}
}

// Make some on-chain remark and emit event.
func MakeRemarkWithEventCall(remark0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsRemarkWithEvent:        true,
			AsRemarkWithEventRemark0: remark0,
		},
	}
}

// Authorize an upgrade to a given `code_hash` for the runtime. The runtime can be supplied
// later.
//
// This call requires Root origin.
func MakeAuthorizeUpgradeCall(codeHash0 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsAuthorizeUpgrade:          true,
			AsAuthorizeUpgradeCodeHash0: codeHash0,
		},
	}
}

// Authorize an upgrade to a given `code_hash` for the runtime. The runtime can be supplied
// later.
//
// WARNING: This authorizes an upgrade that will take place without any safety checks, for
// example that the spec name remains the same and that the version number increases. Not
// recommended for normal use. Use `authorize_upgrade` instead.
//
// This call requires Root origin.
func MakeAuthorizeUpgradeWithoutChecksCall(codeHash0 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsAuthorizeUpgradeWithoutChecks:          true,
			AsAuthorizeUpgradeWithoutChecksCodeHash0: codeHash0,
		},
	}
}

// Provide the preimage (runtime binary) `code` for an upgrade that has been authorized.
//
// If the authorization required a version check, this call will ensure the spec name
// remains unchanged and that the spec version has increased.
//
// Depending on the runtime's `OnSetCode` configuration, this function may directly apply
// the new `code` in the same block or attempt to schedule the upgrade.
//
// All origins are allowed.
func MakeApplyAuthorizedUpgradeCall(code0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsApplyAuthorizedUpgrade:      true,
			AsApplyAuthorizedUpgradeCode0: code0,
		},
	}
}
