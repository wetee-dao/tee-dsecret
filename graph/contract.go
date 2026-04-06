package graph

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
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

// SubmitContractCall 提交合约调用
func SubmitContractCall(caller []byte, contract, method string, args []string, signatureType uint32, signature string) error {
	argsBytes := make([][]byte, len(args))
	for i, arg := range args {
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return fmt.Errorf("dao: decode arg: %w", err)
		}
		argsBytes[i] = argBytes
	}

	contractPayload, err := codec.Encode(model.ContractCall{
		Contract: contract,
		Method:   method,
		Args:     argsBytes,
	})
	if err != nil {
		return fmt.Errorf("dao: encode contract payload: %w", err)
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("signature 需为 hex 编码的 32 字节: %w", err)
	}

	tx := &model.Tx{
		Caller:        caller,
		Payload:       &model.Tx_Contract{Contract: contractPayload},
		SignatureType: signatureType,
		Signature:     signatureBytes,
	}

	_, err = sidechain.SubmitTx(tx)
	return err
}

// ContractQuery 只读查询
func ContractDryRun(caller []byte, contract string, mut bool, method string, args []string) (string, error) {
	argsBytes := make([][]byte, len(args))
	for i, arg := range args {
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return "", fmt.Errorf("dao: decode arg: %w", err)
		}
		argsBytes[i] = argBytes
	}

	var err error
	var out []byte
	if !mut {
		out, err = sideChain.ContractQuery(caller, contract, method, argsBytes)
		if err != nil {
			return "", err
		}
	} else {
		err = sideChain.ContractDryRun(caller, &model.ContractCall{
			Contract: contract,
			Method:   method,
			Args:     argsBytes,
		})
	}

	return base64.StdEncoding.EncodeToString(out), nil
}
