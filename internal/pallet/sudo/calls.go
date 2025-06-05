package sudo

import types "github.com/wetee-dao/go-sdk/pallet/types"

// Authenticates the sudo key and dispatches a function call with `Root` origin.
func MakeSudoCall(call0 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudo:      true,
			AsSudoCall0: &call0,
		},
	}
}

// Authenticates the sudo key and dispatches a function call with `Root` origin.
// This function does not check the weight of the call, and instead allows the
// Sudo user to specify the weight of the call.
//
// The dispatch origin for this call must be _Signed_.
func MakeSudoUncheckedWeightCall(call0 types.RuntimeCall, weight1 types.Weight) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudoUncheckedWeight:        true,
			AsSudoUncheckedWeightCall0:   &call0,
			AsSudoUncheckedWeightWeight1: weight1,
		},
	}
}

// Authenticates the current sudo key and sets the given AccountId (`new`) as the new sudo
// key.
func MakeSetKeyCall(new0 types.MultiAddress) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSetKey:     true,
			AsSetKeyNew0: &new0,
		},
	}
}

// Authenticates the sudo key and dispatches a function call with `Signed` origin from
// a given account.
//
// The dispatch origin for this call must be _Signed_.
func MakeSudoAsCall(who0 types.MultiAddress, call1 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudoAs:      true,
			AsSudoAsWho0:  who0,
			AsSudoAsCall1: &call1,
		},
	}
}

// Permanently removes the sudo key.
//
// **This cannot be un-done.**
func MakeRemoveKeyCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsRemoveKey: true,
		},
	}
}
