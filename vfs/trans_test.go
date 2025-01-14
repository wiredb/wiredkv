package vfs

import (
	"testing"
)

// 测试 Transformer 类的压缩、加密和解密功能
func TestTransformerWithComplexData(t *testing.T) {
	// 创建一个新的 Transformer
	transformer := NewTransformer()

	// 构造复杂数据结构，包括 uint 和字符串
	originalString := "example-data"

	// 启用压缩
	transformer.SetCompressor(&Snappy{})

	// 对数据进行编码（压缩 + 加密）
	encodedData, err := transformer.Encode([]byte(originalString))
	if err != nil {
		t.Fatalf("failed to encode data: %v", err)
	}

	// 解码数据
	decodedData, err := transformer.Decode(encodedData)
	if err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if originalString != string(decodedData) {
		t.Fatalf("failed to decode data: got %s, want %s", decodedData, originalString)
	}
}
