package dao

import "github.com/wetee-dao/ink.go/util"

// 以下类型为 DAO 合约状态与调用的原生 Go 表示，不依赖 pkg/model/dao.pb.go。

// Member 成员账户与余额（u128 为大端字节）。
type Member struct {
	Account []byte
	Balance []byte
}

// TrackData 治理轨道参数。
type TrackData struct {
	Name               string
	PreparePeriod      uint32
	MaxDeciding        uint32
	ConfirmPeriod      uint32
	DecisionPeriod     uint32
	MinEnactmentPeriod uint32
	DecisionDeposit    []byte
	MaxBalance         []byte
}

// CallContent 提案内嵌调用内容。
type CallContent struct {
	Contract     []byte
	Selector     []byte
	Input        []byte
	Amount       []byte
	RefTimeLimit uint64
	AllowReentry bool
}

type ProposalStatus struct {
	State ProposalState
	Block int64
}

type ProposalDeposit struct {
	Depositor []byte
	Amount    []byte
	Block     int64
}

type Proposal struct {
	ID            uint32
	Call          util.Option[CallContent]
	TrackID       uint32
	Caller        []byte
	Status        ProposalStatus
	SubmitBlock   int64
	DecisionBlock int64
	Deposit       ProposalDeposit
}

type Vote struct {
	ID          uint64
	ProposalID  uint32
	Caller      []byte
	Pledge      []byte
	OpinionYes  bool
	VoteWeight  []byte
	UnlockBlock int64
	VoteBlock   int64
	Deleted     bool
}

type Spend struct {
	ID      uint64
	Caller  []byte
	To      []byte
	Amount  []byte
	TrackID uint32
	Payout  bool
}
