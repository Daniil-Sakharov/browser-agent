package agent

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// executeAction выполняет конкретное действие
func (a *Agent) executeAction(ctx context.Context, action domain.Action) (*domain.ActionResult, error) {
	logger.Info(ctx, "⚡ Executing action",
		zap.String("type", string(action.Type)),
		zap.String("selector", action.Selector),
		zap.String("url", action.URL))

	switch action.Type {
	case domain.ActionTypeQueryDOM:
		return a.executeQueryDOM(ctx, action)
	case domain.ActionTypeAnalyzePage:
		return a.executeAnalyzePage(ctx, action)
	default:
		result, err := a.browser.ExecuteAction(ctx, action)
		if err != nil {
			return nil, err
		}
		if a.interactive {
			logger.Info(ctx, "⏸️  Interactive mode - press Enter to continue")
		}
		return result, nil
	}
}

// executeAnalyzePage выполняет глубокий анализ страницы через Sub-Agent
func (a *Agent) executeAnalyzePage(ctx context.Context, action domain.Action) (*domain.ActionResult, error) {
	logger.Info(ctx, "🧠 Analyze Page", zap.String("question", action.Question))

	a.emitProgress(ProgressEvent{
		Type: "subagent", Tool: "analyze_page",
		Result: "🧠 Глубокий анализ страницы...",
	})

	html, _ := a.browser.GetHTML(ctx)
	liveElements, _ := a.browser.FindElementsLive(ctx, "")

	analysis, err := a.domSubAgent.Analyze(ctx, html, liveElements, action.Question)
	if err != nil {
		return &domain.ActionResult{
			Success: false, Action: string(action.Type),
			Message: fmt.Sprintf("Ошибка анализа: %s", err.Error()),
		}, nil
	}

	a.emitProgress(ProgressEvent{
		Type: "subagent_result", Tool: "analyze_page",
		Result: truncateForProgress(analysis, 150), Success: true,
	})

	return &domain.ActionResult{
		Success: true, Action: string(action.Type),
		Message: analysis, QueryResult: analysis,
	}, nil
}

// executeQueryDOM выполняет query_dom - FindElementsLive как главный источник
func (a *Agent) executeQueryDOM(ctx context.Context, action domain.Action) (*domain.ActionResult, error) {
	logger.Info(ctx, "🔍 DOM Query", zap.String("query", action.Query))

	a.emitProgress(ProgressEvent{
		Type: "subagent", Tool: "dom_search",
		Result: "Поиск элементов в реальном DOM...",
	})

	liveElements, err := a.browser.FindElementsLive(ctx, action.Query)
	if err != nil {
		liveElements = "Ошибка при поиске элементов."
	}

	a.emitProgress(ProgressEvent{
		Type: "subagent_result", Tool: "dom_search",
		Result: truncateForProgress(liveElements, 100), Success: true,
	})

	var response strings.Builder
	response.WriteString("🎯 ИСПОЛЬЗУЙ СЕЛЕКТОРЫ ТОЛЬКО ИЗ ЭТОГО СПИСКА:\n\n")
	response.WriteString(liveElements)
	response.WriteString("\n\n💡 Если нужный элемент не найден - используй click_at_position с координатами.")

	result := response.String()
	return &domain.ActionResult{
		Success: true, Action: string(action.Type),
		Message: result, QueryResult: result,
	}, nil
}
