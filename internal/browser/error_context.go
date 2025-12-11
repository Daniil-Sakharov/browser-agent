package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

func (c *Controller) BuildErrorContext(ctx context.Context, failedSelector string, err error) *domain.ErrorContext {
	similar := c.findSimilarElements(ctx)
	return &domain.ErrorContext{
		FailedSelector:  failedSelector,
		SimilarElements: similar,
		Suggestion:      c.generateSuggestion(failedSelector, err, similar),
	}
}

func (c *Controller) findSimilarElements(ctx context.Context) []string {
	var results []string
	selectors := []string{"button", "a[href]", "input", "[role='button']", "[data-testid]"}

	for _, sel := range selectors {
		elements, err := c.page.Timeout(time.Second).Elements(sel)
		if err != nil { continue }

		for _, elem := range elements {
			result, err := elem.Eval(elemInfoJS)
			if err != nil { continue }
			if info := result.Value.String(); info != "" && info != `""` {
				results = append(results, info)
			}
			if len(results) >= 15 { break }
		}
		if len(results) >= 15 { break }
	}
	return results
}

const elemInfoJS = `() => {
	const el = this;
	const tag = el.tagName.toLowerCase();
	const id = el.id ? '#'+el.id : '';
	const cls = el.className && typeof el.className === 'string' ? '.'+el.className.split(' ').filter(c=>c).slice(0,2).join('.') : '';
	const tid = el.getAttribute('data-testid') ? '[data-testid="'+el.getAttribute('data-testid')+'"]' : '';
	const text = (el.innerText||'').substring(0,30).trim().replace(/\n/g,' ');
	let sel = tag + (tid || id || cls);
	return text ? sel+' ("'+text+'")' : sel;
}`

func (c *Controller) generateSuggestion(selector string, err error, similar []string) string {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "not found") || strings.Contains(errStr, "timeout"):
		if strings.Contains(selector, "modal") { return "Модальное окно закрыто. Открой его снова." }
		if len(similar) > 0 { return "Элемент не найден. Используй альтернативный селектор." }
		return "Элемент не найден. Страница загружается или элемент скрыт."
	case strings.Contains(errStr, "not visible"):
		return "Элемент скрыт. Прокрути страницу или открой меню."
	case strings.Contains(errStr, "not clickable"):
		return "Элемент перекрыт. Закрой модальное окно или прокрути."
	}
	return "Попробуй альтернативный селектор."
}

func FormatErrorContextMessage(ctx *domain.ErrorContext) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("❌ Не найден: %s\n\n", ctx.FailedSelector))
	if len(ctx.SimilarElements) > 0 {
		sb.WriteString("📋 Похожие элементы:\n")
		for i, e := range ctx.SimilarElements {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, e))
		}
	}
	sb.WriteString(fmt.Sprintf("\n💡 %s", ctx.Suggestion))
	return sb.String()
}
