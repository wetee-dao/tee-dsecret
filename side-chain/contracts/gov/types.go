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
	Name               string
	PreparePeriod      uint32
	MaxDeciding        uint32
	ConfirmPeriod      uint32
	DecisionPeriod     uint32
	MinEnactmentPeriod uint32
	DecisionDeposit    model.Amount // 保持 []byte，因为从链上读取
	MaxBalance         model.Amount // 保持 []byte，因为从链上读取
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
	ID            uint32
	Call          util.Option[CallContent]
	TrackID       uint32
	Caller        model.UniAddr
	Status        ProposalStatus
	SubmitBlock   int64
	DecisionBlock int64
	Deposit       ProposalDeposit
}

type Vote struct {
	ID          uint64
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
