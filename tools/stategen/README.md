# stategen：状态按字段存取代码生成器

为任意 struct 的**每个字段**生成独立的 get/set 函数，使存储按「单字段 key」读写，避免「改一个字段就重写整块 state」带来的写放大。

## 设计

- **输入**：Go 源文件 + 结构体名 + namespace/prefix
- **输出**：同包下的 `*_gen.go`，包含：
  - **小写访问器类型**：`type xxxState struct{ txn *model.Txn }`（如 `daoStateState`）
  - **构造器**：`newXxxState(txn)`，返回 `*xxxState`，变量名建议小写如 `state`
  - **读**：`state.Members()`、`state.TotalIssuance()` 等（方法名与字段名一致）
  - **写**：`state.SetMembers(v)`、`state.SetTotalIssuance(v)` 等
- **类型**：
  - 简单类型（string, bool, uint32, uint64, int64, []byte）：直接 `txn.Get`/`txn.SetKey` + strconv
  - `*big.Int`：按 u128 字节存取（bytesToU128 / u128ToBytes）
  - 复杂类型（slice、struct、pointer）：`model.TxnGetJson` / `model.TxnSetJson`
  - **StoreMapping[K]**：在 state 上生成字段（如 `Tracks`、`Proposals` 等），通过 `state.Tracks.Get(txn, key)` / `state.Proposals.Set(txn, key, val)` 等存取

## 用法

### 单结构体

```bash
go run ./tools/stategen \
  -pkg=sidechain \
  -file=./side-chain/dao.go \
  -struct=daoState \
  -namespace=dao \
  -prefix=state \
  -out=./side-chain/dao_state_gen.go
```

### 多结构体（YAML 配置）

```yaml
# spec.yaml
gen:
  - pkg: sidechain
    file: ./side-chain/dao.go
    struct: daoState
    namespace: dao
    prefix: state
    out: ./side-chain/dao_state_gen.go
```

```bash
go run ./tools/stategen -config=spec.yaml
```

### 可选：自定义 key 名

在 struct 的字段 tag 里用 `state:"key=xxx"` 指定存储 key，否则用 `json` tag 或字段名转 snake_case。

## 使用示例

```go
// 创建访问器（变量名小写）
state := newDaoStateStorageState(txn)

// 读
members := state.Members()
total := state.TotalIssuance()

// 写
state.SetMembers(newMembers)
state.SetTotalIssuance("1000")
```

## 与现有 dao 的配合

当前 `dao.go` 已手动拆出「配置块」与「计数器」分离存储；对新 struct 或希望统一风格时，可对目标 struct 跑 stategen，用生成的 `state.Xxx()` / `state.SetXxx(v)` 替代整块 load/save，进一步减少写放大。
