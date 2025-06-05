package gov

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// create a proposal
// 创建一个提案
func MakeSubmitProposalCall(daoId0 uint64, memberData1 types.MemberData, proposal2 types.RuntimeCall, periodIndex3 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsSubmitProposal:             true,
			AsSubmitProposalDaoId0:       daoId0,
			AsSubmitProposalMemberData1:  memberData1,
			AsSubmitProposalProposal2:    &proposal2,
			AsSubmitProposalPeriodIndex3: periodIndex3,
		},
	}
}

// Open a prop.
// 开始全民公投
func MakeDepositProposalCall(daoId0 uint64, proposeId1 uint32, deposit2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsDepositProposal:           true,
			AsDepositProposalDaoId0:     daoId0,
			AsDepositProposalProposeId1: proposeId1,
			AsDepositProposalDeposit2:   deposit2,
		},
	}
}

// Vote for the prop
// 为全民公投投票
func MakeVoteForPropCall(daoId0 uint64, propIndex1 uint32, pledge2 types1.UCompact, opinion3 types.Opinion) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsVoteForProp:           true,
			AsVoteForPropDaoId0:     daoId0,
			AsVoteForPropPropIndex1: propIndex1,
			AsVoteForPropPledge2:    pledge2,
			AsVoteForPropOpinion3:   opinion3,
		},
	}
}

// Cancel a vote on a prop
// 取消一个投票
func MakeCancelVoteCall(daoId0 uint64, index1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsCancelVote:       true,
			AsCancelVoteDaoId0: daoId0,
			AsCancelVoteIndex1: index1,
		},
	}
}

// Vote and execute the transaction corresponding to the proposa
// 执行一个投票通过提案
func MakeRunProposalCall(daoId0 uint64, index1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsRunProposal:       true,
			AsRunProposalDaoId0: daoId0,
			AsRunProposalIndex1: index1,
		},
	}
}

// Unlock
func MakeUnlockCall(daoId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsUnlock:       true,
			AsUnlockDaoId0: daoId0,
		},
	}
}

// Set the maximum number of proposals at the same time
func MakeSetMaxPrePropsCall(daoId0 uint64, max1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsSetMaxPreProps:       true,
			AsSetMaxPrePropsDaoId0: daoId0,
			AsSetMaxPrePropsMax1:   max1,
		},
	}
}
func MakeUpdateVoteModelCall(daoId0 uint64, model1 byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsUpdateVoteModel:       true,
			AsUpdateVoteModelDaoId0: daoId0,
			AsUpdateVoteModelModel1: model1,
		},
	}
}
func MakeSetPeriodsCall(daoId0 uint64, periods1 []types.Period) types.RuntimeCall {
	return types.RuntimeCall{
		IsGov: true,
		AsGovField0: &types.WeteeGovPalletCall{
			IsSetPeriods:         true,
			AsSetPeriodsDaoId0:   daoId0,
			AsSetPeriodsPeriods1: periods1,
		},
	}
}
