package env

import (
	"os"
	"strconv"
)

type AgentConfig struct {
	maxSteps       int
	interactive    bool
	screenshots    bool
	screenshotsDir string
}

func (a *AgentConfig) MaxSteps() int          { return a.maxSteps }
func (a *AgentConfig) Interactive() bool      { return a.interactive }
func (a *AgentConfig) Screenshots() bool      { return a.screenshots }
func (a *AgentConfig) ScreenshotsDir() string { return a.screenshotsDir }

func LoadAgentConfig() *AgentConfig {
	maxSteps, _ := strconv.Atoi(os.Getenv("AGENT_MAX_STEPS"))
	if maxSteps == 0 {
		maxSteps = 30
	}
	screenshotsDir := os.Getenv("AGENT_SCREENSHOTS_DIR")
	if screenshotsDir == "" {
		screenshotsDir = "screenshots"
	}
	return &AgentConfig{
		maxSteps:       maxSteps,
		interactive:    os.Getenv("AGENT_INTERACTIVE") != "false",
		screenshots:    os.Getenv("AGENT_SCREENSHOTS") != "false",
		screenshotsDir: screenshotsDir,
	}
}
