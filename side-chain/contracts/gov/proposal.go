package gov

import (
	"bytes"

	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d GovQuery) Proposal(id uint32) (util.Option[ProposalWithID], error) {
	prop, err := d.proposals.Get(d.api.GetTxn(), id)
	if err != nil {
		return util.NewNone[ProposalWithID](), err
	}
	if prop == nil {
		return util.NewNone[ProposalWithID](), ErrInvalidProposal
	}

	// 填充 ID（StoreList 不自动存储 key）
	p := ProposalWithID{
		ID:       id,
		Proposal: *prop,
	}

	return util.NewSome(p), nil
}

func (d GovQuery) Proposals(startKey util.Option[uint32], size uint32) ([]ProposalWithID, error) {
	var startKeyPtr = new(uint32)
	if startKey.IsSome() {
		startKeyPtr = &startKey.V
	}

	ids, props, err := d.proposals.DescList(d.api.GetTxn(), startKeyPtr, size)
	if err != nil {
		return nil, err
	}

	// 填充每个 proposal 的 ID
	ps := make([]ProposalWithID, len(ids))
	for i := range ids {
		ps = append(ps, ProposalWithID{
			ID:       ids[i],
			Proposal: props[i],
		})
	}

	return ps, nil
}

func (d GovQuery) ProposalStatus(id uint32) (util.Option[ProposalStatus], error) {
	height := d.api.GetHeight()
	prop, err := d.proposal(id)
	if err != nil {
		return util.NewNone[ProposalStatus](), err
	}
	if prop.Status.State != ProposalOngoing {
		return util.NewSome(prop.Status), nil
	}

	confirmed, end, _, err := d.calculateProposalStatus(id, prop)
	if err != nil {
		return util.NewNone[ProposalStatus](), err
	}
	if !confirmed && height > end {
		status := ProposalStatus{State: ProposalRejected, Block: end}
		return util.NewSome(status), nil
	}
	if confirmed {
		status := ProposalStatus{State: ProposalApproved, Block: 0}
		return util.NewSome(status), nil
	}
	status := ProposalStatus{State: ProposalOngoing, Block: 0}
	return util.NewSome(status), nil
}

func (d GovMutation) SubmitProposal(call CallContent, trackID uint32) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	member, err := d.members.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if member.Int.Sign() == 0 {
		return ErrMemberNotExisted
	}

	track, err := d.tracks.Get(d.api.GetTxn(), trackID)
	if err != nil {
		return err
	}
	if track == nil {
		return ErrNoTrack
	}

	amountInt := call.Amount
	maxBalanceInt := track.MaxBalance
	if amountInt.Cmp(maxBalanceInt.Int) > 0 {
		return ErrLowBalance
	}

	// 使用 StoreList.Insert 自动分配 ID
	record := Proposal{
		Call:        util.NewSome(call),
		TrackID:     trackID,
		Caller:      caller,
		Status:      ProposalStatus{State: ProposalPending},
		SubmitBlock: height,
		Deposit: ProposalDeposit{
			Depositor: caller,
			Amount:    model.ZeroAmount,
			Block:     height,
		},
	}
	id, err := d.proposals.Insert(d.api.GetTxn(), record)
	if err != nil {
		return err
	}
	_ = id // ID 已由 StoreList 自动分配，无需额外处理
	return nil
}

func (d GovMutation) CancelProposal(proposalID uint32) error {
	caller := d.api.GetCaller()
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalPending {
		return ErrInvalidProposalStatus
	}
	if !bytes.Equal(prop.Caller.V, caller.V) {
		return ErrInvalidProposalCaller
	}
	prop.Status = ProposalStatus{State: ProposalCanceled}
	return d.proposals.Update(d.api.GetTxn(), proposalID, prop)
}

func (d GovMutation) DepositProposal(proposalID uint32, amount model.Amount) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalPending {
		return ErrInvalidProposalStatus
	}
	track, err := d.track(prop.TrackID)
	if err != nil {
		return err
	}
	if height < prop.Deposit.Block+int64(track.PreparePeriod) {
		return ErrInvalidDepositTime
	}

	depositInt := track.DecisionDeposit
	if amount.Int.Cmp(depositInt.Int) < 0 {
		return ErrInvalidDeposit
	}

	prop.Deposit = ProposalDeposit{Depositor: caller, Amount: amount, Block: height}
	prop.Status = ProposalStatus{State: ProposalOngoing}
	return d.proposals.Update(d.api.GetTxn(), proposalID, prop)
}

func (d GovMutation) ExecProposal(proposalID uint32) error {
	height := d.api.GetHeight()
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}

	if prop.Status.State != ProposalOngoing {
		return ErrPropNotOngoing
	}

	confirmed, end, track, err := d.calculateProposalStatus(proposalID, prop)
	if err != nil {
		return err
	}

	if !confirmed {
		if height > end {
			prop.Status = ProposalStatus{State: ProposalRejected, Block: end}
			return d.proposals.Update(d.api.GetTxn(), proposalID, prop)
		}
		return ErrProposalNotConfirmed
	}

	if height < end+int64(track.DecisionPeriod) {
		return ErrProposalInDecision
	}

	// 执行提案中的调用并记录结果
	result := ProposalResult{}
	if prop.Call.IsSome() {
		call, _ := prop.Call.UnWrap()
		res, err := d.api.Call(model.UniAddr{}, model.ContractCall{
			Contract: string(call.Contract),
			Method:   call.Selector,
			Args:     call.Args,
		})
		if err != nil {
			// 记录错误但不返回，提案仍标记为已批准
			result.ExecError = util.NewSome([]byte(err.Error()))
		} else {
			result.Result = util.NewSome(res)
		}

		// 保存执行结果
		if err := d.proposalResults.Set(d.api.GetTxn(), proposalID, result); err != nil {
			return err
		}
	}

	prop.Status = ProposalStatus{State: ProposalApproved, Block: height}
	// prop.DecisionBlock = height
	return d.proposals.Update(d.api.GetTxn(), proposalID, prop)
}
