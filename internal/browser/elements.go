package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"github.com/ysmood/gson"
	"go.uber.org/zap"
)

// FindElementsLive возвращает кликабельные элементы со страницы
func (c *Controller) FindElementsLive(ctx context.Context, query string) (string, error) {
	result, err := c.page.Timeout(5 * time.Second).Eval(findElementsJS)
	if err != nil {
		logger.Error(ctx, "❌ FindElementsLive error", zap.Error(err))
		return "Ошибка поиска.", nil
	}
	elements := result.Value.Arr()
	if len(elements) == 0 {
		return "Элементы не найдены.", nil
	}
	logger.Info(ctx, "🔍 FindElementsLive", zap.Int("found", len(elements)))
	return formatElements(elements), nil
}

func formatElements(elements []gson.JSON) string {
	var out strings.Builder
	out.WriteString("🎯 КЛИКАБЕЛЬНЫЕ ЭЛЕМЕНТЫ:\n\n")
	for i, elem := range elements {
		obj := elem.Map()
		sel := obj["displaySelector"].String()
		if sel == "" {
			continue
		}
		out.WriteString(fmt.Sprintf("%d. %s", i+1, sel))
		if css := obj["cssSelector"].String(); css != "" && css != sel {
			out.WriteString(fmt.Sprintf(" [%s]", css))
		}
		out.WriteString("\n")
	}
	out.WriteString("\n💡 Используй text:ТекстКнопки для клика")
	return out.String()
}

const findElementsJS = `() => {
	const results = [], seen = new Set(), counts = {};
	function isVisible(el) {
		if (!el) return false;
		const s = window.getComputedStyle(el), r = el.getBoundingClientRect();
		return s.display !== 'none' && s.visibility !== 'hidden' && s.opacity !== '0' && r.width > 0 && r.height > 0 && r.top < window.innerHeight && r.bottom > 0;
	}
	function getSelector(el) {
		const tag = el.tagName.toLowerCase();
		if (el.id && el.id.length < 50) return tag + '#' + el.id;
		const aria = el.getAttribute('aria-label');
		if (aria && aria.length < 100) return tag + '[aria-label="' + aria.replace(/"/g, '\\"') + '"]';
		const tid = el.getAttribute('data-testid');
		if (tid) return tag + '[data-testid="' + tid + '"]';
		if (el.name) return tag + '[name="' + el.name + '"]';
		return tag;
	}
	function getText(el) {
		let t = '';
		for (const n of el.childNodes) if (n.nodeType === Node.TEXT_NODE) t += n.textContent;
		t = t.trim();
		if (t && t.length > 0 && t.length < 50) return t;
		return (el.innerText || '').trim().split('\n')[0].substring(0, 50);
	}
	const all = [], clickable = ['button','a[href]','[role="button"]','[role="menuitem"]','input:not([type="hidden"])','textarea','select','[onclick]','li','label'];
	const modals = ['[role="dialog"]','[aria-modal="true"]','.modal','[class*="modal"]'];
	let hasModal = false;
	for (const ms of modals) { try { for (const m of document.querySelectorAll(ms)) { if (!isVisible(m)) continue; hasModal = true;
		for (const el of m.querySelectorAll(clickable.join(','))) { if (isVisible(el)) { const s = getSelector(el); counts[s] = (counts[s] || 0) + 1; all.push({ el, s }); }}
	}} catch(e){} }
	if (!hasModal) { for (const sel of clickable) { try { for (const el of document.querySelectorAll(sel)) { if (isVisible(el)) { const s = getSelector(el); counts[s] = (counts[s] || 0) + 1; all.push({ el, s }); }}} catch(e){} }}
	for (const item of all) {
		const el = item.el, css = item.s, text = getText(el), r = el.getBoundingClientRect(), key = css + '|' + text;
		if (seen.has(key)) continue; seen.add(key);
		results.push({ displaySelector: text ? 'text:' + text : css, cssSelector: counts[css] === 1 ? css : '', x: Math.round(r.x + r.width / 2), y: Math.round(r.y + r.height / 2) });
		if (results.length >= 60) break;
	}
	return results;
}`
