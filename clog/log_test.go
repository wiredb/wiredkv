package clog

import (
	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	// 在系统临时目录中创建一个临时文件
	tempFile := "./example-log.txt"
	defer os.Remove(tempFile) // 退出时删除

	SetOutput(tempFile)

	Info("info message.")

	Infof("info %s", "message.")

	Warn("warn message.")

	Warnf("warin %s", "message.")

	Error("error message.")

	Errorf("error %s", "message.")

}

// 测试 Failed 函数
func TestFailed(t *testing.T) {
	msg, panicked := capturePanic(func() {
		Failed("Test error message")
	})

	if !panicked {
		t.Errorf("Failed() did not panic as expected")
	}

	if msg == "" {
		t.Errorf("Failed() panic message is empty")
	} else {
		t.Logf("Captured panic message: %s", msg)
	}
}

// 测试 Failedf 函数
func TestFailedf(t *testing.T) {
	msg, panicked := capturePanic(func() {
		Failedf("Test formatted message: %d", 42)
	})

	if !panicked {
		t.Errorf("Failedf() did not panic as expected")
	}

	if msg == "" {
		t.Errorf("Failedf() panic message is empty")
	} else {
		t.Logf("Captured panic message: %s", msg)
	}
}

// 捕获 panic 并返回 panic 信息
func capturePanic(f func()) (message string, didPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			message, didPanic = r.(string), true
		}
	}()
	f()
	return "", false
}
