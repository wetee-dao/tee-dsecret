package model

import "bytes"

func (c *SysCall) BytesForSig() []byte {
	switch call := any(c).(type) {
	case *ContractCall:
		parts := append([][]byte{[]byte("<Bytes>"), call.Name, call.Method[:]}, call.Args...)
		parts = append(parts, []byte("</Bytes>"))
		return bytes.Join(parts, nil)
	default:
		data, err := c.XXX_Marshal(nil, true)
		if err != nil {
			panic(err)
		}
		return data
	}
}
