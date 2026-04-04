#!/usr/bin/env bash
# 从合约包 AST 生成 dao_gen.go（ExecCall / ExecQuery）。
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
exec go run ./tools/contractgen \
  -dir ./side-chain/contracts/dao \
  -out ./side-chain/contracts/dao/dao_gen.go \
  -mutation DaoMutation \
  -query DaoQuery \
  -skip-mutation Delete
