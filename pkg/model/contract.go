package model

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"golang.org/x/crypto/blake2b"
)

// -----------------------------------------------------------------------------
// 对齐的合约基础类型（SCALE 固定宽字节 / U256）。
// -----------------------------------------------------------------------------

// Address 为 20 字节地址（EVM / 账户兼容），SCALE 编码为 20 字节原始数据。
type Address [20]byte

// AddressZero 全零地址（同 Rust Address::zero()）。
var AddressZero Address

// Bytes 返回底层字节数组副本的视图（值拷贝为数组）。
func (a Address) Bytes() [20]byte { return a }

// Slice 返回长度 20 的切片（指向 a 的副本时需先取地址再切片，此处为值接收者下的安全副本）。
func (a Address) Slice() []byte { return a[:] }

// AccountID 为 32 字节账户标识（如 Substrate AccountId）。
type AccountID [32]byte

// AccountIDZero 全零账户 ID。
var AccountIDZero AccountID

func (a AccountID) Bytes() [32]byte { return a }
func (a AccountID) Slice() []byte   { return a[:] }

// H256 为 32 字节哈希（如 Keccak-256），SCALE 编码为 32 字节。
type H256 [32]byte

// H256Zero 全零哈希。
var H256Zero H256

func (h H256) Bytes() [32]byte { return h }
func (h H256) Slice() []byte   { return h[:] }

// U256 为 256 位无符号整数，内部为大端 32 字节，与 wrevive U256 / EVM 一致。
type U256 [32]byte

