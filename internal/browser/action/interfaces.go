package action

import (
	"time"

	"github.com/go-rod/rod"
)

// PageProvider - интерфейс для доступа к странице браузера
type PageProvider interface {
	GetPage() *rod.Page
	WaitStable(timeout time.Duration)
	GetTimeout() time.Duration
}
