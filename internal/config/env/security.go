package env

import "os"

type SecurityConfig struct {
	enabled     bool
	autoConfirm bool
}

func (s *SecurityConfig) Enabled() bool     { return s.enabled }
func (s *SecurityConfig) AutoConfirm() bool { return s.autoConfirm }

func LoadSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		enabled:     os.Getenv("SECURITY_ENABLED") != "false",
		autoConfirm: os.Getenv("SECURITY_AUTO_CONFIRM") == "true",
	}
}
