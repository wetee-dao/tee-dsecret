package dkg

import "fmt"

type Logger struct {
	NodeIndex int
}

func (l Logger) Info(keyvals ...any) {
	logs := []any{"NODE ==========", l.NodeIndex, " Info  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"}
	for _, v := range keyvals {
		logs = append(logs, v)
	}
	fmt.Println(logs...)
}
func (l Logger) Error(keyvals ...any) {
	logs := []any{"NODE ==========", l.NodeIndex, " Error  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"}
	for _, v := range keyvals {
		logs = append(logs, v)
	}
	fmt.Println(logs...)
}
