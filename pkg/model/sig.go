package model

import (
	"bytes"
)

func (c *SysCall) BytesForSig() []byte {
	if call := c.GetContract(); call != nil {
		return ContractForSig(call)
	}

	data, err := c.XXX_Marshal(nil, true)
	if err != nil {
		panic(err)
	}
	return data

}

func ContractForSig(call *ContractCall) []byte {
	parts := append([][]byte{[]byte("<Bytes>"), call.Name, call.Method[:]}, call.Args...)
	parts = append(parts, []byte("</Bytes>"))
	return bytes.Join(parts, nil)
}
