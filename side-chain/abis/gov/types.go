package gov

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go/util"
)

type Error struct { // Enum
	InsufficientAllowance *bool // 0
	InvalidDeposit        *bool // 1
	InvalidDepositTime    *bool // 2
	InvalidProposal       *bool // 3
	InvalidProposalCaller *bool // 4
	InvalidProposalStatus *bool // 5
	InvalidVote           *bool // 6
	InvalidVoteStatus     *bool // 7
	InvalidVoteTime       *bool // 8
	InvalidVoteUnlockTime *bool // 9
	InvalidVoteUser       *bool // 10
	LowBalance            *bool // 11
	MemberBalanceNotZero  *bool // 12
	MemberExisted         *bool // 13
	MemberNotExisted      *bool // 14
	MustCallByGov         *bool // 15
	NoTrack               *bool // 16
	PropNotOngoing        *bool // 17
	ProposalInDecision    *bool // 18
	ProposalNotConfirmed  *bool // 19
	PublicJoinNotAllowed  *bool // 20
	SpendAlreadyExecuted  *bool // 21
	SpendNotFound         *bool // 22
	TransferDisabled      *bool // 23
	VoteAlreadyUnlocked   *bool // 24
}

func (ty Error) Encode(encoder scale.Encoder) (err error) {
	if ty.InsufficientAllowance != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidDeposit != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidDepositTime != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidProposal != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidProposalCaller != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidProposalStatus != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidVote != nil {
		err = encoder.PushByte(6)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidVoteStatus != nil {
		err = encoder.PushByte(7)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidVoteTime != nil {
		err = encoder.PushByte(8)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidVoteUnlockTime != nil {
		err = encoder.PushByte(9)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidVoteUser != nil {
		err = encoder.PushByte(10)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.LowBalance != nil {
		err = encoder.PushByte(11)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MemberBalanceNotZero != nil {
		err = encoder.PushByte(12)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MemberExisted != nil {
		err = encoder.PushByte(13)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MemberNotExisted != nil {
		err = encoder.PushByte(14)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MustCallByGov != nil {
		err = encoder.PushByte(15)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NoTrack != nil {
		err = encoder.PushByte(16)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.PropNotOngoing != nil {
		err = encoder.PushByte(17)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProposalInDecision != nil {
		err = encoder.PushByte(18)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProposalNotConfirmed != nil {
		err = encoder.PushByte(19)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.PublicJoinNotAllowed != nil {
		err = encoder.PushByte(20)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.SpendAlreadyExecuted != nil {
		err = encoder.PushByte(21)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.SpendNotFound != nil {
		err = encoder.PushByte(22)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.TransferDisabled != nil {
		err = encoder.PushByte(23)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.VoteAlreadyUnlocked != nil {
		err = encoder.PushByte(24)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Error) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.InsufficientAllowance = &t
		return
	case 1: // Base
		t := true
		ty.InvalidDeposit = &t
		return
	case 2: // Base
		t := true
		ty.InvalidDepositTime = &t
		return
	case 3: // Base
		t := true
		ty.InvalidProposal = &t
		return
	case 4: // Base
		t := true
		ty.InvalidProposalCaller = &t
		return
	case 5: // Base
		t := true
		ty.InvalidProposalStatus = &t
		return
	case 6: // Base
		t := true
		ty.InvalidVote = &t
		return
	case 7: // Base
		t := true
		ty.InvalidVoteStatus = &t
		return
	case 8: // Base
		t := true
		ty.InvalidVoteTime = &t
		return
	case 9: // Base
		t := true
		ty.InvalidVoteUnlockTime = &t
		return
	case 10: // Base
		t := true
		ty.InvalidVoteUser = &t
		return
	case 11: // Base
		t := true
		ty.LowBalance = &t
		return
	case 12: // Base
		t := true
		ty.MemberBalanceNotZero = &t
		return
	case 13: // Base
		t := true
		ty.MemberExisted = &t
		return
	case 14: // Base
		t := true
		ty.MemberNotExisted = &t
		return
	case 15: // Base
		t := true
		ty.MustCallByGov = &t
		return
	case 16: // Base
		t := true
		ty.NoTrack = &t
		return
	case 17: // Base
		t := true
		ty.PropNotOngoing = &t
		return
	case 18: // Base
		t := true
		ty.ProposalInDecision = &t
		return
	case 19: // Base
		t := true
		ty.ProposalNotConfirmed = &t
		return
	case 20: // Base
		t := true
		ty.PublicJoinNotAllowed = &t
		return
	case 21: // Base
		t := true
		ty.SpendAlreadyExecuted = &t
		return
	case 22: // Base
		t := true
		ty.SpendNotFound = &t
		return
	case 23: // Base
		t := true
		ty.TransferDisabled = &t
		return
	case 24: // Base
		t := true
		ty.VoteAlreadyUnlocked = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}
func (ty *Error) Error() string {
	if ty.InsufficientAllowance != nil {
		return "InsufficientAllowance"
	}

	if ty.InvalidDeposit != nil {
		return "InvalidDeposit"
	}

	if ty.InvalidDepositTime != nil {
		return "InvalidDepositTime"
	}

	if ty.InvalidProposal != nil {
		return "InvalidProposal"
	}

	if ty.InvalidProposalCaller != nil {
		return "InvalidProposalCaller"
	}

	if ty.InvalidProposalStatus != nil {
		return "InvalidProposalStatus"
	}

	if ty.InvalidVote != nil {
		return "InvalidVote"
	}

	if ty.InvalidVoteStatus != nil {
		return "InvalidVoteStatus"
	}

	if ty.InvalidVoteTime != nil {
		return "InvalidVoteTime"
	}

	if ty.InvalidVoteUnlockTime != nil {
		return "InvalidVoteUnlockTime"
	}

	if ty.InvalidVoteUser != nil {
		return "InvalidVoteUser"
	}

	if ty.LowBalance != nil {
		return "LowBalance"
	}

	if ty.MemberBalanceNotZero != nil {
		return "MemberBalanceNotZero"
	}

	if ty.MemberExisted != nil {
		return "MemberExisted"
	}

	if ty.MemberNotExisted != nil {
		return "MemberNotExisted"
	}

	if ty.MustCallByGov != nil {
		return "MustCallByGov"
	}

	if ty.NoTrack != nil {
		return "NoTrack"
	}

	if ty.PropNotOngoing != nil {
		return "PropNotOngoing"
	}

	if ty.ProposalInDecision != nil {
		return "ProposalInDecision"
	}

	if ty.ProposalNotConfirmed != nil {
		return "ProposalNotConfirmed"
	}

	if ty.PublicJoinNotAllowed != nil {
		return "PublicJoinNotAllowed"
	}

	if ty.SpendAlreadyExecuted != nil {
		return "SpendAlreadyExecuted"
	}

	if ty.SpendNotFound != nil {
		return "SpendNotFound"
	}

	if ty.TransferDisabled != nil {
		return "TransferDisabled"
	}

	if ty.VoteAlreadyUnlocked != nil {
		return "VoteAlreadyUnlocked"
	}
	return "Unknown"
}

type UniAddr struct { // Composite
	T uint32
	V []byte
}
type TrackData struct { // Composite
	Name               string
	PreparePeriod      uint32
	MaxDeciding        uint32
	ConfirmPeriod      uint32
	DecisionPeriod     uint32
	MinEnactmentPeriod uint32
	DecisionDeposit    types.U256
	MaxBalance         types.U256
}
type CallContent struct { // Composite
	Contract []byte
	Selector [4]byte
	Args     [][]byte
	Amount   types.U256
}
type Member struct { // Composite
	Account UniAddr
	Balance types.U256
}
type Proposal struct { // Composite
	ID            uint32
	Call          util.Option[CallContent]
	TrackID       uint32
	Caller        UniAddr
	Status        ProposalStatus
	SubmitBlock   int64
	DecisionBlock int64
	Deposit       ProposalDeposit
}
type ProposalStatus struct { // Composite
	State string
	Block int64
}
type ProposalDeposit struct { // Composite
	Depositor UniAddr
	Amount    types.U256
	Block     int64
}
type Vote struct { // Composite
	ID          uint64
	ProposalID  uint32
	Caller      UniAddr
	Pledge      types.U256
	OpinionYes  bool
	VoteWeight  types.U256
	UnlockBlock int64
	VoteBlock   int64
	Deleted     bool
}
