package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// FieldSpec 描述从 AST 解析出的一个结构体字段
type FieldSpec struct {
	Name       string // Go 字段名
	Key        string // 存储 key（snake_case 或 json tag）
	Type       string // 类型字符串
	UseJSON    bool   // 是否用 JSON 序列化（复杂类型）
	Pointer    bool   // 是否指针类型
	Optional   bool   // 是否可空（指针或 omitempty）
	IsMapping  bool   // 是否为 *model.StoreMapping[K]
	MappingKey string // StoreMapping 的 key 类型：uint32、uint64、string 等
}

// Config 单次生成的配置
type Config struct {
	Pkg       string // 生成代码的 package
	File      string // 源 Go 文件路径
	Struct    string // 要处理的 struct 名
	Namespace string // 存储 namespace
	Prefix    string // key 前缀，最终 key = namespace + "_" + prefix + "_" + fieldKey
	Out       string // 输出文件路径
}

// Run 根据配置解析 file 中的 struct，生成按字段存取代码并写入 out
func Run(cfg Config) error {
	src, err := os.ReadFile(cfg.File)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Base(cfg.File), src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	specs, err := extractStructFields(f, cfg.Struct)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("struct %s not found or has no fields", cfg.Struct)
	}
	out := generate(cfg.Pkg, cfg.Struct, cfg.Namespace, cfg.Prefix, specs)
	return os.WriteFile(cfg.Out, out, 0644)
}

func extractStructFields(f *ast.File, structName string) ([]FieldSpec, error) {
	var specs []FieldSpec
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name == nil || ts.Name.Name != structName {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok || st.Fields == nil {
			return true
		}
		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			name := field.Names[0].Name
			typeStr := typeString(field.Type)
			key := name
			useJSON := false
			pointer := false
			optional := false
			if field.Tag != nil {
				tag := strings.Trim(field.Tag.Value, "`")
				for _, part := range strings.Split(tag, " ") {
					if strings.HasPrefix(part, "json:") {
						v := strings.TrimPrefix(part, "json:")
						v = strings.Trim(v, `"`)
						if idx := strings.Index(v, ","); idx >= 0 {
							key = v[:idx]
							if strings.Contains(v[idx:], "omitempty") {
								optional = true
							}
						} else {
							key = v
						}
						break
					}
					if strings.HasPrefix(part, "state:") {
						v := strings.Trim(part, `state:"`)
						for _, kv := range strings.Split(v, ",") {
							if strings.HasPrefix(kv, "key=") {
								key = strings.TrimPrefix(kv, "key=")
								break
							}
						}
						break
					}
				}
			}
			if key == "" || key == "-" {
				key = name // json:"-" 时用字段名作为存储 key
			}
			key = toSnake(key)
			isMapping, mappingKey := isStoreMapping(typeStr)
			if isMapping {
				// StoreMapping 的 key 用 state tag 或字段名转 snake，通常带后缀如 track_
				specs = append(specs, FieldSpec{
					Name:       name,
					Key:        key,
					Type:       typeStr,
					IsMapping:  true,
					MappingKey: mappingKey,
				})
				continue
			}
			_, isSimple := simpleTypes[typeStr]
			if !isSimple || strings.Contains(typeStr, "[]") || strings.Contains(typeStr, "map") {
				useJSON = true
			}
			if strings.HasPrefix(typeStr, "*") {
				pointer = true
			}
			specs = append(specs, FieldSpec{
				Name:     name,
				Key:      key,
				Type:     typeStr,
				UseJSON:  useJSON,
				Pointer:  pointer,
				Optional: optional,
			})
		}
		return false
	})
	return specs, nil
}

var simpleTypes = map[string]bool{
	"string": true, "bool": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"[]byte": true, "byte": true,
}

func typeString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.ArrayType:
		return "[]" + typeString(t.Elt)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + typeString(t.Key) + "]" + typeString(t.Value)
	case *ast.IndexExpr:
		return typeString(t.X) + "[" + typeString(t.Index) + "]"
	case *ast.IndexListExpr:
		var parts []string
		for _, idx := range t.Indices {
			parts = append(parts, typeString(idx))
		}
		return typeString(t.X) + "[" + strings.Join(parts, ", ") + "]"
	default:
		return "interface{}"
	}
}

