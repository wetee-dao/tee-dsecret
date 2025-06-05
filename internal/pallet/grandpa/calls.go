package grandpa

import types "github.com/wetee-dao/go-sdk/pallet/types"

// Report voter equivocation/misbehavior. This method will verify the
// equivocation proof and validate the given key ownership proof
// against the extracted offender. If both are valid, the offence
// will be reported.
func MakeReportEquivocationCall(equivocationProof0 types.EquivocationProof, keyOwnerProof1 struct{}) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsReportEquivocation:                   true,
			AsReportEquivocationEquivocationProof0: &equivocationProof0,
			AsReportEquivocationKeyOwnerProof1:     keyOwnerProof1,
		},
	}
}

// Report voter equivocation/misbehavior. This method will verify the
// equivocation proof and validate the given key ownership proof
// against the extracted offender. If both are valid, the offence
// will be reported.
//
// This extrinsic must be called unsigned and it is expected that only
// block authors will call it (validated in `ValidateUnsigned`), as such
// if the block author is defined it will be defined as the equivocation
// reporter.
func MakeReportEquivocationUnsignedCall(equivocationProof0 types.EquivocationProof, keyOwnerProof1 struct{}) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsReportEquivocationUnsigned:                   true,
			AsReportEquivocationUnsignedEquivocationProof0: &equivocationProof0,
			AsReportEquivocationUnsignedKeyOwnerProof1:     keyOwnerProof1,
		},
	}
}

// Note that the current authority set of the GRANDPA finality gadget has stalled.
//
// This will trigger a forced authority set change at the beginning of the next session, to
// be enacted `delay` blocks after that. The `delay` should be high enough to safely assume
// that the block signalling the forced change will not be re-orged e.g. 1000 blocks.
// The block production rate (which may be slowed down because of finality lagging) should
// be taken into account when choosing the `delay`. The GRANDPA voters based on the new
// authority will start voting on top of `best_finalized_block_number` for new finalized
// blocks. `best_finalized_block_number` should be the highest of the latest finalized
// block of all validators of the new authority set.
//
// Only callable by root.
func MakeNoteStalledCall(delay0 uint32, bestFinalizedBlockNumber1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsNoteStalled:                          true,
			AsNoteStalledDelay0:                    delay0,
			AsNoteStalledBestFinalizedBlockNumber1: bestFinalizedBlockNumber1,
		},
	}
}
