package types

type Message struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}
