package sidechain

import (
	"encoding/hex"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets"
)

// ContractQuery 只读合约查询。contract=dao 时 method 为查询方法名，args 为字符串参数列表（与 ExecuteQuery 约定一致）。
func (app *SideChain) ContractQuery(caller []byte, contract, method string, args []string) ([]byte, error) {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := pallets.NewRuntime(height, txn, caller)
	argsBytes := make([][]byte, len(args))
	for i, arg := range args {
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("dao: decode arg: %w", err)
		}
		argsBytes[i] = argBytes
	}

	return pallets.Query(&model.ContractCall{
		Contract: contract,
		Method:   method,
		Args:     argsBytes,
	}, runtime)
}

func (app *SideChain) ContractMutation(caller []byte, call *model.ContractCall) error {
	height := app.state.Height
	txn := app.onGoingBlock

	runtime := pallets.NewRuntime(height, txn, caller)
	return pallets.Mutation(call, runtime)
}
