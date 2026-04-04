package sidechain

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/dao"
)

// ContractQuery 只读合约查询。contract=dao 时 method 为查询方法名，args 为字符串参数列表（与 ExecuteQuery 约定一致）。
func (app *SideChain) ContractQuery(caller []byte, contract, method string, args []string) ([]byte, error) {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := pallets.NewRuntime(height, txn, caller)

	switch contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		query := dao.DaoQuery{DAO: *ins}
		return query.ExecuteQuery(method, args)
	default:
		return nil, fmt.Errorf("unsupported contract: %s", contract)
	}
}

func (app *SideChain) ContractMutation(caller []byte, call *model.ContractCall) error {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := pallets.NewRuntime(height, txn, caller)

	switch call.Contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		mutation := dao.DaoMutation{DAO: *ins}
		return mutation.ExecCall(call)
	default:
		return fmt.Errorf("unsupported contract: %s", call.Contract)
	}
}
