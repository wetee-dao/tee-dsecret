package model

import (
	"encoding/binary"
	"fmt"

	"github.com/cockroachdb/pebble"
)

func AddToList(namespace string, skey string, val []byte) error {
	indexBt, err := GetKey(namespace, "__"+skey)
	if err != nil {
		return err
	}
	index := unmarshalSize(indexBt)

	txn := DBINS.NewTransaction()
	txn.SetKey(namespace, "__"+skey, marshalSize(index+1))
	txn.SetKey(namespace, skey+fmt.Sprint(index), val)
	return txn.Commit()
}

func GetList(namespace string, subkey string, page int, size int) ([][]byte, error) {
	key := []byte(namespace + "_" + subkey)
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: key,
		UpperBound: keyUpperBound(key),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var list [][]byte
	for iter.First(); iter.Valid(); iter.Next() {
		v := iter.Value()
		value, err := Unseal(v, nil)
		if err != nil {
			return nil, err
		}

		list = append(list, value)
	}

	return list, nil
}

func DeleteList(namespace string, key string) error {
	return DeletekeysByPrefix(namespace, key)
}

func marshalSize(size uint64) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint64(bs, size)
	return bs
}

func unmarshalSize(bz []byte) uint64 {
	return binary.LittleEndian.Uint64(bz)
}
