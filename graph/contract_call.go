package graph

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/dao"
)

// DecodeCaller 支持 32 字节 hex（可带 0x 前缀）或 SS58，返回公钥字节。
func DecodeCaller(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	if b, err := hex.DecodeString(s); err == nil && len(b) == 32 {
		return b, nil
	}
	pub, err := model.PubKeyFromSS58(s)
	if err != nil {
		return nil, fmt.Errorf("caller 需为 32 字节 hex 或 SS58: %w", err)
	}
	return pub.Byte(), nil
}

// SubmitContractCall 根据合约名提交交易。dao 的 payload 为 base64 编码的 model.DaoCall protobuf。
func SubmitContractCall(caller []byte, contract string, payload string) error {
	var tx *model.Tx
	switch contract {
	case "dao":
		raw, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return fmt.Errorf("dao payload 需为 base64 编码的 protobuf: %w", err)
		}
		tx = &model.Tx{
			Caller:  caller,
			Payload: &model.Tx_DaoCall{DaoCall: raw},
		}
	default:
		return fmt.Errorf("不支持的合约: %s", contract)
	}
	_, err := sidechain.SubmitTx(tx)
	return err
}

// ContractQuery 按合约与方法单独查询，返回 JSON。args 为可选 JSON 字符串。
func ContractQuery(contract string, method string, args *string) (string, error) {
	switch contract {
	case "dao":
		return daoQuery(method, args)
	default:
		return "", fmt.Errorf("不支持的合约: %s", contract)
	}
}

// parseArgBytes 从 argMap 中取出 key 对应的地址（hex 或 SS58），返回 32 字节。
func parseArgBytes(argMap map[string]interface{}, key string) ([]byte, error) {
	v, ok := argMap[key]
	if !ok {
		return nil, fmt.Errorf("缺少参数 %s", key)
	}
	s, _ := v.(string)
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		pub, err := model.PubKeyFromSS58(s)
		if err != nil {
			return nil, fmt.Errorf("%s 需为 32 字节 hex 或 SS58", key)
		}
		return pub.Byte(), nil
	}
	return b, nil
}

// parseArgU32 从 argMap 中取出 key 对应的 uint32（数字或数字字符串）。
func parseArgU32(argMap map[string]interface{}, key string) (uint32, error) {
	v, ok := argMap[key]
	if !ok {
		return 0, fmt.Errorf("缺少参数 %s", key)
	}
	switch x := v.(type) {
	case float64:
		return uint32(x), nil
	case string:
		u, err := strconv.ParseUint(x, 10, 32)
		return uint32(u), err
	default:
		return 0, fmt.Errorf("%s 需为数字", key)
	}
}

// parseArgU64 从 argMap 中取出 key 对应的 uint64（数字或数字字符串）。
func parseArgU64(argMap map[string]interface{}, key string) (uint64, error) {
	v, ok := argMap[key]
	if !ok {
		return 0, fmt.Errorf("缺少参数 %s", key)
	}
	switch x := v.(type) {
	case float64:
		return uint64(x), nil
	case string:
		return strconv.ParseUint(x, 10, 64)
	default:
		return 0, fmt.Errorf("%s 需为数字", key)
	}
}

// daoQuery 分发 dao 的只读方法，args 为 JSON 如 {"owner":"0x..."}。
func daoQuery(method string, args *string) (string, error) {
	argMap := make(map[string]interface{})
	if args != nil && *args != "" {
		if err := json.Unmarshal([]byte(*args), &argMap); err != nil {
			return "", fmt.Errorf("args 需为合法 JSON: %w", err)
		}
	}

	var result interface{}
	switch method {
	case "member_list":
		result = dao.MemberList()
	case "get_public_join":
		result = dao.GetPublicJoin()
	case "name":
		result = dao.Name()
	case "symbol":
		result = dao.Symbol()
	case "decimals":
		result = dao.Decimals()
	case "total_supply":
		t := dao.TotalSupply()
		if t == nil {
			result = "0"
		} else {
			result = t.String()
		}
	case "balance_of":
		owner, err := parseArgBytes(argMap, "owner")
		if err != nil {
			return "", err
		}
		b := dao.BalanceOf(owner)
		result = b.String()
	case "allowance":
		owner, err := parseArgBytes(argMap, "owner")
		if err != nil {
			return "", err
		}
		spender, err := parseArgBytes(argMap, "spender")
		if err != nil {
			return "", err
		}
		b := dao.Allowance(owner, spender)
		result = b.String()
	case "lock_balance_of":
		owner, err := parseArgBytes(argMap, "owner")
		if err != nil {
			return "", err
		}
		b := dao.LockBalanceOf(owner)
		result = b.String()
	case "default_track":
		result = dao.DefaultTrack()
	case "track_list":
		page, _ := parseArgU32(argMap, "page")
		if page == 0 {
			page = 1
		}
		size, _ := parseArgU32(argMap, "size")
		if size == 0 {
			size = 10
		}
		result = dao.TrackList(page, size)
	case "track":
		id, err := parseArgU32(argMap, "id")
		if err != nil {
			return "", err
		}
		result = dao.Track(id)
	case "proposals":
		page, _ := parseArgU32(argMap, "page")
		if page == 0 {
			page = 1
		}
		size, _ := parseArgU32(argMap, "size")
		if size == 0 {
			size = 10
		}
		result = dao.Proposals(page, size)
	case "proposal":
		id, err := parseArgU32(argMap, "id")
		if err != nil {
			return "", err
		}
		result = dao.Proposal(id)
	case "vote_list":
		proposalId, err := parseArgU32(argMap, "proposal_id")
		if err != nil {
			return "", err
		}
		result = dao.VoteList(proposalId)
	case "vote":
		voteId, err := parseArgU64(argMap, "vote_id")
		if err != nil {
			return "", err
		}
		result = dao.Vote(voteId)
	case "proposal_status":
		proposalId, err := parseArgU32(argMap, "proposal_id")
		if err != nil {
			return "", err
		}
		result = dao.ProposalStatus(proposalId)
	default:
		return "", fmt.Errorf("未知方法: %s", method)
	}

	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
