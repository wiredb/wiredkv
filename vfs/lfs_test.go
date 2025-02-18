package vfs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/auula/wiredkv/conf"
	"github.com/auula/wiredkv/types"
)

// TestSerializedIndex 测试 serializedIndex 函数
func TestSerializedIndex(t *testing.T) {
	// 创建一个测试的 INode
	inode := &INode{
		RegionID:  1234,
		Position:  5678,
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

	// 验证内容字段进行反序列化并检查
	inum, node, err := deserializedIndex(result)

	if err != nil {
		t.Errorf("failed to deserialized: %v", err)
	}

	// 验证字段是否一致
	if inum != 1001 {
		t.Errorf("expected inum %d, got %d", 1001, inum)
	}
	if node.RegionID != inode.RegionID {
		t.Errorf("expected RegionID %d, got %d", inode.RegionID, node.RegionID)
	}
	if node.Position != inode.Position {
		t.Errorf("expected Offset %d, got %d", inode.Position, node.RegionID)
	}
	if node.Length != inode.Length {
		t.Errorf("expected Length %d, got %d", inode.Length, node.Length)
	}
	if node.ExpiredAt != inode.ExpiredAt {
		t.Errorf("expected ExpiredAt %d, got %d", inode.ExpiredAt, node.ExpiredAt)
	}
	if node.CreatedAt != inode.CreatedAt {
		t.Errorf("expected CreatedAt %d, got %d", inode.CreatedAt, node.CreatedAt)
	}

}

// 测试 readSegment 函数
func TestReadSegment(t *testing.T) {
	// 构造测试数据
	seg := &Segment{
		Tombstone: 0,
		Type:      1,
		ExpiredAt: 123456789,
		CreatedAt: 987654321,
		KeySize:   3,
		ValueSize: 5,
		Key:       []byte("key"),
		Value:     []byte("value"),
	}

	// 将 Segment 数据转化为字节数组
	bytes, err := serializedSegment(seg)
	if err != nil {
		t.Fatalf("failed to serialized segment:%v", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试数据
	_, err = tmpFile.Write(bytes)
	if err != nil {
		t.Fatalf("failed to write test data to temp file: %v", err)
	}

	// 使用 readSegment 读取并测试数据
	offset := uint64(0)
	inum, segment, err := readSegment(tmpFile, offset, 26)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// 校验 Segment 数据
	if segment.Tombstone != seg.Tombstone {
		t.Errorf("expected Tombstone to be %d, but got: %d", seg.Tombstone, segment.Tombstone)
	}
	if segment.Type != seg.Type {
		t.Errorf("expected Type to be %d, but got: %d", seg.Type, segment.Type)
	}
	if segment.ExpiredAt != seg.ExpiredAt {
		t.Errorf("expected ExpiredAt to be %d, but got: %d", seg.ExpiredAt, segment.ExpiredAt)
	}
	if segment.CreatedAt != seg.CreatedAt {
		t.Errorf("expected CreatedAt to be %d, but got: %d", seg.CreatedAt, segment.CreatedAt)
	}
	if segment.KeySize != seg.KeySize {
		t.Errorf("expected KeySize to be %d, but got: %d", seg.KeySize, segment.KeySize)
	}
	if segment.ValueSize != seg.ValueSize {
		t.Errorf("expected ValueSize to be %d, but got: %d", seg.ValueSize, segment.ValueSize)
	}
	if string(segment.Key) != string(seg.Key) {
		t.Errorf("expected Key to be %s, but got: %s", string(seg.Key), string(segment.Key))
	}
	if string(segment.Value) != string(seg.Value) {
		t.Errorf("expected Value to be %s, but got: %s", string(seg.Value), string(segment.Value))
	}

	// 校验返回的 inode number (InodeNum)
	if inum != InodeNum(string(seg.Key)) {
		t.Errorf("expected InodeNum to be '%s', but got: %d", seg.Key, inum)
	}
}

func TestVFSWrite(t *testing.T) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`

	var tables types.Tables
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		t.Fatal(err)
	}

	seg, err := NewSegment("key-01", tables, tables.TTL)
	if err != nil {
		t.Fatal(err)
	}

	err = fss.PutSegment("key-01", seg)
	if err != nil {
		t.Fatal(err)
	}

	_, seg, err = fss.FetchSegment("key-01")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", seg)
}

func BenchmarkVFSWrite(b *testing.B) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		b.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`
	var tables types.Tables
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		b.Fatal(err)
	}

	// 重置计时器，忽略 setup 代码的影响
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)

		seg, err := NewSegment(key, tables, tables.TTL)
		if err != nil {
			b.Fatal(err)
		}

		err = fss.PutSegment(key, seg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestUpdateSegmentWithCAS_Concurrent(t *testing.T) {
	var wg sync.WaitGroup

	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`
	var tables types.Tables
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		t.Fatal(err)
	}

	key := "key-01"
	seg, err := NewSegment(key, tables, tables.TTL)
	if err != nil {
		t.Fatal(err)
	}

	err = fss.PutSegment(key, seg)
	if err != nil {
		t.Fatal(err)
	}

	var failures int32
	var success int32

	concurrentUpdates := rand.Intn(1000)

	for i := 0; i < concurrentUpdates; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 1. 获取当前版本号
			version, _, err := fss.FetchSegment(key)
			if err != nil {
				t.Errorf("goroutine %d: failed to fetch segment: %v", id, err)
				return
			}

			// 2. 创建新的 `Segment`
			newseg, err := NewSegment(key, tables, tables.TTL)
			if err != nil {
				t.Errorf("goroutine %d: failed to create segment: %v", id, err)
				return
			}

			// 3. CAS 更新
			err = fss.UpdateSegmentWithCAS(key, version, newseg)
			if err != nil {
				atomic.AddInt32(&failures, 1)
				t.Logf("goroutine %d: CAS update failed (expected version: %d)", id, version)
			} else {
				atomic.AddInt32(&success, 1)
				t.Logf("goroutine %d: CAS update succeeded (version: %d)", id, version)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Total success: %d, Total failures: %d,Updates concurrent: %d", success, failures, concurrentUpdates)

	// 断言至少有一些失败的情况（正常情况下应该有很多失败）
	if failures == 0 && (failures+success) != int32(concurrentUpdates) {
		t.Error("Expected some CAS failures, but got none")
	}
}
