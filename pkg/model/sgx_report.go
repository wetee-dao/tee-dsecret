package model

import (
	"encoding/json"
	"errors"
)

// GetReport 通过给定的哈希值获取 TeeCall 和 TeeReport
// 参数 hash 用于标识要获取报告的数据哈希值
// 返回值包括 TeeCall、TeeReport 和错误类型
// 可能的错误包括数据获取失败、反序列化 TeeCall 失败以及验证报告失败
func GetReport(hash string) (*TeeCall, *TeeVerifyResult, error) {
	// 根据哈希值获取对应的 secretData
	secretData, err := GetKey("secret", hash)
	// 如果获取数据失败，返回错误
	if err != nil {
		return nil, nil, err
	}

	// 初始化 TeeCall 实例
	TeeCall := &TeeCall{}
	// 尝试将 secretData 反序列化为 TeeCall
	err = json.Unmarshal(secretData, TeeCall)
	// 如果反序列化失败，返回错误
	if err != nil {
		return nil, nil, errors.New("unmarshal tee param: " + err.Error())
	}

	// 使用 TeeCall 验证并获取报告
	report, err := VerifyReport(TeeCall)
	// 如果验证报告失败，返回错误
	if err != nil {
		return nil, nil, errors.New("verify cluster report: " + err.Error())
	}

	// 成功获取 TeeCall 和 TeeReport，返回结果
	return TeeCall, report, nil
}
