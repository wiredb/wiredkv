package conf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoad(t *testing.T) {
	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	// 设置 Settings.Path 为临时目录
	Settings.Path = tmpDir

	// 创建一个配置文件并写入测试数据
	configFile := filepath.Join(tmpDir, "test-config.yaml")
	testConfigData := []byte(`
  port: 8080
  path: "/test/path"
  debug: true
`)

	err := os.WriteFile(configFile, testConfigData, 0644)
	if err != nil {
		t.Fatalf("Error writing test config file: %v", err)
	}

	// 调用 Load 函数
	loadedConfig := new(ServerOptions)
	err = Load(configFile, loadedConfig)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// 检查加载的配置是否正确
	expectedConfig := &ServerOptions{
		Port:  8080,
		Path:  "/test/path",
		Debug: true,
	}

	// 检查比较是否一致
	if !reflect.DeepEqual(loadedConfig, expectedConfig) {
		t.Errorf("Loaded config is not as expected.\nGot: %+v\nExpected: %+v", loadedConfig, expectedConfig)
	}
}

func TestConfigLoad_Error(t *testing.T) {

	// 创建一个临时目录用于测试
	Settings.Path = t.TempDir() + "/aaa/bbb"

	// 创建一个配置文件并写入测试数据
	configFile := filepath.Join(Settings.Path, "test-config.yaml")

	// 调用 Load 函数
	loadedConfig := new(ServerOptions)
	err := Load(configFile, loadedConfig)
	if err != nil {
		t.Log(err)
	}

}

func TestSavedAsConfig(t *testing.T) {

	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	// 创建一个 ServerOptions 实例
	config := &ServerOptions{
		Port:     8080,
		Path:     tmpDir,
		Debug:    true,
		LogPath:  "/tmp/wiredb/out.log",
		Password: "password@123",
		Region: Region{
			Enable:    true,
			Second:    15000,
			Threshold: 3,
		},
		Encryptor: Encryptor{
			Enable: false,
			Secret: "test-secret",
		},
		Compressor: Compressor{
			Enable: false,
		},
		AllowIP: []string{"192.127.0.1", "192.127.0.2"},
	}

	_, err := os.Create(tmpDir + "/config.yaml")
	if err != nil {
		t.Error(err)
	}

	// 调用 Saved 函数
	err = config.SavedAs(tmpDir + "/config.yaml")

	if err != nil {
		t.Fatalf("Error saving config: %v", err)
	}
}

func TestSavedConfig(t *testing.T) {

	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "etc"), FSPerm)

	// 创建一个 ServerOptions 实例
	config := &ServerOptions{

		Port:     8080,
		Path:     tmpDir,
		Debug:    true,
		LogPath:  "/tmp/wiredb/out.log",
		Password: "password@123",
		Region: Region{
			Enable:    true,
			Second:    15000,
			Threshold: 3,
		},
		Encryptor: Encryptor{
			Enable: false,
			Secret: "test-secret",
		},
		Compressor: Compressor{
			Enable: false,
		},
	}

	// 调用 Saved 函数
	err := config.Saved()

	if err != nil {
		t.Fatalf("Error saving config: %v", err)
	}
}

func TestSavedConfig_Error(t *testing.T) {

	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	// 创建一个 ServerOptions 空实例
	var config *ServerOptions = nil

	// 调用 Saved 函数
	err := config.SavedAs(tmpDir)

	if err != nil {
		t.Log(err)
	}
}

