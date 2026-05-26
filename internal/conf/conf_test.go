package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultConfig 验证默认值正确
func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	assert.Equal(t, "0.0.0.0:8080", cfg.Server.HTTP.Addr)
	assert.Equal(t, "30s", cfg.Server.HTTP.Timeout)
	assert.Equal(t, "sqlite", cfg.Database.Driver)
	assert.Equal(t, "./data/comic-translator.db", cfg.Database.DSN)
	assert.Equal(t, "chromem", cfg.VectorStore.Type)
	assert.Equal(t, "./data/chromem", cfg.VectorStore.Path)
	assert.Equal(t, "deepseek", cfg.LLM.Provider)
	assert.Equal(t, "DEEPSEEK_API_KEY", cfg.LLM.APIKeyEnv)
	assert.Equal(t, "deepseek-chat", cfg.LLM.Model)
	assert.Equal(t, "https://api.deepseek.com", cfg.LLM.BaseURL)
	assert.Equal(t, "gpt-image-2", cfg.ImageProcessing.Provider)
	assert.Equal(t, "OPENAI_API_KEY", cfg.ImageProcessing.APIKeyEnv)
	assert.Equal(t, "gpt-image-2", cfg.ImageProcessing.Model)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

// TestLoadYAMLIfExists_NotExist 文件不存在时不报错
func TestLoadYAMLIfExists_NotExist(t *testing.T) {
	var cfg Config
	err := loadYAMLIfExists("/non/existent/path/config.yaml", &cfg)
	assert.NoError(t, err)
	// 结构体保持零值
	assert.Empty(t, cfg.Server.HTTP.Addr)
}

// TestLoadYAMLIfExists_Exists 文件存在时正确解析
func TestLoadYAMLIfExists_Exists(t *testing.T) {
	content := `
logging:
  level: "debug"
  format: "text"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	var cfg Config
	err := loadYAMLIfExists(path, &cfg)
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
	// 未设置的字段应为零值
	assert.Empty(t, cfg.Server.HTTP.Addr)
}

// TestLoad_DefaultsOnly workspaceDir="" 时返回默认值
func TestLoad_DefaultsOnly(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 核心默认值应存在（可能被本机全局配置覆盖，但至少不为空）
	assert.NotEmpty(t, cfg.Server.HTTP.Addr)
	assert.NotEmpty(t, cfg.Database.Driver)
	assert.NotEmpty(t, cfg.LLM.Provider)
	assert.NotEmpty(t, cfg.Logging.Level)
}

// TestLoad_WorkspaceOverride workspace.yaml 覆盖默认值
func TestLoad_WorkspaceOverride(t *testing.T) {
	dir := t.TempDir()
	content := `
llm:
  model: "gpt-4o"
  provider: "openai"
logging:
  level: "warn"
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "workspace.yaml"), []byte(content), 0644))

	cfg, err := Load(dir)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// workspace.yaml 中设置的字段应被覆盖
	assert.Equal(t, "gpt-4o", cfg.LLM.Model)
	assert.Equal(t, "openai", cfg.LLM.Provider)
	assert.Equal(t, "warn", cfg.Logging.Level)

	// 未在 workspace.yaml 中设置的字段保留默认值
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.HTTP.Addr)
	assert.Equal(t, "sqlite", cfg.Database.Driver)
}

// TestLoad_EnvOverride 环境变量覆盖
func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("COMIC_TRANSLATOR_LLM_MODEL", "claude-3-5-sonnet")
	t.Setenv("COMIC_TRANSLATOR_LOGGING_LEVEL", "error")

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "claude-3-5-sonnet", cfg.LLM.Model)
	assert.Equal(t, "error", cfg.Logging.Level)
	// 未设置的字段保留默认/覆盖值
	assert.NotEmpty(t, cfg.Server.HTTP.Addr)
}

// TestLoad_MergePriority 验证优先级：workspace > global > default
// 通过 workspace 覆盖默认值来验证优先级（不写入 home 目录，避免副作用）
func TestLoad_MergePriority(t *testing.T) {
	// 创建 workspace 目录，其中 workspace.yaml 设置特殊值
	wsDir := t.TempDir()
	wsContent := `
database:
  driver: "postgres"
  dsn: "host=localhost port=5432"
llm:
  model: "workspace-model"
`
	require.NoError(t, os.WriteFile(filepath.Join(wsDir, "workspace.yaml"), []byte(wsContent), 0644))

	cfg, err := Load(wsDir)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// workspace 覆盖 > 默认值
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, "host=localhost port=5432", cfg.Database.DSN)
	assert.Equal(t, "workspace-model", cfg.LLM.Model)

	// 未被 workspace 覆盖的字段来自默认值
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.HTTP.Addr)
	assert.Equal(t, "deepseek", cfg.LLM.Provider)
}
