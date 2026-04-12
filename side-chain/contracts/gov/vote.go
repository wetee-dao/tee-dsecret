package gov

import (
	"bytes"
	"math/big"

	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Vote 获取单个投票
func (d GovQuery) Vote(proposalID uint32, index uint32) (util.Option[Vote], error) {
	v, err := d.votes.Get(d.api.GetTxn(), proposalID, index)
	if err != nil {
		return util.NewNone[Vote](), err
	}
	if v == nil {
		return util.NewNone[Vote](), ErrInvalidVote
	}
	return util.NewSome(*v), nil
}

// Votes 获取提案的所有投票（分页）
func (d GovQuery) Votes(proposalID uint32, startKey util.Option[uint32], size uint32) ([]Vote, error) {
	var startKeyPtr *uint32
	if startKey.IsSome() {
		startKeyPtr = &startKey.V
	}
	_, votes, err := d.votes.DescList(d.api.GetTxn(), proposalID, startKeyPtr, size)
	if err != nil {
		return nil, err
	}
	return votes, nil
}

// SubmitVote 提交投票
func (d GovMutation) SubmitVote(proposalID uint32, opinionYes bool, lockAmount model.Amount) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	member, err := d.members.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if member.Int.Sign() == 0 {
		return ErrMemberNotExisted
	}

	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalStatusOngoing {
		return ErrPropNotOngoing
	}
	track, err := d.track(prop.TrackID)
	if err != nil {
		return err
	}
	if height > prop.Deposit.Block+int64(track.MaxDeciding) {
		return ErrInvalidVoteTime
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	free := model.AmountSub(member, lock)
	if free.Int.Cmp(lockAmount.Int) < 0 {
		return ErrLowBalance
	}

	unlockAfter := int64(track.MinEnactmentPeriod)

	one := big.NewInt(1)
	vote := Vote{
		ProposalID:  proposalID,
		Caller:      caller,
		Pledge:      lockAmount,
		OpinionYes:  opinionYes,
		VoteWeight:  model.Amount{Int: one},
		UnlockBlock: unlockAfter,
		VoteBlock:   height,
	}

	// 使用 StoreList2D.Insert 插入投票
	index, err := d.votes.Insert(d.api.GetTxn(), proposalID, vote)
	if err != nil {
		return err
	}

	// 更新投票的 Index 字段
	vote.Index = index
	if err := d.votes.Update(d.api.GetTxn(), proposalID, index, vote); err != nil {
		return err
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountAdd(lock, lockAmount)); err != nil {
		return err
	}

	return nil
}

// CancelVote 取消投票
func (d GovMutation) CancelVote(proposalID uint32, index uint32) error {
	caller := d.api.GetCaller()
	vote, err := d.vote(proposalID, index)
	if err != nil {
		return err
	}
	if !bytes.Equal(vote.Caller.V, caller.V) {
		return ErrInvalidVoteUser
	}
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalStatusOngoing {
		return ErrPropNotOngoing
	}
	vote.Deleted = true

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if lock.Int.Cmp(vote.Pledge.Int) < 0 {
		return ErrLowBalance
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountSub(lock, vote.Pledge)); err != nil {
		return err
	}

	return d.votes.Update(d.api.GetTxn(), proposalID, index, vote)
}

// Unlock 解锁投票
func (d GovMutation) Unlock(proposalID uint32, index uint32) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	// 使用组合 key (proposalID, index) 检查解锁状态
	unlockKey := voteUnlockKey(proposalID, index)
	unlocked, err := d.voteUnlocks.Get(d.api.GetTxn(), unlockKey)
	if err != nil {
		return err
	}
	if unlocked != nil && *unlocked {
		return ErrVoteAlreadyUnlocked
	}

	vote, err := d.vote(proposalID, index)
	if err != nil {
		return err
	}
	if vote.Deleted {
		return ErrInvalidVoteStatus
	}
	if !bytes.Equal(vote.Caller.V, caller.V) {
		return ErrInvalidVoteUser
	}

	end, err := d.calculateProposalEndBlock(proposalID)
	if err != nil {
		return err
	}
	if height < end+vote.UnlockBlock {
		return ErrInvalidVoteUnlockTime
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if lock.Int.Cmp(vote.Pledge.Int) < 0 {
		return ErrLowBalance
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountSub(lock, vote.Pledge)); err != nil {
		return err
	}

	return d.voteUnlocks.Set(d.api.GetTxn(), unlockKey, true)
}

// calculateProposalStatus 计算提案状态
// 返回提案状态码（uint8）
//
// 提案状态流转：
//
//	Pending(0) -> 待存款，提案已提交但尚未支付押金
//	Ongoing(1) -> 进行中，正在投票阶段
//	Confirming(2) -> 确认中，阈值已满足，正在等待确认期结束
//	Confirmed(3) -> 已确认，确认期结束，等待执行
//	Approved(4) -> 已批准，提案通过
//	Rejected(5) -> 已拒绝，提案被否决
//	Canceled(6) -> 已取消，提案被取消
//
// 投票计算公式：
//
//	Approval = yes / (yes + no) * 10000（万分比）
//	Support = support / totalSupply * 10000（万分比）
//
// 通过条件：
//  1. approvalRatio >= minApproval && supportRatio >= minSupport
//  2. 连续满足条件的时间 > confirmPeriod
func (d Gov) calculateProposalStatus(id uint32, prop Proposal) (ProposalStatusQuery, error) {
	track, err := d.track(prop.TrackID)
	if err != nil {
		return ProposalStatusQuery{}, err
	}
	end := prop.Deposit.Block + int64(track.MaxDeciding)
	height := d.api.GetHeight()

	// 使用 ListAll 获取该提案的所有投票
	_, votes, err := d.votes.ListAll(d.api.GetTxn(), id)
	if err != nil {
		return ProposalStatusQuery{}, err
	}

	totalSupply, err := d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
	if err != nil {
		return ProposalStatusQuery{}, err
	}

	// 统计结果
	yes := big.NewInt(0)
	no := big.NewInt(0)
	support := big.NewInt(0)
	all := totalSupply.Int

	// 确认期逻辑
	var lastAchieveBlock int64 = 0
	confirmPeriod := int64(track.ConfirmPeriod)
	decisionPeriod := int64(track.DecisionPeriod)
	confirmed := 0

	for _, vote := range votes {
		if vote.Deleted {
			continue
		}

		// 使用投票时的相对区块偏移量计算阈值
		// x = voteBlock - depositBlock，表示从投票期开始后的第几个区块
		voteBlock := vote.VoteBlock
		voteOffset := voteBlock - prop.Deposit.Block
		minApproval := uint64(track.MinApproval.Y(voteOffset))
		minSupport := uint64(track.MinSupport.Y(voteOffset))

		// 累计投票
		pledge := vote.Pledge.Int
		weightedPledge := new(big.Int).Mul(pledge, vote.VoteWeight.Int)
		support.Add(support, pledge)

		if vote.OpinionYes {
			yes.Add(yes, weightedPledge)
		} else {
			no.Add(no, weightedPledge)
		}

		// 计算比率（万分比）
		totalVotes := new(big.Int).Add(yes, no)
		if totalVotes.Sign() > 0 && all.Sign() > 0 {
			// Approval = yes / (yes + no) * 10000
			approvalRatio := new(big.Int).Mul(yes, big.NewInt(10000))
			approvalRatio.Div(approvalRatio, totalVotes)

			// Support = support / all * 10000
			supportRatio := new(big.Int).Mul(support, big.NewInt(10000))
			supportRatio.Div(supportRatio, all)

			// 检查是否满足阈值
			if approvalRatio.Uint64() >= minApproval && supportRatio.Uint64() >= minSupport {
				// 检查确认期
				if lastAchieveBlock > 0 && voteBlock-lastAchieveBlock > confirmPeriod {
					// 已确认，等待执行
					return ProposalStatusQuery{
						State:              ProposalStatusConfirmed,
						BlockHeight:        height,
						ConfirmedNumber:    uint32(confirmed),
						LastConfirmedBlock: lastAchieveBlock,
					}, nil
				}

				// 记录首次满足阈值的区块
				if lastAchieveBlock == 0 {
					lastAchieveBlock = voteBlock
					confirmed++
				}
			} else {
				// 不满足阈值，重置
				lastAchieveBlock = 0
			}
		}
	}

	// 检查当前区块与最后满足阈值区块之间是否已超过确认期
	if lastAchieveBlock > 0 && height-lastAchieveBlock > confirmPeriod {
		// 检查是否有足够时间完成确认期（必须在 decisionPeriod 结束前）
		if lastAchieveBlock+confirmPeriod <= end+decisionPeriod {
			// 已确认，等待执行
			return ProposalStatusQuery{
				State:              ProposalStatusConfirmed,
				BlockHeight:        height,
				ConfirmedNumber:    uint32(confirmed),
				LastConfirmedBlock: lastAchieveBlock,
			}, nil
		}
		// 确认期时间不够，拒绝
		return ProposalStatusQuery{
			State:              ProposalStatusRejected,
			BlockHeight:        height,
			ConfirmedNumber:    uint32(confirmed),
			LastConfirmedBlock: 0,
		}, nil
	}

	// 超过截止区块，拒绝
	if height > end {
		return ProposalStatusQuery{
			State:              ProposalStatusRejected,
			BlockHeight:        height,
			ConfirmedNumber:    uint32(confirmed),
			LastConfirmedBlock: 0,
		}, nil
	}

	// 正在确认中，检查是否有足够时间完成确认期
	if lastAchieveBlock > 0 {
		if lastAchieveBlock+confirmPeriod > end+decisionPeriod {
			// 确认期时间不够，拒绝
			return ProposalStatusQuery{
				State:              ProposalStatusRejected,
				BlockHeight:        height,
				ConfirmedNumber:    uint32(confirmed),
				LastConfirmedBlock: 0,
			}, nil
		}
		return ProposalStatusQuery{
			State:              ProposalStatusConfirming,
			BlockHeight:        height,
			ConfirmedNumber:    uint32(confirmed),
			LastConfirmedBlock: lastAchieveBlock,
		}, nil
	}

	// 进行中
	return ProposalStatusQuery{
		State:              ProposalStatusOngoing,
		BlockHeight:        height,
		ConfirmedNumber:    uint32(confirmed),
		LastConfirmedBlock: lastAchieveBlock,
	}, nil
}

// calculateProposalEndBlock 计算提案结束区块
func (d Gov) calculateProposalEndBlock(id uint32) (int64, error) {
	prop, err := d.proposal(id)
	if err != nil {
		return 0, err
	}
	switch prop.Status.State {
	case ProposalStatusRejected, ProposalStatusApproved:
		return prop.Status.Block, nil
	case ProposalStatusOngoing:
		status, _ := d.calculateProposalStatus(id, prop)
		if status.State == ProposalStatusRejected {
			track, _ := d.track(prop.TrackID)
			return prop.Deposit.Block + int64(track.MaxDeciding), nil
		}
	}
	return 0, ErrInvalidProposalStatus
}

// voteUnlockKey 生成投票解锁状态的 key
func voteUnlockKey(proposalID uint32, index uint32) uint64 {
	return uint64(proposalID)<<32 | uint64(index)
}
