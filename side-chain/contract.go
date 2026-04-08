package sidechain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts"
)

func (app *SideChain) ContractMutation(caller []byte, callerType uint32, call *model.ContractCall) error {
	height := app.state.Height
	txn := app.onGoingBlock

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

func (app *SideChain) ContractDryRun(caller []byte, call *model.ContractCall) error {
	height := app.state.Height
	txn := model.DBINS.NewTransaction()
	defer txn.Rollback()

	address := new(model.UniAddr)
	err := codec.Decode(caller, &address)
	if err != nil {
		return err
	}

	runtime := contracts.NewRuntime(height, txn, *address, model.UniAddr{
		V: app.dkg.DkgPubKey.Byte(),
	})
	return contracts.Mutation(call, runtime)
}
