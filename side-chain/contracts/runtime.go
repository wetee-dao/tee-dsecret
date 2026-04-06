// Package contracts 提供合约运行时环境和合约调用入口。
// 本文件定义了合约执行时的运行时上下文，包含区块高度、交易信息和调用者身份。
package contracts

import (
	"fmt"
	"os"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Runtime 合约运行时上下文，封装了合约执行所需的环境信息。
// 它实现了 model.ContractApi 接口，在 FinalizeTx / Query 等路径中传入区块高度、交易和调用方信息。
// 每次合约调用都会创建一个新的 Runtime 实例，确保执行上下文的隔离性。
type Runtime struct {
	height int64      // 当前区块高度，用于时间相关的业务逻辑判断
	txn    *model.Txn // 当前交易对象，包含交易的签名和原始数据
	caller []byte     // 调用者地址，标识合约方法的调用发起方
	sudo   []byte
}

// NewRuntime 创建一个新的合约运行时实例。
// 参数:
//   - height: 当前区块高度
//   - txn: 交易对象，包含交易的完整信息
//   - caller: 调用者的区块链地址
//
// 返回:
//   - Runtime: 初始化后的运行时上下文对象
func NewRuntime(height int64, txn *model.Txn, caller []byte, sudo []byte) *Runtime {
	return &Runtime{height: height, txn: txn, caller: caller}
}

// GetHeight 获取当前区块高度。
// 该方法用于合约内部判断区块时间或执行与区块高度相关的业务逻辑。
// 返回:
//   - int64: 当前区块高度
func (r *Runtime) GetHeight() int64 {
	return r.height
}

// GetTxn 获取当前交易对象。
// 交易对象包含签名数据和原始交易数据，可用于验证交易有效性。
// 返回:
//   - *model.Txn: 当前交易对象的指针
func (r *Runtime) GetTxn() *model.Txn {
	return r.txn
}

// GetCaller 获取合约调用者的地址。
// 该地址用于权限验证和身份识别，确保只有授权的调用者能执行特定操作。
// 返回:
//   - []byte: 调用者的区块链地址字节切片
func (r *Runtime) GetCaller() []byte {
	return r.caller
}

func (r *Runtime) GetSudoAccount() []byte {
	return r.sudo
}

// is contract inited
func IsContractIsInit(txn *model.Txn, contract string) bool {
	b, err := model.TxnGetCodec[bool](txn, "CONTRACT___INITED", contract)
	if err != nil {
		fmt.Println("contract init db error", err)
		os.Exit(2)
		return false
	}

	if b == nil {
		return false
	}

	return *b
}

func SetContractInited(txn *model.Txn, contract string) error {
	return model.TxnSetCodec(txn, "CONTRACT___INITED", contract, true)
}
