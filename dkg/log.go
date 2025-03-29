package dkg

import "fmt"

type Logger struct {
	NodeIndex int
}

func (l Logger) Info(keyvals ...any) {
	fmt.Println("NODE ==========", l.NodeIndex, " INFO  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(keyvals...)
}
func (l Logger) Error(keyvals ...any) {
	fmt.Println("NODE ==========", l.NodeIndex, " Error  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(keyvals...)
}
