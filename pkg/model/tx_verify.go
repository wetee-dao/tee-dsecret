package model

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/gogo/protobuf/proto"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

// TxBytesForSigning 返回用于签名的 Tx 序列化结果（不包含 signature 字段，保证验签时可复现）
func TxBytesForSigning(tx *Tx) ([]byte, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	// 复制一份并清空 signature，再序列化
	signTx := proto.Clone(tx).(*Tx)
	signTx.Signature = nil
	return proto.Marshal(signTx)
}

// VerifyTxSigner 验证交易发起方：caller 对 Tx（不含 signature）的签名必须有效
func VerifyTxSigner(tx *Tx) error {
	if tx == nil {
		return errors.New("tx is nil")
	}
	caller := tx.GetCaller()
	sig := tx.GetSignature()
	if len(caller) == 0 {
		return errors.New("tx: missing caller")
	}
	if len(sig) == 0 {
		return errors.New("tx: missing signature")
	}

	switch tx.CallerType {
	case 0:
		msg, err := TxBytesForSigning(tx)
		if err != nil {
			return err
		}
		if !SignVerify(caller, msg, sig) {
			return errors.New("tx: invalid signature")
		}
	case 1:
		contractBt := tx.GetContract()
		call := &ContractCall{}
		codec.Decode(contractBt, call)
		if err := SideContractPolkadotVerify(caller, call, sig); err != nil {
			return errors.New("tx: contrtact invalid signature")
		}
	default:
		return errors.New("unknown caller type")
	}

	return nil
}

func SideContractVerify(caller []byte, callerType uint32, call *ContractCall, signature []byte) error {
	switch callerType {
	case 1:
		return SideContractPolkadotVerify(caller, call, signature)
	default:
		return errors.New("unknown caller type")
	}
}

func SideContractPolkadotVerify(caller []byte, call *ContractCall, signature []byte) error {
	// 解析公钥
	pubkey, err := sr25519.Scheme{}.FromPublicKey(caller)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	parts := append([][]byte{[]byte("<Bytes>"), []byte(call.Contract), call.Method[:]}, call.Args...)
	parts = append(parts, []byte("</Bytes>"))

	// 合并所有部分
	result := bytes.Join(parts, nil)

	// 验证签名
	ok := pubkey.Verify(result, signature)
	if !ok {
		return errors.New("invalid signature")
	}

	return nil
}
