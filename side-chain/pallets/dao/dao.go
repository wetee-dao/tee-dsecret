// Package dao 实现 DAO 状态存储与交易应用逻辑，对应 ink DAO 智能合约能力，
// 在 CometBFT 侧链上提供治理、成员、代币、提案与国库功能。
package dao

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// U128 使用 model.U128，与合约 u128 对齐。
type U128 = model.U128

const (
	DaoNamespace = "dao"
)

// DaoCallPayload 为 dao_call 的 JSON 载荷，与 ink DAO 消息一一对应。
type DaoCallPayload struct {
	Op string `json:"op"` // 见下 OpXxx
	// 以下按 op 使用
	InitialMembers []DaoMember     `json:"initial_members,omitempty"`
	PublicJoin     *bool           `json:"public_join,omitempty"`
	SudoAccount    []byte          `json:"sudo_account,omitempty"`
	DefaultTrack   *DaoTrackData   `json:"default_track,omitempty"`
	NewUser        []byte          `json:"new_user,omitempty"`
	Balance        U128            `json:"balance,omitempty"`
	ProposalId     uint32          `json:"proposal_id,omitempty"`
	TrackId        uint32          `json:"track_id,omitempty"`
	Call           *DaoCallContent `json:"call,omitempty"`
	Amount         U128            `json:"amount,omitempty"`
	OpinionYes     bool            `json:"opinion_yes,omitempty"`
	LockAmount     U128            `json:"lock_amount,omitempty"`
	VoteId         uint64          `json:"vote_id,omitempty"`
	To             []byte          `json:"to,omitempty"`
	Value          U128            `json:"value,omitempty"`
	Spender        []byte          `json:"spender,omitempty"`
	From           []byte          `json:"from,omitempty"`
	SpendId        uint64          `json:"spend_id,omitempty"`
	Track          *DaoTrackData   `json:"track,omitempty"`
}

const (
	OpDaoInit            = "dao_init"
	OpDaoPublicJoin      = "dao_public_join"
	OpDaoJoin            = "dao_join"
	OpDaoLeave           = "dao_leave"
	OpDaoLeaveWithBurn   = "dao_leave_with_burn"
	OpDaoSubmitProposal  = "dao_submit_proposal"
	OpDaoDepositProposal = "dao_deposit_proposal"
	OpDaoSubmitVote      = "dao_submit_vote"
	OpDaoCancelVote      = "dao_cancel_vote"
	OpDaoUnlock          = "dao_unlock"
	OpDaoExecProposal    = "dao_exec_proposal"
	OpDaoCancelProposal  = "dao_cancel_proposal"
	OpDaoTransfer        = "dao_transfer"
	OpDaoApprove         = "dao_approve"
	OpDaoTransferFrom    = "dao_transfer_from"
	OpDaoSpend           = "dao_spend"
	OpDaoPayout          = "dao_payout"
	OpDaoSetPublicJoin   = "dao_set_public_join"
	OpDaoAddTrack        = "dao_add_track"
	OpDaoSetDefaultTrack = "dao_set_default_track"
)

type DaoMember struct {
	Account []byte `json:"account"`
	Balance U128   `json:"balance"`
}

type DaoTrackData struct {
	Name               string `json:"name"`
	PreparePeriod      uint32 `json:"prepare_period"`
	MaxDeciding        uint32 `json:"max_deciding"`
	ConfirmPeriod      uint32 `json:"confirm_period"`
	DecisionPeriod     uint32 `json:"decision_period"`
	MinEnactmentPeriod uint32 `json:"min_enactment_period"`
	DecisionDeposit    U128   `json:"decision_deposit"`
	MaxBalance         U128   `json:"max_balance"`
}

type DaoCallContent struct {
	Contract     []byte `json:"contract,omitempty"`
	Selector     []byte `json:"selector,omitempty"`
	Input        []byte `json:"input,omitempty"`
	Amount       U128   `json:"amount,omitempty"`
	RefTimeLimit uint64 `json:"ref_time_limit,omitempty"`
	AllowReentry bool   `json:"allow_reentry,omitempty"`
}

type daoState struct {
	Members        [][]byte `json:"members"`
	TotalIssuance  *big.Int `json:"-"` // u128，不参与 JSON
	PublicJoin     bool     `json:"public_join"`
	SudoAccount    []byte   `json:"sudo_account,omitempty"`
	Transfer       bool     `json:"transfer"`
	DefaultTrack   *uint32  `json:"default_track,omitempty"`
	NextProposalId uint32   `json:"next_proposal_id"`
	NextVoteId     uint64   `json:"next_vote_id"`
	NextSpendId    uint64   `json:"next_spend_id"`

	Tracks         *model.StoreMapping[uint32] `state:"key=track_"`
	Proposals      *model.StoreMapping[uint32] `state:"key=proposal_"`
	ProposalStatus *model.StoreMapping[uint32] `state:"key=proposal_status_"`
	ProposalTrack  *model.StoreMapping[uint32] `state:"key=proposal_track_"`
	ProposalCaller *model.StoreMapping[uint32] `state:"key=proposal_caller_"`
	SubmitBlock    *model.StoreMapping[uint32] `state:"key=submit_block_"`
	Deposits       *model.StoreMapping[uint32] `state:"key=deposit_"`
	Votes          *model.StoreMapping[uint64] `state:"key=vote_"`
	VoteOfMember   *model.StoreMapping[[]byte] `state:"key=vote_of_member_"`
	Unlock         *model.StoreMapping[uint64] `state:"key=unlock_"`
	Spends         *model.StoreMapping[uint64] `state:"key=spend_"`
	MemberBalance  *model.StoreMapping[[]byte] `state:"key=member_balance_"`
	MemberLock     *model.StoreMapping[[]byte] `state:"key=member_lock_"`
	Allowance      *model.StoreMapping[string] `state:"key=allowance_"`
}

