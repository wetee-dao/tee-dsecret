// stategen 对指定 struct 的每个字段生成独立的「按字段存取」状态函数（get/set），
// 使用 model.Txn，每字段一个 key，避免整块 state 写放大。
// 用法: go run ./tools/stategen -pkg=sidechain -file=./side-chain/dao.go -struct=daoState -namespace=dao -prefix=state -out=./side-chain/dao_state_gen.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	pkg := flag.String("pkg", "", "生成代码的 package 名")
	file := flag.String("file", "", "源 Go 文件路径（含待解析的 struct）")
	structName := flag.String("struct", "", "要生成存取代码的 struct 名")
	namespace := flag.String("namespace", "", "存储 namespace（与 model.Txn 一致）")
	prefix := flag.String("prefix", "", "key 前缀，最终 key = namespace + \"_\" + prefix + \"_\" + fieldKey")
	out := flag.String("out", "", "输出 .go 文件路径")
	configPath := flag.String("config", "", "可选：YAML 配置文件路径，可含多组 pkg/file/struct/namespace/prefix/out")
	flag.Parse()

	if *configPath != "" {
		if err := runConfig(*configPath); err != nil {
			log.Fatalf("stategen: %v", err)
		}
		return
	}

	if *pkg == "" || *file == "" || *structName == "" || *namespace == "" || *out == "" {
		fmt.Fprintf(os.Stderr, "用法:\n  stategen -pkg=<pkg> -file=<go文件> -struct=<结构体名> -namespace=<ns> -prefix=<前缀> -out=<输出.go>\n")
		fmt.Fprintf(os.Stderr, "或:\n  stategen -config=<spec.yaml>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := Config{
		Pkg:       *pkg,
		File:      *file,
		Struct:    *structName,
		Namespace: *namespace,
		Prefix:    *prefix,
		Out:       *out,
	}
	if err := Run(cfg); err != nil {
		log.Fatalf("stategen: %v", err)
	}
	log.Printf("wrote %s", cfg.Out)
}

type configFile struct {
	Gen []Config `yaml:"gen"`
}

func runConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	for _, c := range cf.Gen {
		if err := Run(c); err != nil {
			return fmt.Errorf("%s: %w", c.Out, err)
		}
		log.Printf("wrote %s", c.Out)
	}
	return nil
}
