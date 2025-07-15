package model

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type PodVersion struct {
	PodId   uint64
	Version uint32
	Status  uint8
}

type Container struct { // Composite
	Image   []byte
	Command Command
	Port    []Service
	Cr      CR
}

type CR struct { // Composite
	Cpu  uint32
	Mem  uint32
	Disk []Disk
	Gpu  uint32
}

type Disk struct { // Composite
	Path DiskClass
	Size uint32
}

type DiskClass struct { // Enum
	SSD *[]byte // 0
}

type Service struct { // Enum
	Tcp        *uint16 // 0
	Udp        *uint16 // 1
	Http       *uint16 // 2
	Https      *uint16 // 3
	ProjectTcp *uint16 // 4
	ProjectUdp *uint16 // 5
}

type Command struct { // Enum
	SH   *[]byte // 0
	BASH *[]byte // 1
	ZSH  *[]byte // 2
	NONE *bool   // 3
}

type Pod struct {
	PodId               uint64
	Containers          []Container
	Owner               types.H160
	Ptype               PodType
	TeeType             TEEType
	Version             uint32
	Status              uint8
	LastMintBlockNumber uint32
	Meta                []byte
}

type PodType struct { // Enum
	CPU    *bool // 0
	GPU    *bool // 1
	SCRIPT *bool // 2
}

type TEEType struct { // Enum
	SGX *bool // 0
	CVM *bool // 1
}
