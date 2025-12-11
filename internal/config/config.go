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
	// Загружаем .env файл (если существует)
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	browserCfg, err := env.NewBrowserConfig()
	if err != nil {
		return err
	}

	anthropicCfg, err := env.NewAnthropicConfig()
	if err != nil {
		return err
	}

	agentCfg, err := env.NewAgentConfig()
	if err != nil {
		return err
	}

	securityCfg, err := env.NewSecurityConfig()
	if err != nil {
		return err
	}

	appConfig = &Config{
		Logger:    loggerCfg,
		Browser:   browserCfg,
		Anthropic: anthropicCfg,
		Agent:     agentCfg,
		Security:  securityCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию
func AppConfig() *Config {
	return appConfig
}
