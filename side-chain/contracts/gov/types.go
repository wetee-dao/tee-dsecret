package gov

import (
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// 以下类型为 Gov 合约状态与调用的原生 Go 表示，不依赖 pkg/model/Gov.pb.go。

// Member 成员账户与余额。
type Member struct {
	Account model.UniAddr
	Balance model.Amount
}

// TrackData 治理轨道参数。
type TrackData struct {
	Name               string       // 轨道名称
	PreparePeriod      uint32       // 准备期（区块数），提案提交后需要等待的准备时间
	MaxDeciding        uint32       // 最大决策期（区块数），投票持续的最大时间
	ConfirmPeriod      uint32       // 确认期（区块数），阈值满足后需要持续的时间才能通过
	DecisionPeriod     uint32       // 决定期（区块数）
	MinEnactmentPeriod uint32       // 最小执行期（区块数），提案通过后到执行的最短时间
	DecisionDeposit    model.Amount // 决定押金金额，提案人需要质押的金额
	MaxBalance         model.Amount // 最大能执行的金额，提案涉及金额的上限
	MinApproval        Curve        // 投票通过阈值曲线，yes/(yes+no) 需达到的比例
	MinSupport         Curve        // 投票参与率曲线，support/totalSupply 需达到的比例
}

type TrackWithID struct {
	ID    uint32
	Track TrackData
}

// CallContent 提案内嵌调用内容。
type CallContent struct {
	Contract []byte
	Selector [4]byte
	Args     [][]byte
	Amount   model.Amount
}

type ProposalStatus struct {
	State uint8
	Block int64
}

type ProposalDeposit struct {
	Depositor model.UniAddr
	Amount    model.Amount
	Block     int64
}

type Proposal struct {
	Call        util.Option[CallContent]
	TrackID     uint32
	Caller      model.UniAddr
	Status      ProposalStatus
	SubmitBlock int64
	Deposit     ProposalDeposit
}

type ProposalWithID struct {
	ID       uint32
	Proposal Proposal
}

type ProposalResult struct {
	Result    util.Option[[]byte]
	ExecError util.Option[[]byte]
}

type Vote struct {
	Index       uint32
	ProposalID  uint32
	Caller      model.UniAddr
	Pledge      model.Amount
	OpinionYes  bool
	VoteWeight  model.Amount
	UnlockBlock int64
	VoteBlock   int64
	Deleted     bool
}

type Spend struct {
	ID      uint64
	Caller  model.UniAddr
	To      model.UniAddr
	Amount  model.Amount
	TrackID uint32
	Payout  bool
}
