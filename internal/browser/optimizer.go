package browser

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// optimizeElements оптимизирует список элементов (сортировка, без лимита)
func (e *Extractor) optimizeElements(ctx context.Context, elements []domain.Element) []domain.Element {
	// Фильтруем элементы без текста и селектора
	filtered := make([]domain.Element, 0, len(elements))
	for _, elem := range elements {
		if elem.Selector != "" && (elem.Text != "" || elem.Type == "input") {
			filtered = append(filtered, elem)
		}
	}

	// Сортируем по приоритету (кнопки > inputs > ссылки)
	prioritized := e.prioritizeElements(filtered)
	
	logger.Info(ctx, "📊 Elements", zap.Int("count", len(prioritized)))
	return prioritized
}

// prioritizeElements сортирует элементы по важности
func (e *Extractor) prioritizeElements(elements []domain.Element) []domain.Element {
	buttons := []domain.Element{}
	inputs := []domain.Element{}
	links := []domain.Element{}
	others := []domain.Element{}

	for _, elem := range elements {
		switch elem.Type {
		case "button":
			buttons = append(buttons, elem)
		case "input":
			inputs = append(inputs, elem)
		case "link":
			links = append(links, elem)
		default:
			others = append(others, elem)
		}
	}

	// Объединяем в порядке приоритета
	result := make([]domain.Element, 0, len(elements))
	result = append(result, buttons...)
	result = append(result, inputs...)
	result = append(result, links...)
	result = append(result, others...)

	return result
}

// extractVisibleText извлекает видимый текст страницы
func (e *Extractor) extractVisibleText(doc *goquery.Document) string {
	// Удаляем скрипты и стили
	doc.Find("script, style, noscript").Remove()

	// Получаем текст body
	bodyText := doc.Find("body").Text()

	// Очищаем от лишних пробелов и переносов
	lines := strings.Split(bodyText, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 2 {
			cleaned = append(cleaned, line)
		}
	}

	text := strings.Join(cleaned, " ")

	// Ограничиваем длину
	if len(text) > e.maxTextChars {
		text = text[:e.maxTextChars] + "..."
	}

	return text
}
