#!/usr/bin/env bash
# 从合约包 AST 生成 dao_gen.go（ExecCall / ExecQuery）及 side-chain/abis/dao.json（ink metadata 在内存中由 AST + 内建类型注册表构建，无默认 JSON 模板）。
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
exec go run ./tools/contractgen \
  -dir ./side-chain/contracts/dao \
  -out ./side-chain/contracts/dao/dao_gen.go \
  -mutation DaoMutation \
  -query DaoQuery \
  -struct DAO \
  -skip-mutation Delete
