package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level   string `env:"LOG_LEVEL" envDefault:"info"`
	AsJson  bool   `env:"LOG_AS_JSON" envDefault:"false"`
	LogFile string `env:"LOG_FILE"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &loggerConfig{raw: raw}, nil
}

func (c *loggerConfig) Level() string   { return c.raw.Level }
func (c *loggerConfig) AsJson() bool    { return c.raw.AsJson }
func (c *loggerConfig) LogFile() string { return c.raw.LogFile }
