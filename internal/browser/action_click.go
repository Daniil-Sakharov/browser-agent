package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

func (c *Controller) Click(ctx context.Context, selector string) error {
	timeout := 10 * time.Second
	
	// text: селекторы - приоритет
	if strings.HasPrefix(selector, "text:") {
		text := strings.TrimPrefix(selector, "text:")
		// Сначала пробуем Rod
		if elem, err := c.findByText(ctx, text, timeout); err == nil {
			if err := c.doClick(ctx, elem, selector); err == nil {
				logger.Info(ctx, "✅ Click via Rod", zap.String("text", text))
				return nil
			}
		}
		// Fallback на JavaScript
		if err := c.jsClickText(ctx, text); err == nil {
			logger.Info(ctx, "✅ Click via JS", zap.String("text", text))
			return nil
		}
		return fmt.Errorf("element not found: %s", selector)
	}
	
	// CSS селекторы
	if elem, err := c.page.Timeout(timeout).Element(selector); err == nil {
		if err := c.doClick(ctx, elem, selector); err == nil {
			logger.Info(ctx, "✅ Click via Rod", zap.String("selector", selector))
			return nil
		}
	}
	// Fallback на JavaScript
	if err := c.jsClick(ctx, selector); err == nil {
		logger.Info(ctx, "✅ Click via JS", zap.String("selector", selector))
		return nil
	}
	return fmt.Errorf("element not found: %s", selector)
}

func (c *Controller) doClick(ctx context.Context, elem *rod.Element, sel string) error {
	elem.ScrollIntoView()
	if elem.Timeout(5*time.Second).WaitVisible() != nil {
		return c.jsFallback(ctx, sel)
	}
	if elem.Click(proto.InputMouseButtonLeft, 1) != nil {
		return c.jsFallback(ctx, sel)
	}
	c.page.Timeout(5*time.Second).WaitStable(300*time.Millisecond)
	return nil
}

func (c *Controller) jsFallback(ctx context.Context, sel string) error {
	if strings.HasPrefix(sel, "text:") {
		return c.jsClickText(ctx, strings.TrimPrefix(sel, "text:"))
	}
	return c.jsClick(ctx, sel)
}

func (c *Controller) jsClick(ctx context.Context, sel string) error {
	js := fmt.Sprintf(`()=>{const e=document.querySelector('%s');if(!e)return{ok:0};e.scrollIntoView();e.click();return{ok:1}}`, strings.ReplaceAll(sel, "'", "\\'"))
	r, _ := c.page.Eval(js)
	if r == nil || !r.Value.Get("ok").Bool() {
		return fmt.Errorf("click failed")
	}
	return nil
}

func (c *Controller) jsClickText(ctx context.Context, text string) error {
	// Более широкий поиск: все кликабельные элементы + любые с текстом
	js := fmt.Sprintf(`()=>{
		const t='%s'.toLowerCase();
		const tags=['button','a','div','span','li','p','label','input','td','th','h1','h2','h3','h4','h5','h6'];
		for(const s of tags){
			for(const e of document.querySelectorAll(s)){
				const x=(e.innerText||e.textContent||'').trim().toLowerCase();
				if(x===t||x.includes(t)||(x.length<100&&t.includes(x)&&x.length>3)){
					e.scrollIntoView({block:'center'});
					e.click();
					return {ok:true,tag:s,text:x.substring(0,50)};
				}
			}
		}
		return {ok:false};
	}`, strings.ReplaceAll(text, "'", "\\'"))
	r, err := c.page.Eval(js)
	if err != nil || r == nil {
		return fmt.Errorf("js eval error")
	}
	if !r.Value.Get("ok").Bool() {
		return fmt.Errorf("text not found")
	}
	logger.Info(ctx, "🎯 JS found element", 
		zap.String("tag", r.Value.Get("tag").String()),
		zap.String("text", r.Value.Get("text").String()))
	return nil
}

func (c *Controller) findByText(ctx context.Context, text string, timeout time.Duration) (*rod.Element, error) {
	js := fmt.Sprintf(`()=>{const t="%s";for(const s of['button','a','div','span','li']){for(const e of document.querySelectorAll(s)){const x=(e.innerText||'').trim();if(x===t||x.includes(t))return e}}return null}`, strings.ReplaceAll(text, `"`, `\"`))
	r, err := c.page.Timeout(timeout).Eval(js)
	if err != nil || r.Value.Nil() {
		return nil, fmt.Errorf("not found")
	}
	id := r.Value.Get("objectId").String()
	if id == "" {
		return nil, fmt.Errorf("no id")
	}
	return c.page.ElementFromObject(&proto.RuntimeRemoteObject{ObjectID: proto.RuntimeRemoteObjectID(id)})
}
