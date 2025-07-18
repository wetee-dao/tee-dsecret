package model

import "go.dedis.ch/kyber/v4/share"

type ReencryptSecretRequest struct {
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
}

type ReencryptedSecretShare struct {
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
	// 密钥碎片索引
	Index int32 `json:"index,omitempty"`
	// Re-encrypted secret share
	// 重新加密的秘密份额
	XncSki []byte `json:"xnc_ski,omitempty"`
	// Random oracle challenge
	// 随机神谕挑战
	Chlgi []byte `json:"chlgi,omitempty"`
	// NIZK proofi of re-encryption
	// NIZK 重新加密证明
	Proofi []byte `json:"proofi,omitempty"`
}

type ReencryptedSecretShareReply struct {
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
	// 密钥碎片
	Share share.PubShare
}

// ReencryptSecret 函数处理重新加密的结果
type ReencryptSecret struct {
	// 密文解码数据，需配合私钥使用
	XncCmt []byte `json:"xnc_cmt,omitempty"`
	// 密文
	EncScrt [][]byte `json:"enc_scrt,omitempty"`
}

// LaunchRequest 函数处理启动请求
type LaunchRequest struct {
	// libos tee report
	Libos *TeeParam
	// cluster tee report
	Cluster *TeeParam
	// worker tee report
	WorkID string
}
