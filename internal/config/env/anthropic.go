package env

import "github.com/caarlos0/env/v11"

type anthropicEnvConfig struct {
	APIKey      string  `env:"ANTHROPIC_API_KEY,required"`
	BaseURL     string  `env:"ANTHROPIC_BASE_URL"`
	Model       string  `env:"ANTHROPIC_MODEL" envDefault:"claude-sonnet-4-5-20250929"`
	MaxTokens   int     `env:"ANTHROPIC_MAX_TOKENS" envDefault:"4096"`
	Temperature float64 `env:"ANTHROPIC_TEMPERATURE" envDefault:"0"`
}

type anthropicConfig struct {
	raw anthropicEnvConfig
}

func NewAnthropicConfig() (*anthropicConfig, error) {
	var raw anthropicEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &anthropicConfig{raw: raw}, nil
}

func (c *anthropicConfig) APIKey() string       { return c.raw.APIKey }
func (c *anthropicConfig) BaseURL() string      { return c.raw.BaseURL }
func (c *anthropicConfig) Model() string        { return c.raw.Model }
func (c *anthropicConfig) MaxTokens() int       { return c.raw.MaxTokens }
func (c *anthropicConfig) Temperature() float64 { return c.raw.Temperature }
