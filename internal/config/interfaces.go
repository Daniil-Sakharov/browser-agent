package config

// LoggerConfig конфигурация логгера
type LoggerConfig interface {
	Level() string
	AsJson() bool
	LogFile() string
}

// BrowserConfig конфигурация браузера
type BrowserConfig interface {
	Headless() bool
	UserDataDir() string
	Timeout() int
}

// AnthropicConfig конфигурация Anthropic API
type AnthropicConfig interface {
	APIKey() string
	BaseURL() string
	Model() string
	MaxTokens() int
	Temperature() float64
}

// AgentConfig конфигурация агента
type AgentConfig interface {
	MaxSteps() int
	Interactive() bool
	Screenshots() bool
	ScreenshotsDir() string
}

// SecurityConfig конфигурация security checker
type SecurityConfig interface {
	Enabled() bool
	AutoConfirm() bool
}
