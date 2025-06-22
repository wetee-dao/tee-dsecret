package util

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/wetee-dao/tee-dsecret/internal/model"
)

func Int64ToBytes(time int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(time))
	return b
}

// GetReport 通过给定的哈希值获取 TeeParam 和 TeeReport
// 参数 hash 用于标识要获取报告的数据哈希值
// 返回值包括 TeeParam、TeeReport 和错误类型
// 可能的错误包括数据获取失败、反序列化 TeeParam 失败以及验证报告失败
func GetReport(hash string) (*model.TeeParam, *model.TeeReport, error) {
	// 根据哈希值获取对应的 secretData
	secretData, err := model.GetKey("secret", hash)
	// 如果获取数据失败，返回错误
	if err != nil {
		return nil, nil, err
	}

	// 初始化 TeeParam 实例
	teeParam := &model.TeeParam{}
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
