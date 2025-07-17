package model

import (
	"encoding/json"
	"errors"
)

// 请求开始TEE容器
type TeeParam struct {
	// sign address
	Address string
	// report time
	Time int64
	// 0: sgx, 1: sev 2: tdx 3: sev-snp
	TeeType uint8
	// report data
	Data []byte
	// report
	Report []byte
}

type TeeReport struct {
	// 0: sgx, 1: sev 2: tdx 3: sev-snp
	TeeType uint8
	// report code signer
	CodeSigner []byte
	// report code signature
	CodeSignature []byte
	// report ProductID
	CodeProductID []byte
}

// GetReport 通过给定的哈希值获取 TeeParam 和 TeeReport
// 参数 hash 用于标识要获取报告的数据哈希值
// 返回值包括 TeeParam、TeeReport 和错误类型
// 可能的错误包括数据获取失败、反序列化 TeeParam 失败以及验证报告失败
func GetReport(hash string) (*TeeParam, *TeeReport, error) {
	// 根据哈希值获取对应的 secretData
	secretData, err := GetKey("secret", hash)
	// 如果获取数据失败，返回错误
	if err != nil {
		return nil, nil, err
	}

	// 初始化 TeeParam 实例
	teeParam := &TeeParam{}
	// 尝试将 secretData 反序列化为 teeParam
	err = json.Unmarshal(secretData, teeParam)
	// 如果反序列化失败，返回错误
	if err != nil {
		return nil, nil, errors.New("unmarshal tee param: " + err.Error())
	}

	// 使用 teeParam 验证并获取报告
	report, err := VerifyReport(teeParam)
	// 如果验证报告失败，返回错误
	if err != nil {
		return nil, nil, errors.New("verify cluster report: " + err.Error())
	}

	// 成功获取 TeeParam 和 TeeReport，返回结果
	return teeParam, report, nil
}