// U256Zero、U256One、U256Max 常量语义同 Rust。
var (
	U256Zero = U256{}
	U256Max  = U256{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	U256One  = func() U256 { var b [32]byte; b[31] = 1; return U256(b) }()
)

// AsBytes 返回大端 32 字节（与内部表示相同）。
func (u U256) AsBytes() [32]byte { return u }

// FromBEBytes 从大端 32 字节构造 U256。
func U256FromBEBytes(b [32]byte) U256 { return U256(b) }

// FromU64 将 v 置于最低有效 8 字节（大端）。
func U256FromU64(v uint64) U256 {
	var b [32]byte
	binary.BigEndian.PutUint64(b[24:32], v)
	return U256(b)
}

// ToU64 取最低有效 8 字节大端为 u64（高位截断）。
func (u U256) ToU64() uint64 {
	return binary.BigEndian.Uint64(u[24:32])
}

// ToBEBytes 与 AsBytes 相同。
func (u U256) ToBEBytes() [32]byte { return u }

func u256Mask() *big.Int {
	return new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
}

func u256ToBig(u U256) *big.Int {
	return new(big.Int).SetBytes(u[:])
}

func u256FromBig(i *big.Int) U256 {
	m := u256Mask()
	x := new(big.Int).And(i, m)
	b := x.Bytes()
	var out [32]byte
	if len(b) == 0 {
		return U256(out)
	}
	if len(b) > 32 {
		copy(out[:], b[len(b)-32:])
		return U256(out)
	}
	copy(out[32-len(b):], b)
	return U256(out)
}

// ShlBits 左移 n 位，结果模 2^256；n≥256 为零。
func (u U256) ShlBits(n uint32) U256 {
	if n >= 256 {
		return U256Zero
	}
	if n == 0 {
		return u
	}
	x := u256ToBig(u)
	x.Lsh(x, uint(n))
	return u256FromBig(x)
}

// Bitor 按位或（模 2^256）。
func (u U256) Bitor(v U256) U256 {
	return u256FromBig(new(big.Int).Or(u256ToBig(u), u256ToBig(v)))
}

// WrappingAdd 模 2^256 加法。
func (u U256) WrappingAdd(v U256) U256 {
	x := u256ToBig(u)
	x.Add(x, u256ToBig(v))
	return u256FromBig(x)
}

// WrappingSub 模 2^256 减法。
func (u U256) WrappingSub(v U256) U256 {
	x := u256ToBig(u)
	x.Sub(x, u256ToBig(v))
	return u256FromBig(x)
}

// WrappingMul 模 2^256 乘法（低 256 位）。
func (u U256) WrappingMul(v U256) U256 {
	x := u256ToBig(u)
	x.Mul(x, u256ToBig(v))
	return u256FromBig(x)
}

// CheckedDiv 整除；除数为零时 ok 为 false。
func (u U256) CheckedDiv(v U256) (q U256, ok bool) {
	if v == U256Zero {
		return U256Zero, false
	}
	x := u256ToBig(u)
	y := u256ToBig(v)
	x.Div(x, y)
	return u256FromBig(x), true
}

// BlockNumber 区块高度（与 Rust BlockNumber = u32 一致）。
type BlockNumber uint32

// Bytes 为存储/消息中的变长字节，SCALE 为 compact 长度前缀 + 负载；语义同 types.rs 的 Vec<u8>。
type Bytes []byte

type ContractCall struct {
	Contract string
	Method   string
	Args     [][]byte
}

// ContractApi 侧链合约运行时：区块高度、读写事务、已解析的调用方（如由 TxBox 回填）。
type ContractApi interface {
	GetHeight() int64
	GetTxn() *Txn
	GetCaller() []byte
}

func RequireArgLen[T any](args []T, n int, method string) error {
	if len(args) < n {
		return fmt.Errorf("dao: %s expects at least %d arg(s), got %d", method, n, len(args))
	}
	return nil
}

// hexToBytes 解析十六进制（可带 0x），不允许 SS58；空串返回 nil 字节切片。
func hexToBytes(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}

// DecodeScaleArg 将单个参数视为 hex( SCALE(T) )，解码为 T（泛型入口，供合约 dispatch / 查询共用）。
func DecodeScaleArg[T any](s string) (T, error) {
	var z T
	raw, err := hexToBytes(s)
	if err != nil {
		return z, fmt.Errorf("hex: %w", err)
	}
	if len(raw) == 0 {
		return z, fmt.Errorf("dao: empty arg after hex decode")
	}
	return DecodeScaleArgBytes[T](raw)
}

func DecodeScaleArgBytes[T any](s []byte) (T, error) {
	var z T
	if err := codec.Decode(s, &z); err != nil {
		return z, fmt.Errorf("scale: %w", err)
	}
	return z, nil
}

// MethodToSelector 将 method 字符串转换为 4 字节 selector
// 支持两种格式：
// 1. 十六进制字符串（如 "0x1d623327" 或 "1d623327"）→ 直接解析
// 2. 普通字符串（如 "add_track"）→ 计算 blake2b 哈希取前 4 字节
// 注意：普通字符串默认使用 mutates=true，因为无法从字符串判断是 mutation 还是 query
// 如果需要区分，请使用 MethodToSelectorWithMutates 或直接传入十六进制 selector
func MethodToSelector(method string) [4]byte {
	return MethodToSelectorWithMutates(method, true)
}

// MethodToSelectorWithMutates 将 method 字符串转换为 4 字节 selector，可指定 mutates 参数
// 支持两种格式：
// 1. 十六进制字符串（如 "0x1d623327" 或 "1d623327"）→ 直接解析，忽略 mutates 参数
// 2. 普通字符串（如 "add_track"）→ 计算 blake2b 哈希取前 4 字节，使用指定的 mutates 值
func MethodToSelectorWithMutates(method string, mutates bool) [4]byte {
	// 尝试解析十六进制格式
	s := strings.TrimPrefix(method, "0x")
	if len(s) == 8 {
		if b, err := hex.DecodeString(s); err == nil && len(b) == 4 {
			return [4]byte(b)
		}
	}
	// 计算 blake2b 哈希
	var buf strings.Builder
	fmt.Fprintf(&buf, "wetee/contractgen/v1:%s:mutates=%v", method, mutates)
	h := blake2b.Sum256([]byte(buf.String()))
	var sel [4]byte
	copy(sel[:], h[:4])
	return sel
}
