package env

import (
	"os"
	"strconv"
)

type BrowserConfig struct {
	headless    bool
	userDataDir string
	timeout     int
}

func (b *BrowserConfig) Headless() bool    { return b.headless }
func (b *BrowserConfig) UserDataDir() string { return b.userDataDir }
func (b *BrowserConfig) Timeout() int      { return b.timeout }

func LoadBrowserConfig() *BrowserConfig {
	userDataDir := os.Getenv("BROWSER_USER_DATA_DIR")
	if userDataDir == "" {
		userDataDir = ".browser-data"
	}
	timeout, _ := strconv.Atoi(os.Getenv("BROWSER_TIMEOUT"))
	if timeout == 0 {
		timeout = 30
	}
	return &BrowserConfig{
		headless:    os.Getenv("BROWSER_HEADLESS") == "true",
		userDataDir: userDataDir,
		timeout:     timeout,
	}
}
