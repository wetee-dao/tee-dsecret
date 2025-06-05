package timestamp

import (
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// Set the current time.
//
// This call should be invoked exactly once per block. It will panic at the finalization
// phase, if this call hasn't been invoked by that time.
//
// The timestamp should be greater than the previous one by the amount specified by
// [`Config::MinimumPeriod`].
//
// The dispatch origin for this call must be _None_.
//
// This dispatch class is _Mandatory_ to ensure it gets executed in the block. Be aware
// that changing the complexity of this call could result exhausting the resources in a
// block to execute any other calls.
//
// ## Complexity
//   - `O(1)` (Note that implementations of `OnTimestampSet` must also be `O(1)`)
//   - 1 storage read and 1 storage mutation (codec `O(1)` because of `DidUpdate::take` in
//     `on_finalize`)
//   - 1 event handler `on_timestamp_set`. Must be `O(1)`.
func MakeSetCall(now0 types.UCompact) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsTimestamp: true,
		AsTimestampField0: &types1.PalletTimestampPalletCall{
			IsSet:     true,
			AsSetNow0: now0,
		},
	}
}
