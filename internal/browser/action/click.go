package action

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// Click выполняет умный клик с цепочкой fallback
func Click(ctx context.Context, p PageProvider, selector string) error {
	page := p.GetPage()

	if strings.HasPrefix(selector, "text:") {
		text := strings.TrimPrefix(selector, "text:")
		return smartClickText(ctx, p, page, text)
	}
	return smartClickCSS(ctx, p, page, selector)
}

// smartClickText - умный клик по тексту с fallback цепочкой
func smartClickText(ctx context.Context, p PageProvider, page *rod.Page, text string) error {
	// Убираем ... в конце если есть
	text = strings.TrimSuffix(text, "...")

	attempts := []struct {
		name string
		fn   func() error
	}{
		{"exact", func() error { return tryClickText(ctx, p, page, text, true) }},
		{"partial", func() error { return tryClickText(ctx, p, page, text, false) }},
		{"short", func() error { return tryClickText(ctx, p, page, getShortText(text), false) }},
		{"js_smart", func() error { return jsSmartClick(ctx, page, text) }},
	}

	for _, a := range attempts {
		logger.Debug(ctx, "🔄 Trying click", zap.String("method", a.name), zap.String("text", text))
		if err := a.fn(); err == nil {
			logger.Info(ctx, "✅ Click success", zap.String("method", a.name), zap.String("text", text))
			return nil
		}
	}

	return fmt.Errorf("element not found: text:%s", text)
}

// smartClickCSS - умный клик по CSS с fallback
func smartClickCSS(ctx context.Context, p PageProvider, page *rod.Page, selector string) error {
	attempts := []struct {
		name string
		fn   func() error
	}{
		{"rod", func() error { return tryClickCSS(ctx, p, page, selector) }},
		{"js", func() error { return jsClickCSS(page, selector) }},
	}

	for _, a := range attempts {
		if err := a.fn(); err == nil {
			logger.Info(ctx, "✅ Click CSS", zap.String("method", a.name), zap.String("selector", selector))
			return nil
		}
	}
	return fmt.Errorf("element not found: %s", selector)
}

// tryClickText пробует кликнуть по тексту через Rod
func tryClickText(ctx context.Context, p PageProvider, page *rod.Page, text string, exact bool) error {
	elem, err := findElementByText(page, text, exact)
	if err != nil {
		return err
	}
	return doClick(ctx, p, elem)
}

// tryClickCSS пробует кликнуть по CSS через Rod
func tryClickCSS(ctx context.Context, p PageProvider, page *rod.Page, selector string) error {
	elem, err := page.Timeout(5 * time.Second).Element(selector)
	if err != nil {
		return err
	}
	return doClick(ctx, p, elem)
}

// doClick выполняет клик с hover и scroll
func doClick(ctx context.Context, p PageProvider, elem *rod.Element) error {
	// Scroll к элементу
	if err := elem.ScrollIntoView(); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	// Hover перед кликом (важно для dropdown/меню)
	elem.Hover()
	time.Sleep(50 * time.Millisecond)

	// Проверяем видимость
	if err := elem.Timeout(3 * time.Second).WaitVisible(); err != nil {
		return fmt.Errorf("not visible")
	}

	// Клик
	if err := elem.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	p.WaitStable(3 * time.Second)
	return nil
}

// findElementByText ищет элемент по тексту
func findElementByText(page *rod.Page, text string, exact bool) (*rod.Element, error) {
	matchType := "includes"
	if exact {
		matchType = "exact"
	}

	js := `(text, matchType) => {
		const t = text.toLowerCase();
		const selectors = ['button','a','[role="button"]','[role="link"]','div[onclick]','span[onclick]','li','label','h1','h2','h3','h4'];
		
		for (const sel of selectors) {
			for (const el of document.querySelectorAll(sel)) {
				const rect = el.getBoundingClientRect();
				if (rect.width === 0 || rect.height === 0) continue;
				
				const elText = (el.innerText || el.textContent || '').trim().toLowerCase();
				const ariaLabel = (el.getAttribute('aria-label') || '').toLowerCase();
				const title = (el.getAttribute('title') || '').toLowerCase();
				
				let match = false;
				if (matchType === 'exact') {
					match = elText === t || ariaLabel === t;
				} else {
					match = elText.includes(t) || t.includes(elText.substring(0, 20)) || ariaLabel.includes(t) || title.includes(t);
				}
				
				if (match && elText.length > 0 && elText.length < 200) {
					return el;
				}
			}
		}
		return null;
	}`

	result, err := page.Timeout(5*time.Second).Eval(js, text, matchType)
	if err != nil || result.Value.Nil() {
		return nil, fmt.Errorf("not found")
	}

	return page.ElementFromObject(&proto.RuntimeRemoteObject{
		ObjectID: proto.RuntimeRemoteObjectID(result.Value.Get("objectId").String()),
	})
}

// jsSmartClick - умный JS клик с несколькими стратегиями
func jsSmartClick(ctx context.Context, page *rod.Page, text string) error {
	js := `(searchText) => {
		const t = searchText.toLowerCase();
		const words = t.split(' ').filter(w => w.length > 2);
		
		// Стратегия 1: точное совпадение
		// Стратегия 2: начинается с текста
		// Стратегия 3: содержит все слова
		const strategies = [
			el => (el.innerText||'').trim().toLowerCase() === t,
			el => (el.innerText||'').trim().toLowerCase().startsWith(t),
			el => words.every(w => (el.innerText||'').toLowerCase().includes(w)),
			el => (el.getAttribute('aria-label')||'').toLowerCase().includes(t),
		];
		
		const selectors = ['a','button','[role="button"]','div','span','li'];
		
		for (const strategy of strategies) {
			for (const sel of selectors) {
				for (const el of document.querySelectorAll(sel)) {
					const rect = el.getBoundingClientRect();
					if (rect.width === 0 || rect.height === 0) continue;
					if (rect.top < 0 || rect.top > window.innerHeight) continue;
					
					if (strategy(el)) {
						el.scrollIntoView({block: 'center'});
						el.click();
						return {ok: true, text: (el.innerText||'').substring(0,50)};
					}
				}
			}
		}
		return {ok: false};
	}`

	result, err := page.Eval(js, text)
	if err != nil || result == nil || !result.Value.Get("ok").Bool() {
		return fmt.Errorf("js click failed")
	}

	logger.Info(ctx, "🎯 JS smart click", zap.String("found", result.Value.Get("text").String()))
	return nil
}

// jsClickCSS - JS клик по CSS селектору
func jsClickCSS(page *rod.Page, selector string) error {
	js := `(sel) => {
		const el = document.querySelector(sel);
		if (!el) return {ok: false};
		el.scrollIntoView({block: 'center'});
		el.click();
		return {ok: true};
	}`
	result, _ := page.Eval(js, selector)
	if result == nil || !result.Value.Get("ok").Bool() {
		return fmt.Errorf("js css click failed")
	}
	return nil
}

// getShortText возвращает первые 2-3 слова
func getShortText(text string) string {
	words := strings.Fields(text)
	if len(words) <= 2 {
		return text
	}
	return strings.Join(words[:2], " ")
}
