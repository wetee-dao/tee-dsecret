package model

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

//go:generate protoc --proto_path=. --gogofast_out=. tx.proto

func ToIndexCall(call *types.Call, index int64) *IndexCall {
	bt, _ := codec.Encode(call)
	return &IndexCall{
		Index: index,
		Call:  bt,
	}
}
