package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// inkRegistry 在内存中构建 ink metadata v6 的 types 数组：先注册与历史模板一致的 0–46，
// 再追加 DAO 合约用到的复合类型与 Result<T, LangError>，不依赖任何外部 JSON 模板文件。
type inkRegistry struct {
	types        []map[string]any
	resultQLang  map[int]int // inner type id -> Result<T, LangError> type id（含预置 14/16/33）
	optCallID    int
	propStatusID int
	propDepositID int
	proposalID   int
	voteID       int
	optU32ID     int
	optVoteID    int
	optProposalID int
	optPropStatID int
	vecTrackID   int
	vecProposalID int
	vecVoteID    int
}

func newInkRegistry() *inkRegistry {
	r := &inkRegistry{resultQLang: make(map[int]int)}
	b := &inkBuilder{}
	registerInkV6BaseTypes(b)
	r.types = b.types
	// 预置与旧模板一致的查询返回 Result 类型
	r.resultQLang[6] = 14  // bool
	r.resultQLang[3] = 16  // U256
	r.resultQLang[32] = 33 // Option<Track>
	r.registerDAOComposites()
	return r
}

type inkBuilder struct {
	types []map[string]any
}

func (b *inkBuilder) addID(id int, typ map[string]any) int {
	if id != len(b.types) {
		panic(fmt.Sprintf("inkBuilder: expected id %d, have %d entries", id, len(b.types)))
	}
	b.types = append(b.types, map[string]any{"id": id, "type": typ})
	return id
}

func mustJSONType(s string) map[string]any {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	return m
}

