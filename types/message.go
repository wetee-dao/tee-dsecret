package types

type Message struct {
	Id      uint64 `json:"type"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}
