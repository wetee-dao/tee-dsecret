package types

type User struct {
	// User SS58 address
	Address string `json:"address"`
	// User sign time
	Timestamp int64 `json:"timestamp"`
}
