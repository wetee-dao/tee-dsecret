package model

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

// VerifyTxSigner 验证交易发起方：caller 对 Tx（不含 signature）的签名必须有效
func VerifyTxSigner(tx *Tx) (*SysCall, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}

	caller := tx.GetCaller()
	sig := tx.GetSignature()
	if len(caller) == 0 {
		return nil, errors.New("tx: missing caller")
	}
	if len(sig) == 0 {
		return nil, errors.New("tx: missing signature")
	}

	// 处理系统请求
	call := new(SysCall)
	err := protoio.ReadMessage(bytes.NewBuffer(tx.Call), call)
	if err != nil {
		return nil, err
	}

	switch tx.CallerType {
	case 0: // node ed25519
		if !SignVerify(caller, call.BytesForSig(), sig) {
			return nil, errors.New("tx: invalid signature")
		}
	case 1: // side contract
		if err := SideContractPolkadotVerify(caller, call.BytesForSig(), sig); err != nil {
			return nil, errors.New("tx: contrtact invalid signature")
		}
	default:
		return nil, errors.New("unknown caller type")
	}

	return call, nil
}

func SideContractVerify(caller []byte, callerType uint32, call *ContractCall, signature []byte) error {
	bt, err := call.XXX_Marshal(nil, true)
	if err != nil {
		return err
	}

	switch callerType {
	case 1:
		return SideContractPolkadotVerify(caller, bt, signature)
	default:
		return errors.New("unknown caller type")
	}
}

func SideContractPolkadotVerify(caller []byte, call []byte, signature []byte) error {
	// 解析公钥
	pubkey, err := sr25519.Scheme{}.FromPublicKey(caller)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 验证签名
	ok := pubkey.Verify(call, signature)
	if !ok {
		return errors.New("invalid signature")
	}

	return nil
}

func SideContractPolkadotSign(caller *ink.Signer, call *ContractCall) ([]byte, error) {
	parts := append([][]byte{[]byte("<Bytes>"), call.Name, call.Method[:]}, call.Args...)
	parts = append(parts, []byte("</Bytes>"))

	// 合并所有部分
	result := bytes.Join(parts, nil)

	return caller.Sign(result)
}