func TestIsDefault(t *testing.T) {
	tests := []struct {
		name string
		flag string
		want bool
	}{
		{
			name: "successful", flag: "default.yaml", want: true,
		},
		{
			name: "failed", flag: "", want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasCustom(tt.flag); got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	t.Run("Test DefaultConfig Unmarshal", func(t *testing.T) {
		err := Default.Unmarshal([]byte(DefaultConfigJSON))
		if err != nil {
			t.Log(err)
		}
	})

	t.Run("Test Settings Unmarshal", func(t *testing.T) {
		err := Settings.Unmarshal([]byte(DefaultConfigJSON))
		if err != nil {
			t.Log(err)
		}
	})

}

func TestServerOptions_Marshal(t *testing.T) {

	err := Settings.Unmarshal([]byte(DefaultConfigJSON))
	if err != nil {
		t.Error(err)
	}

	bytes, err := Settings.Marshal()

	if err != nil {
		t.Error(err)
	}

	if err := Default.Unmarshal(bytes); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(Settings, Default) {
		t.Errorf("ServerOptions.Marshal() = %v, want %v", string(bytes), DefaultConfigJSON)
	}

}

func TestDefaultConfigInitialization(t *testing.T) {

	// 检查 DefaultConfig 是否被正确初始化
	if Default.Port != 2668 {
		t.Errorf("Expected DefaultConfig.Port to be 2668, but got %d", Default.Port)
	}

	// 检查 Settings 是否被正确初始化
	if Settings.Port != 2668 {
		t.Errorf("Expected Settings.Port to be 2668, but got %d", Settings.Port)
	}

}

func TestServerOptions_ToString(t *testing.T) {

	type fields struct {
		TestDB ServerOptions
	}

	vdb := ServerOptions{
		Port:     8080,
		Path:     "",
		Debug:    true,
		LogPath:  "/tmp/wiredb/out.log",
		Password: "password@123",
		Region: Region{
			Enable:    true,
			Second:    15000,
			Threshold: 3,
		},
		Encryptor: Encryptor{
			Enable: false,
			Secret: "test-secret",
		},
		Compressor: Compressor{
			Enable: false,
		},
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "successful", fields: fields{TestDB: vdb}, want: vdb.String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.TestDB.String(); got != tt.want {
				t.Errorf("ServerOptions.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasCustom tests HasCustom method to check if custom config is provided
func TestHasCustom(t *testing.T) {
	assert.True(t, HasCustom("/path/to/custom/config.yaml"))
	assert.False(t, HasCustom(defaultFilePath))
}

// TestVaildated tests the configuration validation
func TestVaildated(t *testing.T) {
	// Valid configuration
	validConfig := &ServerOptions{
		Port:     2668,
		Path:     "/tmp/wiredb",
		Password: "securepassword",
		Region: Region{
			Enable:    true,
			Second:    18000,
			Threshold: 3,
		},
		Encryptor: Encryptor{
			Enable: true,
			Secret: "1234567890123456",
		},
		Compressor: Compressor{
			Enable: false,
		},
	}

	// Should pass validation
	err := Vaildated(validConfig)
	assert.NoError(t, err)

	// Invalid configuration: port out of range
	invalidConfig := *validConfig
	invalidConfig.Port = 70000 // Invalid port number

	err = Vaildated(&invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port range must be between 1025 and 65534")

	// Invalid configuration: empty path
	invalidConfig = *validConfig
	invalidConfig.Path = ""

	err = Vaildated(&invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data directory path cannot be empty")

	// Invalid configuration: invalid secret key length
	invalidConfig = *validConfig
	invalidConfig.Encryptor.Secret = "short" // Invalid secret key length

	err = Vaildated(&invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid secret key length it must be 16, 24, or 32 bytes")
}

// TestSaved tests saving the configuration to a file
func TestSaved(t *testing.T) {
	// Prepare server options to save
	opt := &ServerOptions{
		Port:     8080,
		Path:     "/tmp/myconfig",
		Password: "testpassword",
	}

	// Use a temporary directory to save the config
	tempFile := filepath.Join(t.TempDir(), "test_config.yaml")
	err := opt.SavedAs(tempFile)
	require.NoError(t, err)

	// Check if the file exists
	_, err = os.Stat(tempFile)
	require.NoError(t, err)

	// Check if the file contains valid data
	var loadedOpt ServerOptions
	err = Load(tempFile, &loadedOpt)
	require.NoError(t, err)

	// Ensure the saved config matches the original
	assert.Equal(t, opt.Port, loadedOpt.Port)
	assert.Equal(t, opt.Path, loadedOpt.Path)
	assert.Equal(t, opt.Password, loadedOpt.Password)
}

// TestMarshal tests the Marshal method to convert ServerOptions to JSON
func TestMarshal(t *testing.T) {
	opt := &ServerOptions{
		Port:     8080,
		Path:     "/tmp/myconfig",
		Password: "testpassword",
	}

	// Marshal the options to JSON
	data, err := opt.Marshal()
	require.NoError(t, err)

	// Verify the marshaled data is correct
	expectedJSON := `{"port":8080,"path":"/tmp/myconfig","debug":false,"logpath":"","auth":"testpassword","region":{"enable":false,"second":0,"threshold":0},"encryptor":{"enable":false,"secret":""},"compressor":{"enable":false},"allowip":null}`
	assert.JSONEq(t, expectedJSON, string(data))
}

// TestString tests the String method to convert ServerOptions to string
func TestString(t *testing.T) {
	opt := &ServerOptions{
		Port:     8080,
		Path:     "/tmp/myconfig",
		Password: "testpassword",
	}

	// Get the string representation of the options
	str := opt.String()

	// Ensure the string contains expected values
	assert.Contains(t, str, "8080")
	assert.Contains(t, str, "/tmp/myconfig")
	assert.Contains(t, str, "testpassword")
}
