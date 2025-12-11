package config

import (
	"os"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *Config

// Config глобальная конфигурация приложения
type Config struct {
	Logger    LoggerConfig
	Browser   BrowserConfig
	Anthropic AnthropicConfig
	Agent     AgentConfig
	Security  SecurityConfig
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil && !os.IsNotExist(err) {
		return err
	}

	anthropicCfg, err := env.LoadAnthropicConfig()
	if err != nil {
		return err
	}

	appConfig = &Config{
		Logger:    env.LoadLoggerConfig(),
		Browser:   env.LoadBrowserConfig(),
		Anthropic: anthropicCfg,
		Agent:     env.LoadAgentConfig(),
		Security:  env.LoadSecurityConfig(),
	}
	return nil
}

// AppConfig возвращает глобальную конфигурацию
func AppConfig() *Config {
	return appConfig
}