// registerInkV6BaseTypes 注册 ink v6 基础类型 0–46（与既有链上 metadata 环境字段引用一致；类型体由代码拼装，不读外部文件）。
func registerInkV6BaseTypes(b *inkBuilder) {
	_ = b.addID(0, mustJSONType(`{"def":{"primitive":"u8"}}`))
	_ = b.addID(1, mustJSONType(`{"def":{"array":{"len":20,"type":0}}}`))
	_ = b.addID(2, mustJSONType(`{"def":{"composite":{"fields":[{"type":1,"typeName":"[u8; 20]"}]}},"path":["primitive_types","H160"]}`))
	_ = b.addID(3, mustJSONType(`{"def":{"composite":{}},"path":["U256"]}`))
	_ = b.addID(4, mustJSONType(`{"def":{"tuple":[2,3]}}`))
	_ = b.addID(5, mustJSONType(`{"def":{"sequence":{"type":4}}}`))
	_ = b.addID(6, mustJSONType(`{"def":{"primitive":"bool"}}`))
	_ = b.addID(7, mustJSONType(`{"def":{"variant":{"variants":[{"index":0,"name":"None"},{"fields":[{"type":2}],"index":1,"name":"Some"}]}},"params":[{"name":"T","type":2}],"path":["Option"]}`))
	_ = b.addID(8, mustJSONType(`{"def":{"tuple":[]}}`))
	_ = b.addID(9, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[],"index":0,"name":"TokenNotFound"},{"fields":[],"index":1,"name":"MemberExisted"},{"fields":[],"index":2,"name":"MemberNotExisted"},{"fields":[],"index":3,"name":"MemberBalanceNotZero"},{"fields":[],"index":4,"name":"PublicJoinNotAllowed"},{"fields":[],"index":5,"name":"LowBalance"},{"fields":[],"index":6,"name":"InsufficientAllowance"},{"fields":[],"index":7,"name":"CallFailed"},{"fields":[],"index":8,"name":"InvalidDeposit"},{"fields":[],"index":9,"name":"TransferFailed"},{"fields":[],"index":10,"name":"MustCallByGov"},{"fields":[],"index":11,"name":"PropNotOngoing"},{"fields":[],"index":12,"name":"PropNotEnd"},{"fields":[],"index":13,"name":"InvalidProposal"},{"fields":[],"index":14,"name":"InvalidProposalStatus"},{"fields":[],"index":15,"name":"InvalidProposalCaller"},{"fields":[],"index":16,"name":"InvalidDepositTime"},{"fields":[],"index":17,"name":"InvalidVoteTime"},{"fields":[],"index":18,"name":"InvalidVoteStatus"},{"fields":[],"index":19,"name":"InvalidVoteUser"},{"fields":[],"index":20,"name":"ProposalInDecision"},{"fields":[],"index":21,"name":"VoteAlreadyUnlocked"},{"fields":[],"index":22,"name":"InvalidVoteUnlockTime"},{"fields":[],"index":23,"name":"ProposalNotConfirmed"},{"fields":[],"index":24,"name":"NoTrack"},{"fields":[],"index":25,"name":"MaxBalanceOverflow"},{"fields":[],"index":26,"name":"TransferDisable"},{"fields":[],"index":27,"name":"InvalidVote"},{"fields":[],"index":28,"name":"SetCodeFailed"},{"fields":[],"index":29,"name":"SpendNotFound"},{"fields":[],"index":30,"name":"SpendAlreadyExecuted"},{"fields":[],"index":31,"name":"SpendTransferError"}]}},"path":["Error"]}`))
	_ = b.addID(10, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":8}],"index":0,"name":"Ok"},{"fields":[{"type":9}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":8},{"name":"E","type":9}],"path":["Result"]}`))
	_ = b.addID(11, mustJSONType(`{"def":{"sequence":{"type":2}}}`))
	_ = b.addID(12, mustJSONType(`{"def":{"variant":{"variants":[{"index":1,"name":"CouldNotReadInput"}]}},"path":["ink_primitives","LangError"]}`))
	_ = b.addID(13, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":11}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":11},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(14, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":6}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":6},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(15, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":10}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":10},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(16, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":3}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":3},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(17, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":7}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":7},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(18, mustJSONType(`{"def":{"array":{"len":4,"type":0}},"path":["Selector"]}`))
	_ = b.addID(19, mustJSONType(`{"def":{"sequence":{"type":0}}}`))
	_ = b.addID(20, mustJSONType(`{"def":{"primitive":"u64"}}`))
	_ = b.addID(21, mustJSONType(`{"def":{"composite":{"fields":[{"name":"contract","type":7,"typeName":"Option<Address>"},{"name":"selector","type":18,"typeName":"Selector"},{"name":"input","type":19,"typeName":"Vec<u8>"},{"name":"amount","type":3,"typeName":"U256"},{"name":"ref_time_limit","type":20,"typeName":"u64"},{"name":"allow_reentry","type":6,"typeName":"bool"}]}},"path":["Call"]}`))
	_ = b.addID(22, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":19}],"index":0,"name":"Ok"},{"fields":[{"type":9}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":19},{"name":"E","type":9}],"path":["Result"]}`))
	_ = b.addID(23, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":22}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":22},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(24, mustJSONType(`{"def":{"primitive":"u16"}}`))
	_ = b.addID(25, mustJSONType(`{"def":{"variant":{"variants":[{"index":0,"name":"None"},{"fields":[{"type":24}],"index":1,"name":"Some"}]}},"params":[{"name":"T","type":24}],"path":["Option"]}`))
	_ = b.addID(26, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":25}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":25},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(27, mustJSONType(`{"def":{"primitive":"u32"},"path":["BlockNumber"]}`))
	_ = b.addID(28, mustJSONType(`{"def":{"primitive":"u32"}}`))
	_ = b.addID(29, mustJSONType(`{"def":{"primitive":"i64"}}`))
	_ = b.addID(30, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":28,"typeName":"u32"},{"type":28,"typeName":"u32"},{"type":27,"typeName":"BlockNumber"}],"index":0,"name":"LinearDecreasing"},{"fields":[{"type":28,"typeName":"u32"},{"type":28,"typeName":"u32"},{"type":28,"typeName":"u32"},{"type":27,"typeName":"BlockNumber"}],"index":1,"name":"SteppedDecreasing"},{"fields":[{"type":28,"typeName":"u32"},{"type":28,"typeName":"u32"},{"type":29,"typeName":"i64"},{"type":29,"typeName":"i64"}],"index":2,"name":"Reciprocal"}]}},"path":["Curve"]}`))
	_ = b.addID(31, mustJSONType(`{"def":{"composite":{"fields":[{"name":"name","type":19,"typeName":"Vec<u8>"},{"name":"prepare_period","type":27,"typeName":"BlockNumber"},{"name":"decision_deposit","type":3,"typeName":"U256"},{"name":"max_deciding","type":27,"typeName":"BlockNumber"},{"name":"confirm_period","type":27,"typeName":"BlockNumber"},{"name":"decision_period","type":27,"typeName":"BlockNumber"},{"name":"min_enactment_period","type":27,"typeName":"BlockNumber"},{"name":"max_balance","type":3,"typeName":"U256"},{"name":"min_approval","type":30,"typeName":"Curve"},{"name":"min_support","type":30,"typeName":"Curve"}]}},"path":["Track"]}`))
	_ = b.addID(32, mustJSONType(`{"def":{"variant":{"variants":[{"index":0,"name":"None"},{"fields":[{"type":31}],"index":1,"name":"Some"}]}},"params":[{"name":"T","type":31}],"path":["Option"]}`))
	_ = b.addID(33, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":32}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":32},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(34, mustJSONType(`{"def":{"tuple":[24,31]}}`))
	_ = b.addID(35, mustJSONType(`{"def":{"sequence":{"type":34}}}`))
	_ = b.addID(36, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":35}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":35},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(37, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":24}],"index":0,"name":"Ok"},{"fields":[{"type":9}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":24},{"name":"E","type":9}],"path":["Result"]}`))
	_ = b.addID(38, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":37}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":37},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(39, mustJSONType(`{"def":{"variant":{"variants":[{"index":0,"name":"None"},{"fields":[{"type":18}],"index":1,"name":"Some"}]}},"params":[{"name":"T","type":18}],"path":["Option"]}`))
	_ = b.addID(40, mustJSONType(`{"def":{"composite":{"fields":[{"name":"name","type":19,"typeName":"Vec<u8>"},{"name":"symbol","type":19,"typeName":"Vec<u8>"},{"name":"decimals","type":0,"typeName":"u8"}]}},"path":["TokenInfo"]}`))
	_ = b.addID(41, mustJSONType(`{"def":{"variant":{"variants":[{"index":0,"name":"None"},{"fields":[{"type":40}],"index":1,"name":"Some"}]}},"params":[{"name":"T","type":40}],"path":["Option"]}`))
	_ = b.addID(42, mustJSONType(`{"def":{"variant":{"variants":[{"fields":[{"type":41}],"index":0,"name":"Ok"},{"fields":[{"type":12}],"index":1,"name":"Err"}]}},"params":[{"name":"T","type":41},{"name":"E","type":12}],"path":["Result"]}`))
	_ = b.addID(43, mustJSONType(`{"def":{"composite":{}},"path":["H256"]}`))
	_ = b.addID(44, mustJSONType(`{"def":{"primitive":"u128"}}`))
	_ = b.addID(45, mustJSONType(`{"def":{"array":{"len":32,"type":0}}}`))
	_ = b.addID(46, mustJSONType(`{"def":{"composite":{"fields":[{"type":45,"typeName":"[u8; 32]"}]}},"path":["ink_primitives","types","Hash"]}`))
}

func (r *inkRegistry) appendType(typ map[string]any) int {
	id := len(r.types)
	r.types = append(r.types, map[string]any{"id": id, "type": typ})
	return id
}

func (r *inkRegistry) resultVariantLang(okType, errType int) map[string]any {
	return map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"fields": []any{map[string]any{"type": okType}}, "index": 0, "name": "Ok"},
					map[string]any{"fields": []any{map[string]any{"type": errType}}, "index": 1, "name": "Err"},
				},
			},
		},
		"params": []any{
			map[string]any{"name": "T", "type": okType},
			map[string]any{"name": "E", "type": errType},
		},
		"path": []any{"Result"},
	}
}

// ensureResultQueryLang 返回 Result<T, ink::LangError> 的 type id（T 为 SCALE 返回值类型）。
func (r *inkRegistry) ensureResultQueryLang(inner int) int {
	if id, ok := r.resultQLang[inner]; ok {
		return id
	}
	id := r.appendType(r.resultVariantLang(inner, 12))
	r.resultQLang[inner] = id
	return id
}

// registerDAOComposites 追加 side-chain/contracts/dao/types.go 中提案/投票等复合类型（与 SCALE 字段顺序一致）。
func (r *inkRegistry) registerDAOComposites() {
	// Option<Call> — Call 为 id 21
	r.optCallID = r.appendType(map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"index": 0, "name": "None"},
					map[string]any{"fields": []any{map[string]any{"type": 21}}, "index": 1, "name": "Some"},
				},
			},
		},
		"params": []any{map[string]any{"name": "T", "type": 21}},
		"path":   []any{"Option"},
	})
	// ProposalStatus: ProposalState 编码为 Vec<u8>，Block 为 i64
	r.propStatusID = r.appendType(map[string]any{
		"def": map[string]any{
			"composite": map[string]any{
				"fields": []any{
					map[string]any{"name": "state", "type": 19, "typeName": "Vec<u8>"},
					map[string]any{"name": "block", "type": 29, "typeName": "i64"},
				},
			},
		},
		"path": []any{"ProposalStatus"},
	})
	r.propDepositID = r.appendType(map[string]any{
		"def": map[string]any{
			"composite": map[string]any{
				"fields": []any{
					map[string]any{"name": "depositor", "type": 2, "typeName": "Address"},
					map[string]any{"name": "amount", "type": 3, "typeName": "U256"},
					map[string]any{"name": "block", "type": 29, "typeName": "i64"},
				},
			},
		},
		"path": []any{"ProposalDeposit"},
	})
	r.proposalID = r.appendType(map[string]any{
		"def": map[string]any{
			"composite": map[string]any{
				"fields": []any{
					map[string]any{"name": "id", "type": 28, "typeName": "u32"},
					map[string]any{"name": "call", "type": r.optCallID, "typeName": "Option<Call>"},
					map[string]any{"name": "track_id", "type": 28, "typeName": "u32"},
					map[string]any{"name": "caller", "type": 2, "typeName": "Address"},
					map[string]any{"name": "status", "type": r.propStatusID, "typeName": "ProposalStatus"},
					map[string]any{"name": "submit_block", "type": 29, "typeName": "i64"},
					map[string]any{"name": "decision_block", "type": 29, "typeName": "i64"},
					map[string]any{"name": "deposit", "type": r.propDepositID, "typeName": "ProposalDeposit"},
				},
			},
		},
		"path": []any{"Proposal"},
	})
	r.voteID = r.appendType(map[string]any{
		"def": map[string]any{
			"composite": map[string]any{
				"fields": []any{
					map[string]any{"name": "id", "type": 20, "typeName": "u64"},
					map[string]any{"name": "proposal_id", "type": 28, "typeName": "u32"},
					map[string]any{"name": "caller", "type": 2, "typeName": "Address"},
					map[string]any{"name": "pledge", "type": 3, "typeName": "U256"},
					map[string]any{"name": "opinion_yes", "type": 6, "typeName": "bool"},
					map[string]any{"name": "vote_weight", "type": 3, "typeName": "U256"},
					map[string]any{"name": "unlock_block", "type": 29, "typeName": "i64"},
					map[string]any{"name": "vote_block", "type": 29, "typeName": "i64"},
					map[string]any{"name": "deleted", "type": 6, "typeName": "bool"},
				},
			},
		},
		"path": []any{"Vote"},
	})
	// Option<u32> — DefaultTrack *uint32
	r.optU32ID = r.appendType(map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"index": 0, "name": "None"},
					map[string]any{"fields": []any{map[string]any{"type": 28}}, "index": 1, "name": "Some"},
				},
			},
		},
		"params": []any{map[string]any{"name": "T", "type": 28}},
		"path":   []any{"Option"},
	})
	r.optVoteID = r.appendType(map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"index": 0, "name": "None"},
					map[string]any{"fields": []any{map[string]any{"type": r.voteID}}, "index": 1, "name": "Some"},
				},
			},
		},
		"params": []any{map[string]any{"name": "T", "type": r.voteID}},
		"path":   []any{"Option"},
	})
	r.optProposalID = r.appendType(map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"index": 0, "name": "None"},
					map[string]any{"fields": []any{map[string]any{"type": r.proposalID}}, "index": 1, "name": "Some"},
				},
			},
		},
		"params": []any{map[string]any{"name": "T", "type": r.proposalID}},
		"path":   []any{"Option"},
	})
	r.optPropStatID = r.appendType(map[string]any{
		"def": map[string]any{
			"variant": map[string]any{
				"variants": []any{
					map[string]any{"index": 0, "name": "None"},
					map[string]any{"fields": []any{map[string]any{"type": r.propStatusID}}, "index": 1, "name": "Some"},
				},
			},
		},
		"params": []any{map[string]any{"name": "T", "type": r.propStatusID}},
		"path":   []any{"Option"},
	})
	r.vecTrackID = r.appendType(map[string]any{"def": map[string]any{"sequence": map[string]any{"type": 31}}})
	r.vecProposalID = r.appendType(map[string]any{"def": map[string]any{"sequence": map[string]any{"type": r.proposalID}}})
	r.vecVoteID = r.appendType(map[string]any{"def": map[string]any{"sequence": map[string]any{"type": r.voteID}}})
}

func (r *inkRegistry) ensureInnerQueryReturn(goType string) (int, error) {
	goType = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(goType), "*"))
	// 查询返回 []byte 视为 U256（余额等），与参数里的 Address 区分
	if goType == "[]byte" {
		return 3, nil
	}
	kind, inner := parseGoTypeRoot(goType)
	switch kind {
	case "slice":
		inner = strings.TrimSpace(inner)
		switch inner {
		case "Member":
			return 5, nil // Vec<Member>
		case "TrackData":
			return r.vecTrackID, nil
		case "Proposal":
			return r.vecProposalID, nil
		case "Vote":
			return r.vecVoteID, nil
		}
	case "option":
		inner = strings.TrimSpace(inner)
		switch inner {
		case "TrackData":
			return 32, nil // Option<Track>
		case "Vote":
			return r.optVoteID, nil
		case "Proposal":
			return r.optProposalID, nil
		case "ProposalStatus":
			return r.optPropStatID, nil
		}
	case "ptr":
		inner = strings.TrimSpace(inner)
		if inner == "uint32" {
			return r.optU32ID, nil
		}
	case "named":
		switch goType {
		case "bool":
			return 6, nil
		case "uint32":
			return 28, nil
		case "uint64":
			return 20, nil
		}
	}
	return 0, fmt.Errorf("unknown query return type %q", goType)
}

func parseGoTypeRoot(s string) (kind, inner string) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "[]") {
		return "slice", strings.TrimPrefix(s, "[]")
	}
	if strings.HasPrefix(s, "util.Option[") && strings.HasSuffix(s, "]") {
		return "option", s[len("util.Option[") : len(s)-1]
	}
	if strings.HasPrefix(s, "*") {
		return "ptr", strings.TrimPrefix(s, "*")
	}
	return "named", s
}

func (r *inkRegistry) ensureResultQuery(goType string) (map[string]any, error) {
	inner, err := r.ensureInnerQueryReturn(goType)
	if err != nil {
		return nil, err
	}
	rid := r.ensureResultQueryLang(inner)
	return map[string]any{"displayName": []any{"Result"}, "type": rid}, nil
}

func (r *inkRegistry) goTypeToABIRefParam(goType string) (abiTypeRef, error) {
	goType = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(goType), "*"))
	if goType == "[]byte" {
		return abiTypeRef{Display: []string{"Address"}, TypeID: 2}, nil
	}
	return r.goTypeToABIRefNamed(goType)
}

func (r *inkRegistry) goTypeToABIRefNamed(goType string) (abiTypeRef, error) {
	switch goType {
	case "bool":
		return abiTypeRef{Display: []string{"bool"}, TypeID: 6}, nil
	case "uint32":
		return abiTypeRef{Display: []string{"u32"}, TypeID: 28}, nil
	case "uint64":
		return abiTypeRef{Display: []string{"u64"}, TypeID: 20}, nil
	case "[]Member":
		return abiTypeRef{Display: []string{"Vec"}, TypeID: 5}, nil
	case "CallContent":
		return abiTypeRef{Display: []string{"Call"}, TypeID: 21}, nil
	case "TrackData":
		return abiTypeRef{Display: []string{"Track"}, TypeID: 31}, nil
	}
	return abiTypeRef{}, fmt.Errorf("unknown Go type %q for ABI param (extend inkRegistry.goTypeToABIRefNamed)", goType)
}

func inkEnvironment() map[string]any {
	return map[string]any{
		"accountId": map[string]any{
			"displayName": []any{"AccountId"},
			"type":          2,
		},
		"balance": map[string]any{
			"displayName": []any{"Balance"},
			"type":          44,
		},
		"blockNumber": map[string]any{
			"displayName": []any{"BlockNumber"},
			"type":          28,
		},
		"hash": map[string]any{
			"displayName": []any{"Hash"},
			"type":          46,
		},
		"nativeToEthRatio": 100000000,
		"staticBufferSize": 16384,
		"timestamp": map[string]any{
			"displayName": []any{"Timestamp"},
			"type":          20,
		},
	}
}

func inkConstructor() []any {
	return []any{
		map[string]any{
			"args": []any{
				map[string]any{
					"docs":  []any{},
					"label": "users",
					"type": map[string]any{
						"displayName": []any{"Vec"},
						"type":        5,
					},
				},
				map[string]any{
					"docs":  []any{},
					"label": "public_join",
					"type": map[string]any{
						"displayName": []any{"bool"},
						"type":        6,
					},
				},
				map[string]any{
					"docs":  []any{},
					"label": "sudo_account",
					"type": map[string]any{
						"displayName": []any{"Option"},
						"type":        7,
					},
				},
			},
			"default": false,
			"docs":    []any{},
			"label":   "new_with_default_track",
			"payable": false,
			"returnType": map[string]any{
				"displayName": []any{"ink", "MessageResult"},
				"type":        10,
			},
			"selector": "0x00000000",
		},
	}
}

func buildInkRoot(contractName string, mutMethods, qMethods []*methodSig) (map[string]any, error) {
	reg := newInkRegistry()
	messages, err := buildInkMessages(reg, mutMethods, qMethods)
	if err != nil {
		return nil, err
	}
	spec := map[string]any{
		"constructors": inkConstructor(),
		"docs":         []any{},
		"environment":  inkEnvironment(),
		"events":       []any{},
		"lang_error": map[string]any{
			"displayName": []any{"ink", "LangError"},
			"type":        12,
		},
		"messages": messages,
	}
	root := map[string]any{
		"contract": map[string]any{
			"name":    contractNameGoSafe(contractName),
			"version": "0.1.0",
		},
		"image": nil,
		"source": map[string]any{
			"compiler": "contractgen",
			"hash":     "0x0000000000000000000000000000000000000000000000000000000000000000",
			"language": "go",
		},
		"spec":    spec,
		"storage": nil,
		"types":   reg.types,
		"version": 6,
	}
	return root, nil
}

func buildInkMessages(reg *inkRegistry, mutMethods, qMethods []*methodSig) ([]any, error) {
	type tagged struct {
		ms       *methodSig
		mutates  bool
		abiLabel string
	}
	var all []tagged
	for _, m := range mutMethods {
		all = append(all, tagged{ms: m, mutates: true, abiLabel: m.caseName})
	}
	for _, m := range qMethods {
		all = append(all, tagged{ms: m, mutates: false, abiLabel: m.caseName})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].abiLabel < all[j].abiLabel })
	out := make([]any, 0, len(all))
	for _, t := range all {
		msg, err := buildInkMessage(reg, t.ms, t.mutates, t.abiLabel)
		if err != nil {
			return nil, fmt.Errorf("abi %s: %w", t.abiLabel, err)
		}
		out = append(out, msg)
	}
	return out, nil
}

func buildInkMessage(reg *inkRegistry, ms *methodSig, mutates bool, abiLabel string) (map[string]any, error) {
	args, err := abiArgs(reg, ms)
	if err != nil {
		return nil, err
	}
	var ret map[string]any
	if mutates {
		ret = map[string]any{"displayName": []any{"Result"}, "type": 15}
	} else {
		ret, err = reg.ensureResultQuery(ms.queryResult0)
		if err != nil {
			return nil, err
		}
	}
	sel := pickSelectorInk(abiLabel, mutates)
	return abiMessageJSON(args, abiLabel, mutates, ret, sel), nil
}

func abiArgs(reg *inkRegistry, ms *methodSig) ([]any, error) {
	if ms.isInit {
		out := make([]any, 0, 3)
		for i := 0; i < 3 && i < len(ms.params); i++ {
			a, err := abiOneArg(reg, ms.params[i], i)
			if err != nil {
				return nil, err
			}
			out = append(out, a)
		}
		return out, nil
	}
	out := make([]any, 0, len(ms.params))
	for i := range ms.params {
		a, err := abiOneArg(reg, ms.params[i], i)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func abiOneArg(reg *inkRegistry, p param, idx int) (map[string]any, error) {
	ref, err := reg.goTypeToABIRefParam(p.underlying)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"docs": []any{},
		"label": paramABIName(p.name, idx),
		"type": map[string]any{
			"displayName": ref.Display,
			"type":        ref.TypeID,
		},
	}, nil
}

