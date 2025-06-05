package tokens

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// Transfer some liquid free balance to another account.
//
// `transfer` will set the `FreeBalance` of the sender and receiver.
// It will decrease the total issuance of the system by the
// `TransferFee`. If the sender's account is below the existential
// deposit as a result of the transfer, the account will be reaped.
//
// The dispatch origin for this call must be `Signed` by the
// transactor.
//
// - `dest`: The recipient of the transfer.
// - `currency_id`: currency type.
// - `amount`: free balance amount to transfer.
func MakeTransferCall(dest0 types.MultiAddress, currencyId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransfer:            true,
			AsTransferDest0:       dest0,
			AsTransferCurrencyId1: currencyId1,
			AsTransferAmount2:     amount2,
		},
	}
}

// Transfer all remaining balance to the given account.
//
// NOTE: This function only attempts to transfer _transferable_
// balances. This means that any locked, reserved, or existential
// deposits (when `keep_alive` is `true`), will not be transferred by
// this function. To ensure that this function results in a killed
// account, you might need to prepare the account by removing any
// reference counters, storage deposits, etc...
//
// The dispatch origin for this call must be `Signed` by the
// transactor.
//
//   - `dest`: The recipient of the transfer.
//   - `currency_id`: currency type.
//   - `keep_alive`: A boolean to determine if the `transfer_all`
//     operation should send all of the funds the account has, causing
//     the sender account to be killed (false), or transfer everything
//     except at least the existential deposit, which will guarantee to
//     keep the sender account alive (true).
func MakeTransferAllCall(dest0 types.MultiAddress, currencyId1 uint64, keepAlive2 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransferAll:            true,
			AsTransferAllDest0:       dest0,
			AsTransferAllCurrencyId1: currencyId1,
			AsTransferAllKeepAlive2:  keepAlive2,
		},
	}
}

// Same as the [`transfer`] call, but with a check that the transfer
// will not kill the origin account.
//
// 99% of the time you want [`transfer`] instead.
//
// The dispatch origin for this call must be `Signed` by the
// transactor.
//
// - `dest`: The recipient of the transfer.
// - `currency_id`: currency type.
// - `amount`: free balance amount to transfer.
func MakeTransferKeepAliveCall(dest0 types.MultiAddress, currencyId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransferKeepAlive:            true,
			AsTransferKeepAliveDest0:       dest0,
			AsTransferKeepAliveCurrencyId1: currencyId1,
			AsTransferKeepAliveAmount2:     amount2,
		},
	}
}

// Exactly as `transfer`, except the origin must be root and the source
// account may be specified.
//
// The dispatch origin for this call must be _Root_.
//
// - `source`: The sender of the transfer.
// - `dest`: The recipient of the transfer.
// - `currency_id`: currency type.
// - `amount`: free balance amount to transfer.
func MakeForceTransferCall(source0 types.MultiAddress, dest1 types.MultiAddress, currencyId2 uint64, amount3 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsForceTransfer:            true,
			AsForceTransferSource0:     source0,
			AsForceTransferDest1:       dest1,
			AsForceTransferCurrencyId2: currencyId2,
			AsForceTransferAmount3:     amount3,
		},
	}
}

// Set the balances of a given account.
//
// This will alter `FreeBalance` and `ReservedBalance` in storage. it
// will also decrease the total issuance of the system
// (`TotalIssuance`). If the new free or reserved balance is below the
// existential deposit, it will reap the `AccountInfo`.
//
// The dispatch origin for this call is `root`.
func MakeSetBalanceCall(who0 types.MultiAddress, currencyId1 uint64, newFree2 types1.UCompact, newReserved3 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsSetBalance:             true,
			AsSetBalanceWho0:         who0,
			AsSetBalanceCurrencyId1:  currencyId1,
			AsSetBalanceNewFree2:     newFree2,
			AsSetBalanceNewReserved3: newReserved3,
		},
	}
}
