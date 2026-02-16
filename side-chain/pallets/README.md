# Pallets：侧链合约接口抽象

本目录定义**标准接口**（可选）；各 pallet 直接导出 `ApplyXxxCall` 函数，由调用方按需调用，写法简洁。

## 接口（可选）

- **`Pallet`**：可选实现，用于类型约束或统一抽象。
  - `ApplyCall(caller, payload []byte, height int64, txn *model.Txn) error`
- **`CallFunc`**：将函数适配为 `Pallet`。

## 调用方式

推荐直接调用各 pallet 导出的入口函数，例如：

- `dao.ApplyDaoCall(caller, payload, height, txn)`

无需通过 `dao.Pallet.ApplyCall(...)`。

## 新增一个 Pallet 的步骤

1. 在 `pallets/` 下新建目录，如 `pallets/foo/`。
2. 实现并导出入口函数：`func ApplyFooCall(caller, payload []byte, height int64, txn *model.Txn) error`。
3. 在 `tx_finalize.go` 等调用处：导入该 pallet 包，在对应 Tx 分支中调用 `foo.ApplyFooCall(caller, payload, height, txn)`。
