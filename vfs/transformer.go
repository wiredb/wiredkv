package vfs

import (
	"fmt"
	"io"
)

type Encryptor interface {
	Encode(secret, data []byte) ([]byte, error)
	Decode(secret, data []byte) ([]byte, error)
}

type Transformer struct {
	Encryptor
	enable bool
	secret []byte
}

func NewTransformer() *Transformer {
	return &Transformer{
		enable:    false,
		Encryptor: nil,
	}
}

func (t *Transformer) SetEncryptor(encryptor Encryptor, secret []byte) {
	t.enable = true
	t.secret = secret
	t.Encryptor = encryptor
}

// fd 必须实现 io.ReadWriteCloser 接口
func (t *Transformer) Write(fd io.ReadWriteCloser, data []byte) (int, error) {
	if t.enable && t.Encryptor != nil {
		bytes, err := t.Encode(t.secret, data)
		if err != nil {
			return 0, fmt.Errorf("failed to encrypted data: %w", err)
		}
		n, err := fd.Write(bytes)
		if err != nil {
			return 0, fmt.Errorf("failed to write encrypted data: %w", err)
		}
		return n, nil
	}

	// 写入数据到 fd
	return fd.Write(data)
}
