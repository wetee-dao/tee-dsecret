package types

// 去中心化的机密注入
type Secrets struct {
	Envs  map[string]string
	Files map[string][]byte
}
