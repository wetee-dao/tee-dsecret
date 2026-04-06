// Package pallets 定义侧链「合约/Pallet」的标准调用接口（仅抽象结构体与函数），
// 各 pallet 实现该接口并导出实例，由调用方直接使用，无需注册表。
package contracts

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts/dao"
)

type Contract interface {
	Init()
}

func Query(call *model.ContractCall, runtime model.ContractApi) ([]byte, error) {
	switch call.Contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		query := dao.DaoQuery{DAO: *ins}
		return query.ExecQuery(call)
	default:
		return nil, fmt.Errorf("unsupported contract: %s", call.Contract)
	}
}

func Mutation(call *model.ContractCall, runtime model.ContractApi) error {
	switch call.Contract {
	case "dao":
		ins := dao.NewDAO(runtime)
		mutation := dao.DaoMutation{DAO: *ins}
		return mutation.ExecCall(call)
	default:
		return fmt.Errorf("unsupported contract: %s", call.Contract)
	}
}
