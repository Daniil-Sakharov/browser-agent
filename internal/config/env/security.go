package env

import "os"

type SecurityConfig struct {
	enabled     bool
	autoConfirm bool
}

func (s *SecurityConfig) Enabled() bool {
	return s.enabled
}

func (s *SecurityConfig) AutoConfirm() bool {
	return s.autoConfirm
}

// NewSecurityConfig создает конфигурацию security из ENV
func NewSecurityConfig() (*SecurityConfig, error) {
	enabled := os.Getenv("SECURITY_ENABLED") != "false"
	autoConfirm := os.Getenv("SECURITY_AUTO_CONFIRM") == "true"

	return &SecurityConfig{
		enabled:     enabled,
		autoConfirm: autoConfirm,
	}, nil
}
