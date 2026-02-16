// Package model: StoreMapping 提供类似智能合约 mapping 的键值存储能力。
// 泛型 StoreMapping[K] 的 K 为 key 类型（string、[]byte、uint64、uint32），方法直接使用 key 类型。
package model

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/cockroachdb/pebble"
)

// MappingKey 约束 mapping 的 key 类型：string、[]byte、uint64、uint32。
type MappingKey interface {
	string | []byte | uint64 | uint32
}

// StoreMapping 类似 Solidity mapping(key => value)。
// K 为 key 类型，使用 m.Get(txn, key)、m.Set(txn, key, value) 等时 key 类型即 K。
type StoreMapping[K MappingKey] struct {
	Namespace string // 命名空间，如 "dao"
	KeyPrefix string // 键前缀，如 "member_"、"allowance_"
}

func (m *StoreMapping[K]) fullKey(suffix string) []byte {
	return ComboNamespaceKey(m.Namespace, m.KeyPrefix+suffix)
}

// keySuffix 将 key 编码为存储用的字符串后缀。
func (m *StoreMapping[K]) keySuffix(key K) string {
	switch v := any(key).(type) {
	case string:
		return v
	case []byte:
		return hex.EncodeToString(v)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	default:
		panic("unsupported MappingKey type")
	}
}

// StorageKey 返回该 key 对应的完整存储 key。
func (m *StoreMapping[K]) StorageKey(key K) []byte {
	return m.fullKey(m.keySuffix(key))
}

func (m *StoreMapping[K]) getByFullKey(txn *Txn, fullKey []byte) ([]byte, error) {
	v, err := txn.Get(fullKey)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

// Get 按 key 读取 raw 值；不存在时返回 (nil, nil)。
func (m *StoreMapping[K]) Get(txn *Txn, key K) ([]byte, error) {
	return m.getByFullKey(txn, m.StorageKey(key))
}

// Set 按 key 写入 raw 值。
func (m *StoreMapping[K]) Set(txn *Txn, key K, value []byte) error {
	return txn.Set(m.StorageKey(key), value)
}

// Delete 按 key 删除。
func (m *StoreMapping[K]) Delete(txn *Txn, key K) error {
	return txn.Delete(m.StorageKey(key))
}

// Contains 按 key 判断是否存在。
func (m *StoreMapping[K]) Contains(txn *Txn, key K) (bool, error) {
	v, err := m.Get(txn, key)
	if err != nil {
		return false, err
	}
	return len(v) > 0, nil
}

// DeleteByPrefix 删除该 mapping 下以 prefix 开头的所有 key（慎用，会扫前缀）。
// prefix 为字符串后缀前缀，与 key 同编码方式，如 "0" 表示所有数字 key 中以 0 开头的。
func (m *StoreMapping[K]) DeleteByPrefix(txn *Txn, prefix string) error {
	return txn.DeletekeysByPrefix(m.fullKey(prefix))
}

// GetMappingJson 按 key 读取并 JSON 反序列化为 V；不存在时返回 (nil, nil)。
func GetMappingJson[K MappingKey, V any](m *StoreMapping[K], txn *Txn, key K) (*V, error) {
	return TxnGetJson[V](txn, m.StorageKey(key))
}

// SetMappingJson 将 v 序列化为 JSON 并按 key 写入。
func SetMappingJson[K MappingKey, V any](m *StoreMapping[K], txn *Txn, key K, v *V) error {
	return TxnSetJson(txn, m.StorageKey(key), v)
}
