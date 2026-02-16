package graph

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
)

// DecodeCaller 支持 32 字节 hex（可带 0x 前缀）或 SS58，返回公钥字节。
func DecodeCaller(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	if b, err := hex.DecodeString(s); err == nil && len(b) == 32 {
		return b, nil
	}
	pub, err := model.PubKeyFromSS58(s)
	if err != nil {
		return nil, fmt.Errorf("caller 需为 32 字节 hex 或 SS58: %w", err)
	}
	return pub.Byte(), nil
}

// SubmitContractCall 根据合约名提交交易。dao 的 payload 为 base64 编码的 model.DaoCall protobuf。
func SubmitContractCall(caller []byte, contract string, payload string) error {
	var tx *model.Tx
	switch contract {
	case "dao":
		raw, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return fmt.Errorf("dao payload 需为 base64 编码的 protobuf: %w", err)
		}
		tx = &model.Tx{
			Caller:  caller,
			Payload: &model.Tx_DaoCall{DaoCall: raw},
		}
	default:
		return fmt.Errorf("不支持的合约: %s", contract)
	}
	_, err := sidechain.SubmitTx(tx)
	return err
}

// GetContractState 返回指定合约的状态 JSON，未知合约返回 "{}"。
func GetContractState(contract string) (string, error) {
	switch contract {
	case "dao":
		// TODO: 从侧链只读查询 DAO 状态
		return "{}", nil
	default:
		return "{}", nil
	}
}
