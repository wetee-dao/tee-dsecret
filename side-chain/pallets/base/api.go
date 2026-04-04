package base

import "github.com/wetee-dao/tee-dsecret/pkg/model"

// ContractApi 侧链合约运行时：区块高度、读写事务、已解析的调用方（如由 TxBox 回填）。
type ContractApi interface {
	GetHeight() int64
	GetTxn() *model.Txn
	GetCaller() []byte
}

type ContractCall struct {
	contract string
	Txn      *model.Txn
	Caller   []byte
}
