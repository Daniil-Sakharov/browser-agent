package env

import "github.com/caarlos0/env/v11"

type agentEnvConfig struct {
	MaxSteps       int    `env:"AGENT_MAX_STEPS" envDefault:"30"`
	Interactive    bool   `env:"AGENT_INTERACTIVE" envDefault:"true"`
	Screenshots    bool   `env:"AGENT_SCREENSHOTS" envDefault:"true"`
	ScreenshotsDir string `env:"AGENT_SCREENSHOTS_DIR" envDefault:"screenshots"`
}

type agentConfig struct {
	raw agentEnvConfig
}

func NewAgentConfig() (*agentConfig, error) {
	var raw agentEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &agentConfig{raw: raw}, nil
}

func (c *agentConfig) MaxSteps() int          { return c.raw.MaxSteps }
func (c *agentConfig) Interactive() bool      { return c.raw.Interactive }
func (c *agentConfig) Screenshots() bool      { return c.raw.Screenshots }
func (c *agentConfig) ScreenshotsDir() string { return c.raw.ScreenshotsDir }