type proposalDeposit struct {
	Depositor []byte `json:"depositor"`
	Amount    U128   `json:"amount"`
	Block     int64  `json:"block"`
}

type voteInfo struct {
	Pledge      U128   `json:"pledge"`
	OpinionYes  bool   `json:"opinion_yes"`
	CallId      uint32 `json:"call_id"`
	Caller      []byte `json:"caller"`
	VoteBlock   int64  `json:"vote_block"`
	UnlockBlock int64  `json:"unlock_block"`
	Deleted     bool   `json:"deleted"`
}

type spendRecord struct {
	Caller []byte `json:"caller"`
	To     []byte `json:"to"`
	Amount U128   `json:"amount"`
	Payout bool   `json:"payout"`
}

func addrKey(a []byte) string { return hex.EncodeToString(a) }

// ApplyDaoCall 解析 payload 并执行对应 DAO 操作，由 sidechain 在 FinalizeTx 中调用。
func ApplyDaoCall(caller []byte, payload []byte, height int64, txn *model.Txn) error {
	if len(payload) == 0 {
		return errors.New("dao_call: empty payload")
	}
	var p DaoCallPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("dao_call: invalid json: %w", err)
	}
	switch p.Op {
	case OpDaoInit:
		return daoInit(caller, &p, height, txn)
	case OpDaoPublicJoin:
		return daoPublicJoin(caller, txn)
	case OpDaoJoin:
		return daoJoin(caller, &p, txn)
	case OpDaoLeave:
		return daoLeave(caller, txn)
	case OpDaoLeaveWithBurn:
		return daoLeaveWithBurn(caller, txn)
	case OpDaoSubmitProposal:
		return daoSubmitProposal(caller, &p, height, txn)
	case OpDaoDepositProposal:
		return daoDepositProposal(caller, &p, height, txn)
	case OpDaoSubmitVote:
		return daoSubmitVote(caller, &p, height, txn)
	case OpDaoCancelVote:
		return daoCancelVote(caller, &p, height, txn)
	case OpDaoUnlock:
		return daoUnlock(caller, &p, height, txn)
	case OpDaoExecProposal:
		return daoExecProposal(caller, &p, height, txn)
	case OpDaoCancelProposal:
		return daoCancelProposal(caller, &p, height, txn)
	case OpDaoTransfer:
		return daoTransfer(caller, &p, txn)
	case OpDaoApprove:
		return daoApprove(caller, &p, txn)
	case OpDaoTransferFrom:
		return daoTransferFrom(caller, &p, txn)
	case OpDaoSpend:
		return daoSpend(caller, &p, txn)
	case OpDaoPayout:
		return daoPayout(caller, &p, txn)
	case OpDaoSetPublicJoin:
		return daoSetPublicJoin(caller, &p, txn)
	case OpDaoAddTrack:
		return daoAddTrack(caller, &p, txn)
	case OpDaoSetDefaultTrack:
		return daoSetDefaultTrack(caller, &p, txn)
	default:
		return fmt.Errorf("dao_call: unknown op %s", p.Op)
	}
}

func loadDaoState(txn *model.Txn) (*daoState, error) {
	state := newDaoStateState(txn)
	total := state.TotalIssuance()
	if total == nil {
		total = big.NewInt(0)
	}
	return &daoState{
		Members:        state.Members(),
		TotalIssuance:  total,
		PublicJoin:     state.PublicJoin(),
		SudoAccount:    state.SudoAccount(),
		Transfer:       state.Transfer(),
		DefaultTrack:   state.DefaultTrack(),
		NextProposalId: state.NextProposalId(),
		NextVoteId:     state.NextVoteId(),
		NextSpendId:    state.NextSpendId(),
		Tracks:         state.Tracks,
		Proposals:      state.Proposals,
		ProposalStatus: state.ProposalStatus,
		ProposalTrack:  state.ProposalTrack,
		ProposalCaller: state.ProposalCaller,
		SubmitBlock:    state.SubmitBlock,
		Deposits:       state.Deposits,
		Votes:          state.Votes,
		VoteOfMember:   state.VoteOfMember,
		Unlock:         state.Unlock,
		Spends:         state.Spends,
		MemberBalance:  state.MemberBalance,
		MemberLock:     state.MemberLock,
		Allowance:      state.Allowance,
	}, nil
}

func isSudo(caller []byte, st *daoState) bool {
	if len(st.SudoAccount) == 0 {
		return false
	}
	return bytesEqual(caller, st.SudoAccount)
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func memberInList(members [][]byte, k []byte) bool {
	for _, m := range members {
		if bytesEqual(m, k) {
			return true
		}
	}
	return false
}

func removeMember(members [][]byte, k []byte) [][]byte {
	out := make([][]byte, 0, len(members))
	for _, m := range members {
		if !bytesEqual(m, k) {
			out = append(out, m)
		}
	}
	return out
}
