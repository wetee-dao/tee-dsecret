package model

type ContractCall struct {
	Contract string
	Method   string
	Args     [][]byte
}
