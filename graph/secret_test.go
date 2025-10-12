package graph

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestDecryptSecret(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	// 模拟加密数据（使用 jsencrypt 加密后的 Base64 字符串）
	encryptedBase64 := "rSeF3KcQyPyw0DfxF4k4jPnwu72/BvXyQkEofutS5MXSvg7gpVBtE1SWK/0H32UUqFneyR6rNyuUhXbKtPzOPaMCkZextDKc+2CMk1ywj99QfZ6HYwnY2U8K+ZhDtIkuCppJIKtKV45aLa924vf7b73xnejteuZOeT7IbDrYdBGhQjhkZBODKyuq89m1EMChk4ZGi0MEb9V6i1AHKm+G3OtK/4/YIOPBhSIlQhQ7ASIaFsN/H+ugZhjJSV6NOcFPm9aLPt+y4WceVZmzBu4BFr4Wc/YssccJCLMz61TXrCmqJEhtDDtqvXKrpuApKtx/tsAAP1vfQJX+DEwaANcnmA=="

	// 解密
	plaintext, err := rsaDecryptWithKey(privateKey, encryptedBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	// 验证解密结果（这里假设加密前的数据是一个简单的字符串）
	expected := "这是一个测试字符串"
	if string(plaintext) != expected {
		t.Errorf("解密结果错误: 期望 %s, 实际 %s", expected, string(plaintext))
	}
}
