package env

import "github.com/caarlos0/env/v11"

type securityEnvConfig struct {
	Enabled     bool `env:"SECURITY_ENABLED" envDefault:"true"`
	AutoConfirm bool `env:"SECURITY_AUTO_CONFIRM" envDefault:"false"`
}

type securityConfig struct {
	raw securityEnvConfig
}

func NewSecurityConfig() (*securityConfig, error) {
	var raw securityEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &securityConfig{raw: raw}, nil
}

func (c *securityConfig) Enabled() bool     { return c.raw.Enabled }
func (c *securityConfig) AutoConfirm() bool { return c.raw.AutoConfirm }
