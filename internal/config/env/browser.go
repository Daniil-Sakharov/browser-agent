package env

import "github.com/caarlos0/env/v11"

type browserEnvConfig struct {
	Headless    bool   `env:"BROWSER_HEADLESS" envDefault:"false"`
	UserDataDir string `env:"BROWSER_USER_DATA_DIR" envDefault:".browser-data"`
	Timeout     int    `env:"BROWSER_TIMEOUT" envDefault:"30"`
}

type browserConfig struct {
	raw browserEnvConfig
}

func NewBrowserConfig() (*browserConfig, error) {
	var raw browserEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &browserConfig{raw: raw}, nil
}

func (c *browserConfig) Headless() bool      { return c.raw.Headless }
func (c *browserConfig) UserDataDir() string { return c.raw.UserDataDir }
func (c *browserConfig) Timeout() int        { return c.raw.Timeout }
