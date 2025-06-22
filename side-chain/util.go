package sidechain

import (
	"fmt"
	"time"
)

const (
	CodeTypeOK              uint32 = 0
	CodeTypeEncodingError   uint32 = 1
	CodeTypeInvalidTxFormat uint32 = 2
	CodeTypeBanned          uint32 = 3
)

func LogWithTime(a ...any) {
	dim := "\033[2m"
	reset := "\033[0m"
	tag := dim + time.Now().Format("06-01-02 15:04:05") + reset
	a = append([]any{tag}, a...)
	fmt.Println(a...)
}
