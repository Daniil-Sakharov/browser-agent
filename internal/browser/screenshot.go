package browser

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// ScreenshotResult результат создания скриншота
type ScreenshotResult struct {
	Path   string // путь к файлу
	Base64 string // base64 encoded image
}

// TakeScreenshot делает скриншот страницы
func (c *Controller) TakeScreenshot(ctx context.Context, fullPage bool, saveDir string) (*ScreenshotResult, error) {
	logger.Info(ctx, "📸 Taking screenshot", zap.Bool("full_page", fullPage))

	if c.page == nil {
		return nil, fmt.Errorf("page is nil")
	}

	var data []byte
	var err error

	if fullPage {
		// Скриншот всей страницы
		data, err = c.page.Timeout(c.timeout).Screenshot(true, nil)
	} else {
		// Скриншот только видимой части
		data, err = c.page.Timeout(c.timeout).Screenshot(false, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to take screenshot: %w", err)
	}

	// Создаем директорию если не существует
	if saveDir != "" {
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			logger.Warn(ctx, "⚠️ Failed to create screenshots dir", zap.Error(err))
		}
	} else {
		saveDir = "screenshots"
		os.MkdirAll(saveDir, 0755)
	}

	// Генерируем имя файла
	filename := fmt.Sprintf("screenshot-%d.png", time.Now().UnixMilli())
	filePath := filepath.Join(saveDir, filename)

	// Сохраняем файл
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		logger.Warn(ctx, "⚠️ Failed to save screenshot", zap.Error(err))
	}

	// Кодируем в base64 для передачи в Claude
	b64 := base64.StdEncoding.EncodeToString(data)

	logger.Info(ctx, "✅ Screenshot taken",
		zap.String("path", filePath),
		zap.Int("size_bytes", len(data)))

	return &ScreenshotResult{
		Path:   filePath,
		Base64: b64,
	}, nil
}
