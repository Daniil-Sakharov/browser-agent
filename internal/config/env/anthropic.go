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

func (a *AnthropicConfig) APIKey() string      { return a.apiKey }
func (a *AnthropicConfig) BaseURL() string     { return a.baseURL }
func (a *AnthropicConfig) Model() string       { return a.model }
func (a *AnthropicConfig) MaxTokens() int      { return a.maxTokens }
func (a *AnthropicConfig) Temperature() float64 { return a.temperature }

func LoadAnthropicConfig() (*AnthropicConfig, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}
	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = "claude-sonnet-4-5-20250929"
	}
	maxTokens, _ := strconv.Atoi(os.Getenv("ANTHROPIC_MAX_TOKENS"))
	if maxTokens == 0 {
		maxTokens = 4096
	}
	temperature, _ := strconv.ParseFloat(os.Getenv("ANTHROPIC_TEMPERATURE"), 64)

	return &AnthropicConfig{
		apiKey:      apiKey,
		baseURL:     os.Getenv("ANTHROPIC_BASE_URL"),
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}
