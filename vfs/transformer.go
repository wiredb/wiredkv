package vfs

import (
	"errors"
	"fmt"

	"github.com/golang/snappy"
)

const (
	// 使用整数位标志存储状态
	EnabledEncryption  = 1 << iota // 1: 0001
	EnabledCompression             // 2: 0010
)

// 压缩和解密应该针对数据的 VALUE ? 部分进行压缩，这里针对的是不定长部分进行压缩和解密
// | DEL 1 | KIND 1 | EAT 8 | CAT 8 | KLEN 4 | VLEN 4 | KEY ? | VALUE ? | CRC32 4 |
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type Encryptor interface {
	Encode(secret, data []byte) ([]byte, error)
	Decode(secret, data []byte) ([]byte, error)
}

type Transformer struct {
	Encryptor
	Compressor
	flags  int
	secret []byte
}

func NewTransformer() *Transformer {
	return &Transformer{
		flags:      0,
		Encryptor:  nil,
		Compressor: nil,
	}
}

func (t *Transformer) EnableEncryption() {
	t.flags |= EnabledEncryption
}

func (t *Transformer) EnableCompression() {
	t.flags |= EnabledCompression
}

func (t *Transformer) DisableEncryption() {
	t.flags &^= EnabledEncryption
}

func (t *Transformer) DisableCompression() {
	t.flags &^= EnabledCompression
}

func (t *Transformer) IsEncryptionEnabled() bool {
	return t.flags&EnabledEncryption != 0
}

func (t *Transformer) IsCompressionEnabled() bool {
	return t.flags&EnabledCompression != 0
}

func (t *Transformer) DisableAll() {
	t.flags = 0
}

func (t *Transformer) SetEncryptor(encryptor Encryptor, secret []byte) error {
	if len(secret) < 16 {
		return errors.New("secret char length too short")
	}
	t.secret = secret
	t.Encryptor = encryptor
	t.EnableEncryption()
	return nil
}

func (t *Transformer) SetCompressor(compressor Compressor) {
	t.Compressor = compressor
	t.EnableCompression()
}

func (t *Transformer) Encode(data []byte) ([]byte, error) {
	var err error
	// 压缩数据
	if t.IsCompressionEnabled() && t.Compressor != nil {
		data, err = t.Compress(data)
		if err != nil {
			return nil, fmt.Errorf("failed to compress data: %w", err)
		}
	}

	// 加密数据
	if t.IsEncryptionEnabled() && t.Encryptor != nil {
		data, err = t.Encryptor.Encode(t.secret, data)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", err)
		}
	}

	return data, nil
}

// fd 必须实现 io.ReadWriteCloser 接口
func (t *Transformer) Decode(data []byte) ([]byte, error) {
	var err error
	// 解密数据
	if t.IsEncryptionEnabled() && t.Encryptor != nil {
		data, err = t.Encryptor.Decode(t.secret, data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
	}

	// 解压缩数据
	if t.IsCompressionEnabled() && t.Compressor != nil {
		data, err = t.Decompress(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress data: %w", err)
		}
	}

	return data, nil
}

type SnappyCompressor struct{}

func (s *SnappyCompressor) Compress(data []byte) ([]byte, error) {
	// Snappy 压缩数据
	compressed := snappy.Encode(nil, data)
	return compressed, nil
}

func (s *SnappyCompressor) Decompress(data []byte) ([]byte, error) {
	// Snappy 解压数据
	return snappy.Decode(nil, data)
}
