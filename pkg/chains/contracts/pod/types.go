package pod

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type Error struct { // Enum
	SetCodeFailed           *bool // 0
	MustCallByCloudContract *bool // 1
	InsufficientBalance     *bool // 2
	TransferFailed          *bool // 3
	NotOwner                *bool // 4
	NotEnoughAllowance      *bool // 5
	NotEnoughBalance        *bool // 6
}

func (ty Error) Encode(encoder scale.Encoder) (err error) {
	if ty.SetCodeFailed != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MustCallByCloudContract != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InsufficientBalance != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.TransferFailed != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NotOwner != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NotEnoughAllowance != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NotEnoughBalance != nil {
		err = encoder.PushByte(6)
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
		ty.SetCodeFailed = &t
		return
	case 1: // Base
		t := true
		ty.MustCallByCloudContract = &t
		return
	case 2: // Base
		t := true
		ty.InsufficientBalance = &t
		return
	case 3: // Base
		t := true
		ty.TransferFailed = &t
		return
	case 4: // Base
		t := true
		ty.NotOwner = &t
		return
	case 5: // Base
		t := true
		ty.NotEnoughAllowance = &t
		return
	case 6: // Base
		t := true
		ty.NotEnoughBalance = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}
func (ty *Error) Error() string {
	if ty.SetCodeFailed != nil {
		return "SetCodeFailed"
	}

	if ty.MustCallByCloudContract != nil {
		return "MustCallByCloudContract"
	}

	if ty.InsufficientBalance != nil {
		return "InsufficientBalance"
	}

	if ty.TransferFailed != nil {
		return "TransferFailed"
	}

	if ty.NotOwner != nil {
		return "NotOwner"
	}

	if ty.NotEnoughAllowance != nil {
		return "NotEnoughAllowance"
	}

	if ty.NotEnoughBalance != nil {
		return "NotEnoughBalance"
	}
	return "Unknown"
}
