package dkg

import (
	"fmt"

	"wetee.app/dsecret/internal/util"
)

type Logger struct {
	NodeIndex int
}

func (l Logger) Info(keyvals ...any) {
	logs := []any{"|| "}
	logs = append(logs, keyvals...)
	util.LogWithGray("NODE "+fmt.Sprint(l.NodeIndex), logs...)
}
func (l Logger) Error(keyvals ...any) {
	logs := []any{"|| "}
	logs = append(logs, keyvals...)
	util.LogWithRed("NODE "+fmt.Sprint(l.NodeIndex), logs...)
}
