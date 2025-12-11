package action

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

// Click выполняет клик по селектору (text: или CSS)
func Click(ctx context.Context, p PageProvider, selector string) error {
	page := p.GetPage()
	timeout := 10 * time.Second

	if strings.HasPrefix(selector, "text:") {
		return clickByText(ctx, p, page, strings.TrimPrefix(selector, "text:"), timeout)
	}
	return clickByCSS(ctx, p, page, selector, timeout)
}

func clickByText(ctx context.Context, p PageProvider, page *rod.Page, text string, timeout time.Duration) error {
	if elem, err := findByText(page, text, timeout); err == nil {
		if err := doClick(ctx, p, elem); err == nil {
			logger.Info(ctx, "✅ Click via Rod", zap.String("text", text))
			return nil
		}
	}
	if err := jsClickText(ctx, page, text); err == nil {
		logger.Info(ctx, "✅ Click via JS", zap.String("text", text))
		return nil
	}
	return fmt.Errorf("element not found: text:%s", text)
}

func clickByCSS(ctx context.Context, p PageProvider, page *rod.Page, selector string, timeout time.Duration) error {
	if elem, err := page.Timeout(timeout).Element(selector); err == nil {
		if err := doClick(ctx, p, elem); err == nil {
			logger.Info(ctx, "✅ Click via Rod", zap.String("selector", selector))
			return nil
		}
	}
	if err := jsClick(page, selector); err == nil {
		logger.Info(ctx, "✅ Click via JS", zap.String("selector", selector))
		return nil
	}
	return fmt.Errorf("element not found: %s", selector)
}

func doClick(ctx context.Context, p PageProvider, elem *rod.Element) error {
	elem.ScrollIntoView()
	if elem.Timeout(5*time.Second).WaitVisible() != nil {
		return fmt.Errorf("not visible")
	}
	if elem.Click(proto.InputMouseButtonLeft, 1) != nil {
		return fmt.Errorf("click failed")
	}
	p.WaitStable(5 * time.Second)
	return nil
}

func jsClick(page *rod.Page, sel string) error {
	js := fmt.Sprintf(`()=>{const e=document.querySelector('%s');if(!e)return{ok:0};e.scrollIntoView();e.click();return{ok:1}}`,
		strings.ReplaceAll(sel, "'", "\\'"))
	r, _ := page.Eval(js)
	if r == nil || !r.Value.Get("ok").Bool() {
		return fmt.Errorf("click failed")
	}
	return nil
}

func jsClickText(ctx context.Context, page *rod.Page, text string) error {
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
	r, err := page.Eval(js)
	if err != nil || r == nil || !r.Value.Get("ok").Bool() {
		return fmt.Errorf("text not found")
	}
	logger.Info(ctx, "🎯 JS found", zap.String("tag", r.Value.Get("tag").String()))
	return nil
}

func findByText(page *rod.Page, text string, timeout time.Duration) (*rod.Element, error) {
	js := fmt.Sprintf(`()=>{const t="%s";for(const s of['button','a','div','span','li']){for(const e of document.querySelectorAll(s)){const x=(e.innerText||'').trim();if(x===t||x.includes(t))return e}}return null}`,
		strings.ReplaceAll(text, `"`, `\"`))
	r, err := page.Timeout(timeout).Eval(js)
	if err != nil || r.Value.Nil() {
		return nil, fmt.Errorf("not found")
	}
	id := r.Value.Get("objectId").String()
	if id == "" {
		return nil, fmt.Errorf("no id")
	}
	return page.ElementFromObject(&proto.RuntimeRemoteObject{ObjectID: proto.RuntimeRemoteObjectID(id)})
}
