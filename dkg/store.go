package dkg

import (
	"context"
	"encoding/json"

	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
)

// SetData 向Peer发布密钥数据
//
// 该函数将密钥数据序列化为JSON格式，并通过Peer的Pub方法发布到指定的主题“secret”
// 使用了context.Background()来执行发布操作，不设置超时或截止时间
//
// 参数:
//
//	datas ([]types.Kvs): 待发布的密钥数据列表
//
// 返回值:
//
//	error: 表示发布操作中可能遇到的错误，如果发布成功则为nil
func (r *DKG) SetData(datas []types.Kvs) error {
	// 将密钥数据列表序列化为JSON格式的字节切片
	bt, _ := json.Marshal(datas)

	// 通过Peer的Pub方法，向“secret”主题发布序列化后的数据
	// 这里使用了context.Background()来执行发布操作，意味着这是一个无截止时间的异步操作
	return r.Peer.Pub(context.Background(), "secret", bt)
}

// GetData 通过给定的键从存储中获取数据
//
// 参数:
//
//	k - 要检索的数据的键
//
// 返回值:
//
//	[]byte - 存储中与给定键关联的数据，如果没有找到则为nil
//	error  - 如果在存储操作中发生错误，返回该错误；否则返回nil
//
// 该函数通过调用store的getKey方法，使用"k"作为前缀和给定的键(k)来组合查询存储中的"secret"数据
func (r *DKG) GetData(k string) ([]byte, error) {
	return store.GetKey("secret", k)
}
