// Package dao 实现 DAO 状态存储与交易应用逻辑，对应 ink DAO 智能合约能力，
// 在 CometBFT 侧链上提供治理、成员、代币、提案与国库功能。
package dao

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// U128 使用 model.U128，与合约 u128 对齐。
type U128 = model.U128

const (
	DaoNamespace = "dao"
	// 以下 Op 常量仅用于文档/daogen，实际分发由 model.DaoCall 的 oneof 决定。
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

// DaoTrackData、DaoCallContent 用于存储层 JSON 序列化（含 U128），与 model 的 proto 类型对应。
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

// model 的 proto 转为存储用类型（u128 用 bytes 表示）
func trackFromProto(m *model.DaoTrackData) *DaoTrackData {
	if m == nil {
		return nil
	}
	return &DaoTrackData{
		Name:               m.Name,
		PreparePeriod:      m.PreparePeriod,
		MaxDeciding:        m.MaxDeciding,
		ConfirmPeriod:      m.ConfirmPeriod,
		DecisionPeriod:     m.DecisionPeriod,
		MinEnactmentPeriod: m.MinEnactmentPeriod,
		DecisionDeposit:    model.NewU128(model.BytesToU128(m.DecisionDeposit)),
		MaxBalance:         model.NewU128(model.BytesToU128(m.MaxBalance)),
	}
}

func callContentFromProto(m *model.DaoCallContent) *DaoCallContent {
	if m == nil {
		return nil
	}
	return &DaoCallContent{
		Contract:     m.Contract,
		Selector:     m.Selector,
		Input:        m.Input,
		Amount:       model.NewU128(model.BytesToU128(m.Amount)),
		RefTimeLimit: m.RefTimeLimit,
		AllowReentry: m.AllowReentry,
	}
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

// ApplyDaoCall 解析 payload（protobuf model.DaoCall）并执行对应 DAO 操作，由 sidechain 在 FinalizeTx 中调用。
func ApplyDaoCall(caller []byte, payload []byte, height int64, txn *model.Txn) error {
	if len(payload) == 0 {
		return errors.New("dao_call: empty payload")
	}
	var dc model.DaoCall
	if err := dc.Unmarshal(payload); err != nil {
		return fmt.Errorf("dao_call: invalid proto: %w", err)
	}
	switch v := dc.Call.(type) {
	case *model.DaoCall_Init:
		return daoInit(caller, v.Init, height, txn)
	case *model.DaoCall_PublicJoin:
		return daoPublicJoin(caller, txn)
	case *model.DaoCall_Join:
		return daoJoin(caller, v.Join, txn)
	case *model.DaoCall_Leave:
		return daoLeave(caller, txn)
	case *model.DaoCall_LeaveWithBurn:
		return daoLeaveWithBurn(caller, txn)
	case *model.DaoCall_SubmitProposal:
		return daoSubmitProposal(caller, v.SubmitProposal, height, txn)
	case *model.DaoCall_DepositProposal:
		return daoDepositProposal(caller, v.DepositProposal, height, txn)
	case *model.DaoCall_SubmitVote:
		return daoSubmitVote(caller, v.SubmitVote, height, txn)
	case *model.DaoCall_CancelVote:
		return daoCancelVote(caller, v.CancelVote, height, txn)
	case *model.DaoCall_Unlock:
		return daoUnlock(caller, v.Unlock, height, txn)
	case *model.DaoCall_ExecProposal:
		return daoExecProposal(caller, v.ExecProposal, height, txn)
	case *model.DaoCall_CancelProposal:
		return daoCancelProposal(caller, v.CancelProposal, height, txn)
	case *model.DaoCall_Transfer:
		return daoTransfer(caller, v.Transfer, txn)
	case *model.DaoCall_Approve:
		return daoApprove(caller, v.Approve, txn)
	case *model.DaoCall_TransferFrom:
		return daoTransferFrom(caller, v.TransferFrom, txn)
	case *model.DaoCall_Spend:
		return daoSpend(caller, v.Spend, txn)
	case *model.DaoCall_Payout:
		return daoPayout(caller, v.Payout, txn)
	case *model.DaoCall_SetPublicJoin:
		return daoSetPublicJoin(caller, v.SetPublicJoin, txn)
	case *model.DaoCall_AddTrack:
		return daoAddTrack(caller, v.AddTrack, txn)
	case *model.DaoCall_SetDefaultTrack:
		return daoSetDefaultTrack(caller, v.SetDefaultTrack, txn)
	default:
		return errors.New("dao_call: unknown call type")
	}
}

func isSudo(caller []byte, state *daoStateState) bool {
	acc := state.SudoAccount()
	if len(acc) == 0 {
		return false
	}
	return bytesEqual(caller, acc)
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
