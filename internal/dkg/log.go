package dkg

import "fmt"

type Logger struct {
	NodeIndex int
}

func (l Logger) Info(keyvals ...any) {
	logs := []any{"NODE ==========", l.NodeIndex, " Info  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"}
	logs = append(logs, keyvals...)
	fmt.Println(logs...)
}
func (l Logger) Error(keyvals ...any) {
	logs := []any{"NODE ==========", l.NodeIndex, " Error  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"}
	logs = append(logs, keyvals...)
	fmt.Println(logs...)
}