// isStoreMapping 判断类型是否为 *model.StoreMapping[K]，并返回 key 类型 K。
func isStoreMapping(typeStr string) (ok bool, keyType string) {
	const prefix = "*model.StoreMapping["
	if !strings.HasPrefix(typeStr, prefix) {
		return false, ""
	}
	rest := strings.TrimSuffix(typeStr[len(prefix):], "]")
	keyType = strings.TrimSpace(rest)
	switch keyType {
	case "uint32", "uint64", "string", "[]byte":
		return true, keyType
	default:
		return false, ""
	}
}

func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			r = r - 'A' + 'a'
		}
		b.WriteRune(r)
	}
	return b.String()
}

func generate(pkg, structName, namespace, prefix string, specs []FieldSpec) []byte {
	// 存储 key：标量 = namespace + "_" + (prefix + "_" + fieldKey)；StoreMapping 的 KeyPrefix = field key 字面量
	subKeyPrefix := prefix
	if subKeyPrefix != "" {
		subKeyPrefix = prefix + "_"
	}
	// 访问器类型与变量名小写：daoStateStorage -> daoStateStorageState
	accessorType := toLowerFirst(structName) + "State"
	keyConstPrefix := "stateKey" + structName

	var scalarSpecs, mappingSpecs []FieldSpec
	for _, s := range specs {
		if s.IsMapping {
			mappingSpecs = append(mappingSpecs, s)
		} else {
			scalarSpecs = append(scalarSpecs, s)
		}
	}

	needBigImport := false
	for _, s := range scalarSpecs {
		if s.Type == "*big.Int" {
			needBigImport = true
			break
		}
	}

	var buf bytes.Buffer
	buf.WriteString("// Code generated by tools/stategen. DO NOT EDIT.\n")
	buf.WriteString("// 标量：xxx.Members() / xxx.SetMembers(v)；StoreMapping 通过 DaoStorage 按 key 存取。\n\n")
	buf.WriteString("package " + pkg + "\n\n")
	imports := "\t\"strconv\"\n\n\t\"github.com/wetee-dao/tee-dsecret/pkg/model\"\n"
	if needBigImport {
		imports = "\t\"math/big\"\n\n\t\"strconv\"\n\n\t\"github.com/wetee-dao/tee-dsecret/pkg/model\"\n"
	}
	buf.WriteString("import (\n" + imports + ")\n\n")

	// StoreMapping 的 key 常量（先输出，供下方 state 构造器使用）
	if len(mappingSpecs) > 0 {
		for _, s := range mappingSpecs {
			keyConst := keyConstPrefix + s.Name
			keyVal := s.Key
			buf.WriteString("const " + keyConst + " = \"" + keyVal + "\"\n\n")
		}
	}

	// 访问器结构体：txn + StoreMapping 字段（state.Tracks.Get(txn, id) 统一入口）
	buf.WriteString("// " + accessorType + " 按字段读写的 state 访问器，标量用 Methods，mapping 用 state.Tracks.Get(txn, id) 等。\n")
	buf.WriteString("type " + accessorType + " struct {\n\ttxn *model.Txn\n")
	if len(mappingSpecs) > 0 {
		for _, s := range mappingSpecs {
			buf.WriteString("\t" + s.Name + " *model.StoreMapping[" + s.MappingKey + "]\n")
		}
	}
	buf.WriteString("}\n\n")
	buf.WriteString("// new" + toUpperFirst(accessorType) + " 返回访问器，用法: state := new" + toUpperFirst(accessorType) + "(txn); state.Members(); state.Tracks.Get(txn, id)\n")
	buf.WriteString("func new" + toUpperFirst(accessorType) + "(txn *model.Txn) *" + accessorType + " {\n")
	buf.WriteString("\ts := &" + accessorType + "{txn: txn}\n")
	if len(mappingSpecs) > 0 {
		for _, s := range mappingSpecs {
			keyConst := keyConstPrefix + s.Name
			buf.WriteString("\ts." + s.Name + " = &model.StoreMapping[" + s.MappingKey + "]{Namespace: \"" + namespace + "\", KeyPrefix: " + keyConst + "}\n")
		}
	}
	buf.WriteString("\treturn s\n}\n\n")

	// 标量字段：key 常量 + get/set
	for _, s := range scalarSpecs {
		keyConst := keyConstPrefix + s.Name
		keyVal := subKeyPrefix + s.Key
		buf.WriteString("const " + keyConst + " = \"" + keyVal + "\"\n\n")

		// *big.Int 按 u128 字节存取（model.BytesToU128 / model.U128ToBytes）
		if s.Type == "*big.Int" {
			buf.WriteString("func (s *" + accessorType + ") " + s.Name + "() *big.Int {\n")
			buf.WriteString("\tb, _ := s.txn.Get(model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "))\n\treturn model.BytesToU128(b)\n}\n\n")
			buf.WriteString("func (s *" + accessorType + ") Set" + s.Name + "(v *big.Int) error {\n")
			buf.WriteString("\tif v == nil { v = big.NewInt(0) }\n\treturn s.txn.SetKey(\"" + namespace + "\", " + keyConst + ", model.U128ToBytes(v))\n}\n\n")
			continue
		}

		buf.WriteString("func (s *" + accessorType + ") " + s.Name + "() ")
		if s.UseJSON {
			buf.WriteString(s.Type + " {\n")
			baseType := strings.TrimPrefix(s.Type, "*")
			if s.Pointer {
				buf.WriteString("\tv, _ := model.TxnGetJson[" + baseType + "](s.txn, model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "))\n\treturn v\n}\n\n")
			} else {
				buf.WriteString("\tv, _ := model.TxnGetJson[" + baseType + "](s.txn, model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "))\n\tif v == nil { var z " + s.Type + "; return z }\n\treturn *v\n}\n\n")
			}
		} else {
			buf.WriteString(s.Type + " {\n")
			buf.WriteString("\tb, _ := s.txn.Get(model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "))\n")
			switch s.Type {
			case "string":
				buf.WriteString("\treturn string(b)\n}\n\n")
			case "bool":
				buf.WriteString("\tv, _ := strconv.ParseBool(string(b))\n\treturn v\n}\n\n")
			case "uint32":
				buf.WriteString("\tv, _ := strconv.ParseUint(string(b), 10, 32)\n\treturn uint32(v)\n}\n\n")
			case "uint64":
				buf.WriteString("\tv, _ := strconv.ParseUint(string(b), 10, 64)\n\treturn v\n}\n\n")
			case "int64":
				buf.WriteString("\tv, _ := strconv.ParseInt(string(b), 10, 64)\n\treturn v\n}\n\n")
			case "[]byte":
				buf.WriteString("\treturn b\n}\n\n")
			default:
				buf.WriteString("\treturn " + zeroValue(s.Type) + "\n}\n\n")
			}
		}

		buf.WriteString("func (s *" + accessorType + ") Set" + s.Name + "(v " + s.Type + ") error {\n")
		if s.UseJSON {
			if s.Pointer {
				buf.WriteString("\treturn model.TxnSetJson(s.txn, model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "), v)\n}\n\n")
			} else {
				buf.WriteString("\treturn model.TxnSetJson(s.txn, model.ComboNamespaceKey(\"" + namespace + "\", " + keyConst + "), &v)\n}\n\n")
			}
		} else {
			buf.WriteString("\tvar b []byte\n\tswitch x := any(v).(type) {\n\tcase string:\n\t\tb = []byte(x)\n\tcase bool:\n\t\tb = []byte(strconv.FormatBool(x))\n\tcase uint32:\n\t\tb = []byte(strconv.FormatUint(uint64(x), 10))\n\tcase uint64:\n\t\tb = []byte(strconv.FormatUint(x, 10))\n\tcase int64:\n\t\tb = []byte(strconv.FormatInt(x, 10))\n\tcase []byte:\n\t\tb = x\n\tdefault:\n\t\treturn nil\n\t}\n\treturn s.txn.SetKey(\"" + namespace + "\", " + keyConst + ", b)\n}\n\n")
		}
	}

	return buf.Bytes()
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := rune(s[0])
	if r >= 'A' && r <= 'Z' {
		return string(r-'A'+'a') + s[1:]
	}
	return s
}

func toUpperFirst(s string) string {
	if s == "" {
		return s
	}
	r := rune(s[0])
	if r >= 'a' && r <= 'z' {
		return string(r-'a'+'A') + s[1:]
	}
	return s
}

func zeroValue(t string) string {
	switch t {
	case "string":
		return `""`
	case "bool":
		return "false"
	case "uint32", "uint64", "int64", "int":
		return "0"
	case "[]byte":
		return "nil"
	default:
		return "nil"
	}
}
