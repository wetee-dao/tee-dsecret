// Package pallets 定义侧链「合约/Pallet」的标准调用接口（仅抽象结构体与函数），
// 各 pallet 实现该接口并导出实例，由调用方直接使用，无需注册表。
package pallets

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/base"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/dao"
)

// Pallet 是单个 pallet/合约的标准调用接口。
type PalletCall interface {
	// ApplyCall 在给定 txn 与区块高度下执行一次调用：caller 为发起方，payload 为合约自定义载荷。
	ExecCall(tx *model.Tx, runtime base.ContractApi) error
}

func Query(call *model.ContractCall, runtime base.ContractApi) ([]byte, error) {
	switch call.Contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		query := dao.DaoQuery{DAO: *ins}
		return query.ExecQuery(call)
	default:
		return nil, fmt.Errorf("unsupported contract: %s", call.Contract)
	}
}

func Mutation(call *model.ContractCall, runtime base.ContractApi) error {
	switch call.Contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		mutation := dao.DaoMutation{DAO: *ins}
		return mutation.ExecCall(call)
	default:
		return fmt.Errorf("unsupported contract: %s", call.Contract)
	}
}
