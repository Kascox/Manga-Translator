package conf

import (
	"os"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

// Config 顶层配置
type Config struct {
	Server          ServerConf          `yaml:"server" mapstructure:"server"`
	Database        DataConf            `yaml:"database" mapstructure:"database"`
	VectorStore     VectorStoreConf     `yaml:"vectorstore" mapstructure:"vectorstore"`
	LLM             LLMConf             `yaml:"llm" mapstructure:"llm"`
	ImageProcessing ImageProcessingConf `yaml:"image_processing" mapstructure:"image_processing"`
	Logging         LoggingConf         `yaml:"logging" mapstructure:"logging"`
}

// ServerConf HTTP 服务配置
type ServerConf struct {
	HTTP HTTPConf `yaml:"http" mapstructure:"http"`
}

// HTTPConf HTTP 监听与超时配置
type HTTPConf struct {
	Addr    string `yaml:"addr" mapstructure:"addr"`
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

// DataConf 数据库配置
type DataConf struct {
	Driver string `yaml:"driver" mapstructure:"driver"`
	DSN    string `yaml:"dsn" mapstructure:"dsn"`
}

// VectorStoreConf 向量存储配置
type VectorStoreConf struct {
	Type string `yaml:"type" mapstructure:"type"`
	Path string `yaml:"path" mapstructure:"path"`
}

// LLMConf 大语言模型配置
type LLMConf struct {
	Provider  string `yaml:"provider" mapstructure:"provider"`
	APIKeyEnv string `yaml:"api_key_env" mapstructure:"api_key_env"`
	Model     string `yaml:"model" mapstructure:"model"`
	BaseURL   string `yaml:"base_url" mapstructure:"base_url"`
}

// ImageProcessingConf 图片加工配置
type ImageProcessingConf struct {
	Provider  string `yaml:"provider" mapstructure:"provider"`
	APIKeyEnv string `yaml:"api_key_env" mapstructure:"api_key_env"`
	Model     string `yaml:"model" mapstructure:"model"`
}

// LoggingConf 日志配置
type LoggingConf struct {
	Level  string `yaml:"level" mapstructure:"level"`
	Format string `yaml:"format" mapstructure:"format"`
}

// defaultConfig 返回内置默认配置
func defaultConfig() Config {
	return Config{
		Server: ServerConf{HTTP: HTTPConf{
			Addr:    "0.0.0.0:8080",
			Timeout: "30s",
		}},
		Database: DataConf{
			Driver: "sqlite",
			DSN:    "./data/comic-translator.db",
		},
		VectorStore: VectorStoreConf{
			Type: "chromem",
			Path: "./data/chromem",
		},
		LLM: LLMConf{
			Provider:  "deepseek",
			APIKeyEnv: "DEEPSEEK_API_KEY",
			Model:     "deepseek-chat",
			BaseURL:   "https://api.deepseek.com",
		},
		ImageProcessing: ImageProcessingConf{
			Provider:  "gpt-image-2",
			APIKeyEnv: "OPENAI_API_KEY",
			Model:     "gpt-image-2",
		},
		Logging: LoggingConf{
			Level:  "info",
			Format: "json",
		},
	}
}

// loadYAMLIfExists 若文件存在则将其反序列化到 v，否则静默跳过
func loadYAMLIfExists(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return yaml.Unmarshal(data, v)
}

// applyEnvOverrides 将 COMIC_TRANSLATOR_ 前缀的环境变量覆盖到 cfg
// 映射规则：
//
//	COMIC_TRANSLATOR_SERVER_HTTP_ADDR    → cfg.Server.HTTP.Addr
//	COMIC_TRANSLATOR_SERVER_HTTP_TIMEOUT → cfg.Server.HTTP.Timeout
//	COMIC_TRANSLATOR_DATABASE_DRIVER     → cfg.Database.Driver
//	COMIC_TRANSLATOR_DATABASE_DSN        → cfg.Database.DSN
//	COMIC_TRANSLATOR_VECTORSTORE_TYPE    → cfg.VectorStore.Type
//	COMIC_TRANSLATOR_VECTORSTORE_PATH    → cfg.VectorStore.Path
//	COMIC_TRANSLATOR_LLM_PROVIDER        → cfg.LLM.Provider
//	COMIC_TRANSLATOR_LLM_API_KEY_ENV     → cfg.LLM.APIKeyEnv
//	COMIC_TRANSLATOR_LLM_MODEL           → cfg.LLM.Model
//	COMIC_TRANSLATOR_LLM_BASE_URL        → cfg.LLM.BaseURL
//	COMIC_TRANSLATOR_IMAGE_PROVIDER      → cfg.ImageProcessing.Provider
//	COMIC_TRANSLATOR_IMAGE_API_KEY_ENV   → cfg.ImageProcessing.APIKeyEnv
//	COMIC_TRANSLATOR_IMAGE_MODEL         → cfg.ImageProcessing.Model
//	COMIC_TRANSLATOR_LOGGING_LEVEL       → cfg.Logging.Level
//	COMIC_TRANSLATOR_LOGGING_FORMAT      → cfg.Logging.Format
func applyEnvOverrides(cfg *Config) {
	setIfFound := func(dest *string, key string) {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			*dest = v
		}
	}

	setIfFound(&cfg.Server.HTTP.Addr, "COMIC_TRANSLATOR_SERVER_HTTP_ADDR")
	setIfFound(&cfg.Server.HTTP.Timeout, "COMIC_TRANSLATOR_SERVER_HTTP_TIMEOUT")
	setIfFound(&cfg.Database.Driver, "COMIC_TRANSLATOR_DATABASE_DRIVER")
	setIfFound(&cfg.Database.DSN, "COMIC_TRANSLATOR_DATABASE_DSN")
	setIfFound(&cfg.VectorStore.Type, "COMIC_TRANSLATOR_VECTORSTORE_TYPE")
	setIfFound(&cfg.VectorStore.Path, "COMIC_TRANSLATOR_VECTORSTORE_PATH")
	setIfFound(&cfg.LLM.Provider, "COMIC_TRANSLATOR_LLM_PROVIDER")
	setIfFound(&cfg.LLM.APIKeyEnv, "COMIC_TRANSLATOR_LLM_API_KEY_ENV")
	setIfFound(&cfg.LLM.Model, "COMIC_TRANSLATOR_LLM_MODEL")
	setIfFound(&cfg.LLM.BaseURL, "COMIC_TRANSLATOR_LLM_BASE_URL")
	setIfFound(&cfg.ImageProcessing.Provider, "COMIC_TRANSLATOR_IMAGE_PROVIDER")
	setIfFound(&cfg.ImageProcessing.APIKeyEnv, "COMIC_TRANSLATOR_IMAGE_API_KEY_ENV")
	setIfFound(&cfg.ImageProcessing.Model, "COMIC_TRANSLATOR_IMAGE_MODEL")
	setIfFound(&cfg.Logging.Level, "COMIC_TRANSLATOR_LOGGING_LEVEL")
	setIfFound(&cfg.Logging.Format, "COMIC_TRANSLATOR_LOGGING_FORMAT")
}

// Load 加载配置，workspaceDir 为当前工作区目录（可为空字符串）
// 优先级（低→高）：默认值 → 全局配置 → 工作区配置 → 环境变量
func Load(workspaceDir string) (*Config, error) {
	// 1. 从默认值开始
	cfg := defaultConfig()

	// 2. 全局配置：~/.comic-translator/config.yaml
	homeDir, err := os.UserHomeDir()
	if err == nil {
		var globalCfg Config
		globalPath := homeDir + "/.comic-translator/config.yaml"
		if err := loadYAMLIfExists(globalPath, &globalCfg); err != nil {
			return nil, err
		}
		if err := mergo.Merge(&cfg, globalCfg, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	// 3. 工作区配置：<workspaceDir>/workspace.yaml
	if workspaceDir != "" {
		var workspaceCfg Config
		workspacePath := workspaceDir + "/workspace.yaml"
		if err := loadYAMLIfExists(workspacePath, &workspaceCfg); err != nil {
			return nil, err
		}
		if err := mergo.Merge(&cfg, workspaceCfg, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	// 4. 环境变量覆盖
	applyEnvOverrides(&cfg)

	return &cfg, nil
}
