package env

import (
	"fmt"
	"os"
	"strconv"
)

type AnthropicConfig struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	temperature float64
}

func (a *AnthropicConfig) APIKey() string {
	return a.apiKey
}

func (a *AnthropicConfig) BaseURL() string {
	return a.baseURL
}

func (a *AnthropicConfig) Model() string {
	return a.model
}

func (a *AnthropicConfig) MaxTokens() int {
	return a.maxTokens
}

func (a *AnthropicConfig) Temperature() float64 {
	return a.temperature
}

// NewAnthropicConfig создает конфигурацию Anthropic API из ENV
func NewAnthropicConfig() (*AnthropicConfig, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}

	baseURL := os.Getenv("ANTHROPIC_BASE_URL") // пустая строка = официальный API

	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = "claude-sonnet-4-5-20250929"
	}

	maxTokens, err := strconv.Atoi(os.Getenv("ANTHROPIC_MAX_TOKENS"))
	if err != nil || maxTokens == 0 {
		maxTokens = 4096
	}

	temperature, err := strconv.ParseFloat(os.Getenv("ANTHROPIC_TEMPERATURE"), 64)
	if err != nil {
		temperature = 0.0
	}

	return &AnthropicConfig{
		apiKey:      apiKey,
		baseURL:     baseURL,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}
