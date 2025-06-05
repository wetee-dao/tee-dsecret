package balances

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// Transfer some liquid free balance to another account.
//
// `transfer_allow_death` will set the `FreeBalance` of the sender and receiver.
// If the sender's account is below the existential deposit as a result
// of the transfer, the account will be reaped.
//
// The dispatch origin for this call must be `Signed` by the transactor.
func MakeTransferAllowDeathCall(dest0 types.MultiAddress, value1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferAllowDeath:       true,
			AsTransferAllowDeathDest0:  dest0,
			AsTransferAllowDeathValue1: value1,
		},
	}
}

// Exactly as `transfer_allow_death`, except the origin must be root and the source account
// may be specified.
func MakeForceTransferCall(source0 types.MultiAddress, dest1 types.MultiAddress, value2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceTransfer:        true,
			AsForceTransferSource0: source0,
			AsForceTransferDest1:   dest1,
			AsForceTransferValue2:  value2,
		},
	}
}

// Same as the [`transfer_allow_death`] call, but with a check that the transfer will not
// kill the origin account.
//
// 99% of the time you want [`transfer_allow_death`] instead.
//
// [`transfer_allow_death`]: struct.Pallet.html#method.transfer
func MakeTransferKeepAliveCall(dest0 types.MultiAddress, value1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferKeepAlive:       true,
			AsTransferKeepAliveDest0:  dest0,
			AsTransferKeepAliveValue1: value1,
		},
	}
}

// Transfer the entire transferable balance from the caller account.
//
// NOTE: This function only attempts to transfer _transferable_ balances. This means that
// any locked, reserved, or existential deposits (when `keep_alive` is `true`), will not be
// transferred by this function. To ensure that this function results in a killed account,
// you might need to prepare the account by removing any reference counters, storage
// deposits, etc...
//
// The dispatch origin of this call must be Signed.
//
//   - `dest`: The recipient of the transfer.
//   - `keep_alive`: A boolean to determine if the `transfer_all` operation should send all
//     of the funds the account has, causing the sender account to be killed (false), or
//     transfer everything except at least the existential deposit, which will guarantee to
//     keep the sender account alive (true).
func MakeTransferAllCall(dest0 types.MultiAddress, keepAlive1 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferAll:           true,
			AsTransferAllDest0:      dest0,
			AsTransferAllKeepAlive1: keepAlive1,
		},
	}
}

// Unreserve some balance from a user by force.
//
// Can only be called by ROOT.
func MakeForceUnreserveCall(who0 types.MultiAddress, amount1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceUnreserve:        true,
			AsForceUnreserveWho0:    who0,
			AsForceUnreserveAmount1: amount1,
		},
	}
}

// Upgrade a specified account.
//
// - `origin`: Must be `Signed`.
// - `who`: The account to be upgraded.
//
// This will waive the transaction fee if at least all but 10% of the accounts needed to
// be upgraded. (We let some not have to be upgraded just in order to allow for the
// possibility of churn).
func MakeUpgradeAccountsCall(who0 [][32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsUpgradeAccounts:     true,
			AsUpgradeAccountsWho0: who0,
		},
	}
}

// Set the regular balance of a given account.
//
// The dispatch origin for this call is `root`.
func MakeForceSetBalanceCall(who0 types.MultiAddress, newFree1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceSetBalance:         true,
			AsForceSetBalanceWho0:     who0,
			AsForceSetBalanceNewFree1: newFree1,
		},
	}
}

// Adjust the total issuance in a saturating way.
//
// Can only be called by root and always needs a positive `delta`.
//
// # Example
func MakeForceAdjustTotalIssuanceCall(direction0 types.AdjustmentDirection, delta1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceAdjustTotalIssuance:           true,
			AsForceAdjustTotalIssuanceDirection0: direction0,
			AsForceAdjustTotalIssuanceDelta1:     delta1,
		},
	}
}

// Burn the specified liquid free balance from the origin account.
//
// If the origin's account ends up below the existential deposit as a result
// of the burn and `keep_alive` is false, the account will be reaped.
//
// Unlike sending funds to a _burn_ address, which merely makes the funds inaccessible,
// this `burn` operation will reduce total issuance by the amount _burned_.
func MakeBurnCall(value0 types1.UCompact, keepAlive1 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsBurn:           true,
			AsBurnValue0:     value0,
			AsBurnKeepAlive1: keepAlive1,
		},
	}
}
