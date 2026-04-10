package sidechain

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts"
)

func (app *SideChain) ContractMutation(caller []byte, callerType uint32, txn *model.Txn, height int64, call *model.ContractCall) error {
	address := model.UniAddr{
		T: callerType,
		V: caller,
	}

	runtime := contracts.NewRuntime(height, txn, address, model.UniAddr{
		V: app.dkg.DkgPubKey.Byte(),
	})

	return contracts.Mutation(call, runtime)
}

// ContractQuery 只读合约查询。contract=dao 时 method 为查询方法名，args 为字符串参数列表（与 ExecuteQuery 约定一致）。
func (app *SideChain) ContractQuery(caller []byte, callerType uint32, contract string, method string, args [][]byte) ([]byte, error) {
	height := app.state.Height
	txn := model.DBINS.NewTransaction()
	defer txn.Rollback()

	address := model.UniAddr{
		T: callerType,
		V: caller,
	}

	runtime := contracts.NewRuntime(height, txn, address, model.UniAddr{
		V: app.dkg.DkgPubKey.Byte(),
	})

	return contracts.Query(&model.ContractCall{
		Contract: contract,
		Method:   model.MethodToSelector(method),
		Args:     args,
	}, runtime)
}

func (app *SideChain) ContractDryRun(caller []byte, callerType uint32, call *model.ContractCall) error {
	height := app.state.Height
	txn := model.DBINS.NewTransaction()
	defer txn.Rollback()

	address := model.UniAddr{
		T: callerType,
		V: caller,
	}

	runtime := contracts.NewRuntime(height, txn, address, model.UniAddr{
		V: app.dkg.DkgPubKey.Byte(),
	})
	return contracts.Mutation(call, runtime)
}
