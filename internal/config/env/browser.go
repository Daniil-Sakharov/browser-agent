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

func (b *BrowserConfig) Headless() bool {
	return b.headless
}

func (b *BrowserConfig) UserDataDir() string {
	return b.userDataDir
}

func (b *BrowserConfig) Timeout() int {
	return b.timeout
}

// NewBrowserConfig создает конфигурацию браузера из ENV
func NewBrowserConfig() (*BrowserConfig, error) {
	headless := os.Getenv("BROWSER_HEADLESS") == "true"
	userDataDir := os.Getenv("BROWSER_USER_DATA_DIR")
	if userDataDir == "" {
		userDataDir = ".browser-data"
	}

	timeout, err := strconv.Atoi(os.Getenv("BROWSER_TIMEOUT"))
	if err != nil || timeout == 0 {
		timeout = 30
	}

	return &BrowserConfig{
		headless:    headless,
		userDataDir: userDataDir,
		timeout:     timeout,
	}, nil
}
