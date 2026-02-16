package model

import (
	"errors"

	"github.com/gogo/protobuf/proto"
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
	msg, err := TxBytesForSigning(tx)
	if err != nil {
		return err
	}
	if !SignVerify(caller, msg, sig) {
		return errors.New("tx: invalid signature")
	}
	return nil
}
