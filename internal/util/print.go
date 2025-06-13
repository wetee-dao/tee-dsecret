package util

import (
	"encoding/json"
	"fmt"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func LogWithYellow(tag string, a ...any) {
	b := make([]any, 0, len(a)+3)

	b = append(b, Yellow+tag+":")
	b = append(b, Reset)
	b = append(b, a...)

	fmt.Println(b...)
}

func LogWithPurple(tag string, a ...any) {
	b := make([]any, 0, len(a)+3)

	b = append(b, Purple+tag+":")
	b = append(b, Reset)
	b = append(b, a...)

	fmt.Println(b...)
}

func LogWithRed(tag string, a ...any) {
	b := make([]any, 0, len(a)+3)

	b = append(b, Red+tag+":")
	b = append(b, Reset)
	b = append(b, a...)

	fmt.Println(b...)
}

func PrintJson(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println("------------------ json --------------------------")
	fmt.Println(string(b))
	fmt.Println("------------------ json --------------------------")
}
