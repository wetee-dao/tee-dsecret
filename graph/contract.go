package graph

import (
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

func decodeArgs(args []string) ([][]byte, error) {
	argsBytes := make([][]byte, len(args))
	for i, arg := range args {
		arg = strings.TrimPrefix(arg, "0x")
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("decodeArgs: decode arg: %w", err)
		}
		argsBytes[i] = argBytes
	}

	return argsBytes, nil
}

// SubmitContractCall 提交合约调用
func SubmitContractCall(caller []byte, callerType uint32, contract string, method [4]byte, argsBytes [][]byte, signature []byte) error {
	contractPayload, err := codec.Encode(model.ContractCall{
		Contract: contract,
		Method:   method,
		Args:     argsBytes,
	})
	if err != nil {
		return fmt.Errorf("SubmitContractCall: encode contract payload: %w", err)
	}

	tx := &model.Tx{
		Caller:     caller,
		CallerType: callerType,
		Payload:    &model.Tx_Contract{Contract: contractPayload},
		Signature:  signature,
	}

	_, err = sidechain.SubmitTx(tx)
	return err
}

// ContractQuery 只读查询
func ContractDryRun(caller []byte, callerType int, contract string, mut bool, method string, args [][]byte) (string, error) {
	var err error
	var out []byte
	if !mut {
		out, err = sideChain.ContractQuery(caller, uint32(callerType), contract, method, args)
		if err != nil {
			return "", err
		}
	} else {
		err = sideChain.ContractDryRun(caller, &model.ContractCall{
			Contract: contract,
			Method:   model.MethodToSelector(method),
			Args:     args,
		})
	}

	return hex.EncodeToString(out), nil
}
