package pallets

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/base"
)

// Runtime 实现 base.ContractApi，在 FinalizeTx / Query 等路径传入高度、txn 与调用方。
type Runtime struct {
	height int64
	txn    *model.Txn
	caller []byte
}

var _ base.ContractApi = Runtime{}

func NewRuntime(height int64, txn *model.Txn, caller []byte) Runtime {
	return Runtime{height: height, txn: txn, caller: caller}
}

func (r Runtime) GetHeight() int64 {
	return r.height
}

func (r Runtime) GetTxn() *model.Txn {
	return r.txn
}

func (r Runtime) GetCaller() []byte {
	return r.caller
}

type Contrtact struct {
	Method string
	Args   [][]byte
}
