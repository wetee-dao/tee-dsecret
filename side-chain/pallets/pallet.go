// Package pallets 定义侧链「合约/Pallet」的标准调用接口（仅抽象结构体与函数），
// 各 pallet 实现该接口并导出实例，由调用方直接使用，无需注册表。
package pallets

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Pallet 是单个 pallet/合约的标准调用接口。
// 各 pallet 实现此接口并导出变量（如 dao.Pallet），调用方直接对其调用 ApplyCall。
type Pallet interface {
	// ApplyCall 在给定 txn 与区块高度下执行一次调用：caller 为发起方，payload 为合约自定义载荷（如 JSON）。
	ApplyCall(caller, payload []byte, height int64, txn *model.Txn) error
}

// CallFunc 将函数适配为 Pallet，便于将 dao.ApplyDaoCall 等函数作为 Pallet 导出。
type CallFunc func(caller, payload []byte, height int64, txn *model.Txn) error

func (f CallFunc) ApplyCall(caller, payload []byte, height int64, txn *model.Txn) error {
	return f(caller, payload, height, txn)
}
