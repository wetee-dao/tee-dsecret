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
	if prop.Status.State != ProposalOngoing {
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
	if prop.Status.State != ProposalOngoing {
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
// 返回：是否确认通过、结束区块、轨道数据
func (d Gov) calculateProposalStatus(id uint32, prop Proposal) (bool, int64, TrackData, error) {
	track, err := d.track(prop.TrackID)
	if err != nil {
		return false, 0, TrackData{}, err
	}
	end := prop.Deposit.Block + int64(track.MaxDeciding)

	// 使用 ListAll 获取该提案的所有投票
	_, votes, err := d.votes.ListAll(d.api.GetTxn(), id)
	if err != nil {
		return false, 0, TrackData{}, err
	}

	totalSupply, err := d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
	if err != nil {
		return false, 0, TrackData{}, err
	}

	// 统计结果
	yes := big.NewInt(0)
	no := big.NewInt(0)
	support := big.NewInt(0)
	all := totalSupply.Int

	// 确认期逻辑
	isConfirm := false
	var lastAchieveBlock int64 = 0
	confirmPeriod := int64(track.ConfirmPeriod)

	for _, vote := range votes {
		if vote.Deleted {
			continue
		}

		// 使用投票时的区块号计算阈值
		voteBlock := vote.VoteBlock
		minApproval := uint64(track.MinApproval.Y(voteBlock))
		minSupport := uint64(track.MinSupport.Y(voteBlock))

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
					isConfirm = true
					break
				}

				// 记录首次满足阈值的区块
				if lastAchieveBlock == 0 {
					lastAchieveBlock = voteBlock
				}
			} else {
				// 不满足阈值，重置
				lastAchieveBlock = 0
			}
		}
	}

	return isConfirm, end, track, nil
}

// calculateProposalEndBlock 计算提案结束区块
func (d Gov) calculateProposalEndBlock(id uint32) (int64, error) {
	height := d.api.GetHeight()
	prop, err := d.proposal(id)
	if err != nil {
		return 0, err
	}
	switch prop.Status.State {
	case ProposalRejected, ProposalApproved:
		return prop.Status.Block, nil
	case ProposalOngoing:
		confirmed, end, _, err := d.calculateProposalStatus(id, prop)
		if err != nil {
			return 0, err
		}
		if !confirmed && height > end {
			return end, nil
		}
	}
	return 0, ErrInvalidProposalStatus
}

// voteUnlockKey 生成投票解锁状态的 key
func voteUnlockKey(proposalID uint32, index uint32) uint64 {
	return uint64(proposalID)<<32 | uint64(index)
}
