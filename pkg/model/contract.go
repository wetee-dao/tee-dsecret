package model

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

type ContractCall struct {
	Contract string
	Method   string
	Args     [][]byte
}

// ContractApi 侧链合约运行时：区块高度、读写事务、已解析的调用方（如由 TxBox 回填）。
type ContractApi interface {
	GetHeight() int64
	GetTxn() *Txn
	GetCaller() []byte
}

func RequireArgLen[T any](args []T, n int, method string) error {
	if len(args) < n {
		return fmt.Errorf("dao: %s expects at least %d arg(s), got %d", method, n, len(args))
	}
	return nil
}

// hexToBytes 解析十六进制（可带 0x），不允许 SS58；空串返回 nil 字节切片。
func hexToBytes(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}

// DecodeScaleArg 将单个参数视为 hex( SCALE(T) )，解码为 T（泛型入口，供合约 dispatch / 查询共用）。
func DecodeScaleArg[T any](s string) (T, error) {
	var z T
	raw, err := hexToBytes(s)
	if err != nil {
		return z, fmt.Errorf("hex: %w", err)
	}
	if len(raw) == 0 {
		return z, fmt.Errorf("dao: empty arg after hex decode")
	}
	return DecodeScaleArgBytes[T](raw)
}

func DecodeScaleArgBytes[T any](s []byte) (T, error) {
	var z T
	if err := codec.Decode(s, &z); err != nil {
		return z, fmt.Errorf("scale: %w", err)
	}
	return z, nil
}
