package dom

import (
	"time"

	"github.com/go-rod/rod"
)

// PageProvider - интерфейс для доступа к странице
type PageProvider interface {
	GetPage() *rod.Page
	GetTimeout() time.Duration
}
