package cloud

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Pod struct { // Composite
	Name       []byte
	Owner      types.H160
	Contract   PodRef
	Ptype      PodType
	StartBlock uint32
	TeeType    TEEType
}
type PodRef struct { // Composite
	Inner CallBuilder
}
type CallBuilder struct { // Composite
	Addr types.H160
}
type PodType struct { // Enum
	CpuService *bool // 0
	GpuService *bool // 1
	Script     *bool // 2
}

func (ty PodType) Encode(encoder scale.Encoder) (err error) {
	if ty.CpuService != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.GpuService != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Script != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *PodType) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.CpuService = &t
		return
	case 1: // Base
		t := true
		ty.GpuService = &t
		return
	case 2: // Base
		t := true
		ty.Script = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type TEEType struct { // Enum
	SGX *bool // 0
	CVM *bool // 1
}

func (ty TEEType) Encode(encoder scale.Encoder) (err error) {
	if ty.SGX != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.CVM != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *TEEType) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.SGX = &t
		return
	case 1: // Base
		t := true
		ty.CVM = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Service struct { // Enum
	Tcp        *uint16 // 0
	Udp        *uint16 // 1
	Http       *uint16 // 2
	Https      *uint16 // 3
	ProjectTcp *uint16 // 4
	ProjectUdp *uint16 // 5
}

func (ty Service) Encode(encoder scale.Encoder) (err error) {
	if ty.Tcp != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Tcp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Udp != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Udp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Http != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Http)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Https != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Https)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProjectTcp != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ProjectTcp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProjectUdp != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ProjectUdp)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Service) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Inline
		ty.Tcp = new(uint16)
		err = decoder.Decode(ty.Tcp)
		if err != nil {
			return err
		}
		return
	case 1: // Inline
		ty.Udp = new(uint16)
		err = decoder.Decode(ty.Udp)
		if err != nil {
			return err
		}
		return
	case 2: // Inline
		ty.Http = new(uint16)
		err = decoder.Decode(ty.Http)
		if err != nil {
			return err
		}
		return
	case 3: // Inline
		ty.Https = new(uint16)
		err = decoder.Decode(ty.Https)
		if err != nil {
			return err
		}
		return
	case 4: // Inline
		ty.ProjectTcp = new(uint16)
		err = decoder.Decode(ty.ProjectTcp)
		if err != nil {
			return err
		}
		return
	case 5: // Inline
		ty.ProjectUdp = new(uint16)
		err = decoder.Decode(ty.ProjectUdp)
		if err != nil {
			return err
		}
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Disk struct { // Composite
	Path DiskClass
	Size uint32
}
type DiskClass struct { // Enum
	SSD *[]byte // 0
}

func (ty DiskClass) Encode(encoder scale.Encoder) (err error) {
	if ty.SSD != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.SSD)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *DiskClass) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Inline
		ty.SSD = new([]byte)
		err = decoder.Decode(ty.SSD)
		if err != nil {
			return err
		}
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Container struct { // Composite
	Name    []byte
	Image   []byte
	Command Command
	Port    []Service
	Cr      CR
}
type Command struct { // Enum
	SH   *[]byte // 0
	BASH *[]byte // 1
	ZSH  *[]byte // 2
	NONE *bool   // 3
}

func (ty Command) Encode(encoder scale.Encoder) (err error) {
	if ty.SH != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.SH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.BASH != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.BASH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ZSH != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ZSH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NONE != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Command) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Inline
		ty.SH = new([]byte)
		err = decoder.Decode(ty.SH)
		if err != nil {
			return err
		}
		return
	case 1: // Inline
		ty.BASH = new([]byte)
		err = decoder.Decode(ty.BASH)
		if err != nil {
			return err
		}
		return
	case 2: // Inline
		ty.ZSH = new([]byte)
		err = decoder.Decode(ty.ZSH)
		if err != nil {
			return err
		}
		return
	case 3: // Base
		t := true
		ty.NONE = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type CR struct { // Composite
	Cpu  uint32
	Mem  uint32
	Disk []Disk
	Gpu  uint32
}
type Error struct { // Enum
	SetCodeFailed         *bool // 0
	MustCallByGovContract *bool // 1
	WorkerNotFound        *bool // 2
	WorkerLevelNotEnough  *bool // 3
	RegionNotMatch        *bool // 4
	WorkerNotOnline       *bool // 5
}

func (ty Error) Encode(encoder scale.Encoder) (err error) {
	if ty.SetCodeFailed != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MustCallByGovContract != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerNotFound != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerLevelNotEnough != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.RegionNotMatch != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerNotOnline != nil {
		err = encoder.PushByte(5)
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
		ty.MustCallByGovContract = &t
		return
	case 2: // Base
		t := true
		ty.WorkerNotFound = &t
		return
	case 3: // Base
		t := true
		ty.WorkerLevelNotEnough = &t
		return
	case 4: // Base
		t := true
		ty.RegionNotMatch = &t
		return
	case 5: // Base
		t := true
		ty.WorkerNotOnline = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}
func (ty *Error) Error() string {
	if ty.SetCodeFailed != nil {
		return "SetCodeFailed"
	}

	if ty.MustCallByGovContract != nil {
		return "MustCallByGovContract"
	}

	if ty.WorkerNotFound != nil {
		return "WorkerNotFound"
	}

	if ty.WorkerLevelNotEnough != nil {
		return "WorkerLevelNotEnough"
	}

	if ty.RegionNotMatch != nil {
		return "RegionNotMatch"
	}

	if ty.WorkerNotOnline != nil {
		return "WorkerNotOnline"
	}
	return "Unknown"
}

type Tuple_87 struct { // Tuple
	F0 uint64
	F1 Pod
	F2 []Container
}
type Tuple_91 struct { // Tuple
	F0 uint64
	F1 uint32
	F2 byte
}
type Tuple_94 struct { // Tuple
	F0 Pod
	F1 []Container
	F2 uint32
	F3 byte
}
type Tuple_98 struct { // Tuple
	F0 uint64
	F1 Pod
	F2 []Container
	F3 uint32
	F4 byte
}
