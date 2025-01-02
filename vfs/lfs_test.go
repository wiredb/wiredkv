package vfs

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"testing"
)

// TestSerializedIndex 测试 serializedIndex 函数
func TestSerializedIndex(t *testing.T) {
	// 创建一个测试的 INode
	inode := &INode{
		RegionID:  1234,
		Offset:    5678,
		Length:    100,
		ExpiredAt: 1617181723,
		CreatedAt: 1617181623,
	}

	// 计算预期的字节切片
	expectedLength := 48

	// 调用 serializeIndex
	result, err := serializedIndex(1001, inode)
	if err != nil {
		t.Fatalf("serialized index failed: %v", err)
	}

	// 检查返回的字节切片长度
	if len(result) != expectedLength {
		t.Errorf("expected result length %d, got %d", expectedLength, len(result))
	}

	// 验证内容
	buf := bytes.NewReader(result)

	// 逐个字段进行反序列化并检查
	var inum, regionID, offset uint64
	var length, checksum uint32
	var expiredAt, createdAt int64

	// 反序列化每个字段
	binary.Read(buf, binary.LittleEndian, &inum)
	binary.Read(buf, binary.LittleEndian, &regionID)
	binary.Read(buf, binary.LittleEndian, &offset)
	binary.Read(buf, binary.LittleEndian, &length)
	binary.Read(buf, binary.LittleEndian, &expiredAt)
	binary.Read(buf, binary.LittleEndian, &createdAt)
	binary.Read(buf, binary.LittleEndian, &checksum)

	// 验证字段是否一致
	if inum != 1001 {
		t.Errorf("expected inum %d, got %d", 1001, inum)
	}
	if regionID != inode.RegionID {
		t.Errorf("expected RegionID %d, got %d", inode.RegionID, regionID)
	}
	if offset != inode.Offset {
		t.Errorf("expected Offset %d, got %d", inode.Offset, offset)
	}
	if length != inode.Length {
		t.Errorf("expected Length %d, got %d", inode.Length, length)
	}
	if expiredAt != inode.ExpiredAt {
		t.Errorf("expected ExpiredAt %d, got %d", inode.ExpiredAt, expiredAt)
	}
	if createdAt != inode.CreatedAt {
		t.Errorf("expected CreatedAt %d, got %d", inode.CreatedAt, createdAt)
	}

	// 验证 CRC32 校验码
	expectedChecksum := crc32.ChecksumIEEE(result[:len(result)-4])
	if checksum != expectedChecksum {
		t.Errorf("expected checksum %d, got %d", expectedChecksum, checksum)
	}
}
