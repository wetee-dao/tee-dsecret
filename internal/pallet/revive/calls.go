package revive

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// A raw EVM transaction, typically dispatched by an Ethereum JSON-RPC server.
//
// # Parameters
//
//   - `payload`: The encoded [`crate::evm::TransactionSigned`].
//   - `gas_limit`: The gas limit enforced during contract execution.
//   - `storage_deposit_limit`: The maximum balance that can be charged to the caller for
//     storage usage.
//
// # Note
//
// This call cannot be dispatched directly; attempting to do so will result in a failed
// transaction. It serves as a wrapper for an Ethereum transaction. When submitted, the
// runtime converts it into a [`sp_runtime::generic::CheckedExtrinsic`] by recovering the
// signer and validating the transaction.
func MakeEthTransactCall(payload0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsEthTransact:         true,
			AsEthTransactPayload0: payload0,
		},
	}
}

// Makes a call to an account, optionally transferring some balance.
//
// # Parameters
//
//   - `dest`: Address of the contract to call.
//   - `value`: The balance to transfer from the `origin` to `dest`.
//   - `gas_limit`: The gas limit enforced when executing the constructor.
//   - `storage_deposit_limit`: The maximum amount of balance that can be charged from the
//     caller to pay for the storage consumed.
//   - `data`: The input data to pass to the contract.
//
// * If the account is a smart-contract account, the associated code will be
// executed and any value will be transferred.
// * If the account is a regular account, any value will be transferred.
// * If no account exists and the call value is not less than `existential_deposit`,
// a regular account will be created and any value will be transferred.
func MakeCallCall(dest0 [20]byte, value1 types1.UCompact, gasLimit2 types.Weight, storageDepositLimit3 types1.UCompact, data4 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsCall:                     true,
			AsCallDest0:                dest0,
			AsCallValue1:               value1,
			AsCallGasLimit2:            gasLimit2,
			AsCallStorageDepositLimit3: storageDepositLimit3,
			AsCallData4:                data4,
		},
	}
}

// Instantiates a contract from a previously deployed wasm binary.
//
// This function is identical to [`Self::instantiate_with_code`] but without the
// code deployment step. Instead, the `code_hash` of an on-chain deployed wasm binary
// must be supplied.
func MakeInstantiateCall(value0 types1.UCompact, gasLimit1 types.Weight, storageDepositLimit2 types1.UCompact, codeHash3 [32]byte, data4 []byte, salt5 types.OptionTByteArray321) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsInstantiate:                     true,
			AsInstantiateValue0:               value0,
			AsInstantiateGasLimit1:            gasLimit1,
			AsInstantiateStorageDepositLimit2: storageDepositLimit2,
			AsInstantiateCodeHash3:            codeHash3,
			AsInstantiateData4:                data4,
			AsInstantiateSalt5:                salt5,
		},
	}
}

// Instantiates a new contract from the supplied `code` optionally transferring
// some balance.
//
// This dispatchable has the same effect as calling [`Self::upload_code`] +
// [`Self::instantiate`]. Bundling them together provides efficiency gains. Please
// also check the documentation of [`Self::upload_code`].
//
// # Parameters
//
//   - `value`: The balance to transfer from the `origin` to the newly created contract.
//   - `gas_limit`: The gas limit enforced when executing the constructor.
//   - `storage_deposit_limit`: The maximum amount of balance that can be charged/reserved
//     from the caller to pay for the storage consumed.
//   - `code`: The contract code to deploy in raw bytes.
//   - `data`: The input data to pass to the contract constructor.
//   - `salt`: Used for the address derivation. If `Some` is supplied then `CREATE2`
//     semantics are used. If `None` then `CRATE1` is used.
//
// Instantiation is executed as follows:
//
// - The supplied `code` is deployed, and a `code_hash` is created for that code.
// - If the `code_hash` already exists on the chain the underlying `code` will be shared.
// - The destination address is computed based on the sender, code_hash and the salt.
// - The smart-contract account is created at the computed address.
// - The `value` is transferred to the new account.
// - The `deploy` function is executed in the context of the newly-created account.
func MakeInstantiateWithCodeCall(value0 types1.UCompact, gasLimit1 types.Weight, storageDepositLimit2 types1.UCompact, code3 []byte, data4 []byte, salt5 types.OptionTByteArray321) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsInstantiateWithCode:                     true,
			AsInstantiateWithCodeValue0:               value0,
			AsInstantiateWithCodeGasLimit1:            gasLimit1,
			AsInstantiateWithCodeStorageDepositLimit2: storageDepositLimit2,
			AsInstantiateWithCodeCode3:                code3,
			AsInstantiateWithCodeData4:                data4,
			AsInstantiateWithCodeSalt5:                salt5,
		},
	}
}

// Upload new `code` without instantiating a contract from it.
//
// If the code does not already exist a deposit is reserved from the caller
// and unreserved only when [`Self::remove_code`] is called. The size of the reserve
// depends on the size of the supplied `code`.
//
// # Note
//
// Anyone can instantiate a contract from any uploaded code and thus prevent its removal.
// To avoid this situation a constructor could employ access control so that it can
// only be instantiated by permissioned entities. The same is true when uploading
// through [`Self::instantiate_with_code`].
func MakeUploadCodeCall(code0 []byte, storageDepositLimit1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsUploadCode:                     true,
			AsUploadCodeCode0:                code0,
			AsUploadCodeStorageDepositLimit1: storageDepositLimit1,
		},
	}
}

// Remove the code stored under `code_hash` and refund the deposit to its owner.
//
// A code can only be removed by its original uploader (its owner) and only if it is
// not used by any contract.
func MakeRemoveCodeCall(codeHash0 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsRemoveCode:          true,
			AsRemoveCodeCodeHash0: codeHash0,
		},
	}
}

// Privileged function that changes the code of an existing contract.
//
// This takes care of updating refcounts and all other necessary operations. Returns
// an error if either the `code_hash` or `dest` do not exist.
//
// # Note
//
// This does **not** change the address of the contract in question. This means
// that the contract address is no longer derived from its code hash after calling
// this dispatchable.
func MakeSetCodeCall(dest0 [20]byte, codeHash1 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsSetCode:          true,
			AsSetCodeDest0:     dest0,
			AsSetCodeCodeHash1: codeHash1,
		},
	}
}

// Register the callers account id so that it can be used in contract interactions.
//
// This will error if the origin is already mapped or is a eth native `Address20`. It will
// take a deposit that can be released by calling [`Self::unmap_account`].
func MakeMapAccountCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsMapAccount: true,
		},
	}
}

// Unregister the callers account id in order to free the deposit.
//
// There is no reason to ever call this function other than freeing up the deposit.
// This is only useful when the account should no longer be used.
func MakeUnmapAccountCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsUnmapAccount: true,
		},
	}
}

// Dispatch an `call` with the origin set to the callers fallback address.
//
// Every `AccountId32` can control its corresponding fallback account. The fallback account
// is the `AccountId20` with the last 12 bytes set to `0xEE`. This is essentially a
// recovery function in case an `AccountId20` was used without creating a mapping first.
func MakeDispatchAsFallbackAccountCall(call0 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsRevive: true,
		AsReviveField0: &types.PalletRevivePalletCall{
			IsDispatchAsFallbackAccount:      true,
			AsDispatchAsFallbackAccountCall0: &call0,
		},
	}
}
