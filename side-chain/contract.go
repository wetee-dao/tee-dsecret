package sidechain

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts"
)

// ContractQuery 只读合约查询。contract=dao 时 method 为查询方法名，args 为字符串参数列表（与 ExecuteQuery 约定一致）。
func (app *SideChain) ContractQuery(caller []byte, contract string, method string, args [][]byte) ([]byte, error) {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := contracts.NewRuntime(height, txn, caller)

	return contracts.Query(&model.ContractCall{
		Contract: contract,
		Method:   method,
		Args:     args,
	}, runtime)
}

func (app *SideChain) ContractMutation(caller []byte, call *model.ContractCall) error {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := contracts.NewRuntime(height, txn, caller)
	return contracts.Mutation(call, runtime)
}

func (app *SideChain) ContractDryRun(caller []byte, call *model.ContractCall) error {
	height := app.state.Height
	txn := model.DBINS.NewTransaction()
	defer txn.Rollback()

	runtime := contracts.NewRuntime(height, txn, caller)
	return contracts.Mutation(call, runtime)
}
