# daogen：DAO 合约 Graph 与 API 代码生成

根据 `side-chain/pallets/dao/dao.go` 中的 `Op*` 常量与合约定义，生成：

1. **GraphQL 查询与调用**：`graph/dao.graphqls`
   - `Query.daoState`：查询 DAO 状态（JSON 字符串）
   - `Mutation.daoCall(caller, payload)`：执行 DAO 调用

2. **合约函数 API**：`pkg/api/dao/payload_gen.go`
   - `Payload` 与 `MemberInput`、`TrackInput`、`CallInput` 等类型，与合约 `DaoCallPayload` 对齐
   - `MustPayload(p)`：将 `Payload` 序列化为 JSON
   - 每个 op 一个便捷函数：`PayloadDaoInit(p)`、`PayloadDaoTransfer(p)`、…，用于构建对应 op 的 JSON 字符串

## 用法

```bash
go run ./tools/daogen \
  -dao=./side-chain/pallets/dao/dao.go \
  -graph=./graph/dao.graphqls \
  -api=./pkg/api/dao/payload_gen.go
```

- `-dao`：必填，dao.go 路径（用于解析 Op* 常量）
- `-graph`：输出 GraphQL schema 路径
- `-api`：输出 Go API 路径（可选）

生成后需执行 `go run github.com/99designs/gqlgen generate` 以更新 graph 的 resolver 与 generated 代码。

## 调用示例（Go）

```go
import "github.com/wetee-dao/tee-dsecret/pkg/api/dao"

// 构建 dao_transfer 的 payload
p := dao.Payload{ To: toBytes, Value: "1000000" }
payload := dao.PayloadDaoTransfer(p)
// 再通过 GraphQL mutation daoCall(caller, payload) 或侧链 SubmitTx 提交
```

## 调用示例（GraphQL）

```graphql
mutation {
  daoCall(caller: "0x...", payload: "{\"op\":\"dao_transfer\",\"to\":\"...\",\"value\":\"1000000\"}")
}
```

前端可依赖 `pkg/api/dao` 的类型与 `PayloadXxx` 生成 TypeScript 类型或直接拼 JSON。
