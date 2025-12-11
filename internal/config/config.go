package config

import (
	"os"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	Browser   BrowserConfig
	Anthropic AnthropicConfig
	Agent     AgentConfig
	Security  SecurityConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil && !os.IsNotExist(err) {
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

	appConfig = &config{
		Logger:    loggerCfg,
		Browser:   browserCfg,
		Anthropic: anthropicCfg,
		Agent:     agentCfg,
		Security:  securityCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
