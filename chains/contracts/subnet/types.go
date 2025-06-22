package subnet

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go/util"
)

type AccountId = [32]byte // Composite
type K8sCluster struct {  // Composite
	Name          []byte
	Owner         types.H160
	Level         byte
	StartBlock    uint32
	StopBlock     util.Option[uint32]
	TerminalBlock util.Option[uint32]
	ValidatorId   AccountId
	P2pId         AccountId
	Ip            Ip
	Port          uint32
	Status        byte
}
type Ip struct { // Composite
	Ipv4   util.Option[uint32]
	Ipv6   util.Option[types.U128]
	Domain util.Option[[]byte]
}
type SecretNode struct { // Composite
	Name          []byte
	Owner         types.H160
	ValidatorId   AccountId
	P2pId         AccountId
	StartBlock    uint32
	StopBlock     util.Option[uint32]
	TerminalBlock util.Option[uint32]
	Ip            Ip
	Port          uint32
	Status        byte
}
type Error struct { // Enum
	NotEnoughBalance        *bool // 0
	MustCallByMainContract  *bool // 1
	WorkerNotExist          *bool // 2
	WorkerNotOwnedByCaller  *bool // 3
	WorkerStatusNotReady    *bool // 4
	WorkerMortgageNotExist  *bool // 5
	TransferFailed          *bool // 6
	WorkerIsUseByUser       *bool // 7
	NodeNotExist            *bool // 8
	SecretNodeAlreadyExists *bool // 9
	SetCodeFailed           *bool // 10
	EpochNotExpired         *bool // 11
}

func (ty Error) Encode(encoder scale.Encoder) (err error) {
	if ty.NotEnoughBalance != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MustCallByMainContract != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerNotExist != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerNotOwnedByCaller != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerStatusNotReady != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerMortgageNotExist != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.TransferFailed != nil {
		err = encoder.PushByte(6)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerIsUseByUser != nil {
		err = encoder.PushByte(7)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NodeNotExist != nil {
		err = encoder.PushByte(8)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.SecretNodeAlreadyExists != nil {
		err = encoder.PushByte(9)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.SetCodeFailed != nil {
		err = encoder.PushByte(10)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.EpochNotExpired != nil {
		err = encoder.PushByte(11)
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
	case 0:
		t := true
		ty.NotEnoughBalance = &t
		return
	case 1:
		t := true
		ty.MustCallByMainContract = &t
		return
	case 2:
		t := true
		ty.WorkerNotExist = &t
		return
	case 3:
		t := true
		ty.WorkerNotOwnedByCaller = &t
		return
	case 4:
		t := true
		ty.WorkerStatusNotReady = &t
		return
	case 5:
		t := true
		ty.WorkerMortgageNotExist = &t
		return
	case 6:
		t := true
		ty.TransferFailed = &t
		return
	case 7:
		t := true
		ty.WorkerIsUseByUser = &t
		return
	case 8:
		t := true
		ty.NodeNotExist = &t
		return
	case 9:
		t := true
		ty.SecretNodeAlreadyExists = &t
		return
	case 10:
		t := true
		ty.SetCodeFailed = &t
		return
	case 11:
		t := true
		ty.EpochNotExpired = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}
func (ty *Error) Error() string {
	if ty.NotEnoughBalance != nil {
		return "NotEnoughBalance"
	}

	if ty.MustCallByMainContract != nil {
		return "MustCallByMainContract"
	}

	if ty.WorkerNotExist != nil {
		return "WorkerNotExist"
	}

	if ty.WorkerNotOwnedByCaller != nil {
		return "WorkerNotOwnedByCaller"
	}

	if ty.WorkerStatusNotReady != nil {
		return "WorkerStatusNotReady"
	}

	if ty.WorkerMortgageNotExist != nil {
		return "WorkerMortgageNotExist"
	}

	if ty.TransferFailed != nil {
		return "TransferFailed"
	}

	if ty.WorkerIsUseByUser != nil {
		return "WorkerIsUseByUser"
	}

	if ty.NodeNotExist != nil {
		return "NodeNotExist"
	}

	if ty.SecretNodeAlreadyExists != nil {
		return "SecretNodeAlreadyExists"
	}

	if ty.SetCodeFailed != nil {
		return "SetCodeFailed"
	}

	if ty.EpochNotExpired != nil {
		return "EpochNotExpired"
	}
	return "Unknown"
}

type Tuple_69 struct { // Tuple
	F0 uint32
	F1 uint32
	F2 uint32
}
