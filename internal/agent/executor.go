package agent

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// isNegativeResult проверяет что результат негативный
func isNegativeResult(result string) bool {
	lower := strings.ToLower(result)
	negativeWords := []string{
		"не смог", "не удалось", "извините", "частично",
		"к сожалению", "невозможно", "не получилось", "провал",
		"не выполнен", "не завершен", "ошибка",
	}
	for _, word := range negativeWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func (a *Agent) executeStep(ctx context.Context) (bool, error) {
	a.emitProgress(ProgressEvent{Type: "step", Step: a.stepCount, MaxSteps: a.maxSteps})

	pageCtx, err := a.browser.GetPageContext(ctx)
	if err != nil {
		return false, fmt.Errorf("page: %w", err)
	}

	msg := ""
	if a.stepCount == 1 {
		msg = a.currentTask.Description
	}
	a.ai.AddUserMessage(msg, pageCtx)
	a.emitProgress(ProgressEvent{Type: "waiting"}) // Показываем что ждём ответа

	d, err := a.ai.DecideNextAction(ctx)
	if err != nil {
		return false, fmt.Errorf("AI: %w", err)
	}

	// Показываем рассуждения если есть
	if d.Reasoning != "" {
		a.emitProgress(ProgressEvent{Type: "thinking", Reasoning: d.Reasoning, Tool: string(d.Action.Type)})
	}

	if d.Complete {
		// Проверяем что результат не негативный
		if isNegativeResult(d.Result) && a.stepCount < 20 {
			// Отклоняем негативный complete_task и заставляем продолжить
			logger.Warn(ctx, "⚠️ Rejecting negative complete_task", zap.String("result", d.Result))
			a.ai.AddToolResult(d.ToolUseID, "ОТКЛОНЕНО! Нельзя завершать с негативным результатом. Продолжай работу - попробуй другие способы!", true)
			return false, nil
		}
		a.currentTask.Result = d.Result
		a.emitProgress(ProgressEvent{Type: "result", Result: d.Result, Success: true})
		return true, nil
	}

	if a.security != nil {
		if err := a.security.CheckAction(ctx, d.Action, pageCtx); err != nil {
			return a.handleSecurityError(ctx, d, err)
		}
	}

	a.emitToolProgress(d)
	r, err := a.executeAction(ctx, d.Action)
	if err != nil {
		a.ai.AddToolResult(d.ToolUseID, "Error: "+err.Error(), true)
		return false, err
	}

	a.emitProgress(ProgressEvent{Type: "result", Tool: string(d.Action.Type), Result: r.Message, Success: r.Success})
	a.handleActionResult(ctx, d, r)
	return false, nil
}

func (a *Agent) handleSecurityError(ctx context.Context, d *domain.Decision, err error) (bool, error) {
	if err.Error() == "action rejected by user" {
		a.emitProgress(ProgressEvent{Type: "error", Result: "Отменено", Success: false})
		return true, fmt.Errorf("cancelled")
	}
	a.ai.AddToolResult(d.ToolUseID, "Blocked: "+err.Error(), true)
	return false, nil
}

func (a *Agent) emitToolProgress(d *domain.Decision) {
	p := map[string]string{}
	if d.Action.Selector != "" {
		p["selector"] = d.Action.Selector
	}
	if d.Action.URL != "" {
		p["url"] = d.Action.URL
	}
	a.emitProgress(ProgressEvent{Type: "tool", Tool: string(d.Action.Type), Params: p})
}

func (a *Agent) handleActionResult(ctx context.Context, d *domain.Decision, r *domain.ActionResult) {
	if !r.Success {
		a.handleFailedAction(ctx, d, r)
	} else {
		a.consecutiveFailures, a.lastFailedAction = 0, ""
	}

	// Проверяем что ToolUseID не пустой
	if d.ToolUseID == "" {
		logger.Warn(ctx, "⚠️ Empty ToolUseID, skipping tool result")
		return
	}

	switch {
	case !r.Success && r.ErrorContext != nil:
		a.ai.AddToolResult(d.ToolUseID, formatErrorContextMessage(r), true)
	case r.ScreenshotB64 != "":
		a.ai.AddToolResultWithImage(d.ToolUseID, r.Message, r.ScreenshotB64, !r.Success)
	default:
		a.ai.AddToolResult(d.ToolUseID, r.Message, !r.Success)
	}
}

func (a *Agent) handleFailedAction(ctx context.Context, d *domain.Decision, r *domain.ActionResult) {
	a.consecutiveFailures++
	a.lastFailedAction = fmt.Sprintf("%s sel=%s", d.Action.Type, d.Action.Selector)

	logger.Info(ctx, "❌ Action failed",
		zap.String("action", a.lastFailedAction),
		zap.String("error", r.Message),
		zap.Int("consecutive_failures", a.consecutiveFailures))

	if a.consecutiveFailures >= 2 && a.domSubAgent != nil {
		logger.Info(ctx, "🧠 Calling Sub-Agent for analysis", zap.Int("failures", a.consecutiveFailures))

		// Показываем что Sub-Agent анализирует
		a.emitProgress(ProgressEvent{
			Type:   "subagent_thinking",
			Tool:   "error_analysis",
			Result: fmt.Sprintf("Анализирую ошибку: %s", a.lastFailedAction),
		})

		html, _ := a.browser.GetHTML(ctx)
		live, _ := a.browser.FindElementsLive(ctx, "")

		if analysis, err := a.domSubAgent.AnalyzeError(ctx, html, live, a.lastFailedAction, r.Message); err == nil && analysis != "" {
			// Показываем результат анализа Sub-Agent
			a.emitProgress(ProgressEvent{
				Type:    "subagent_result",
				Tool:    "error_analysis",
				Result:  analysis,
				Success: true,
			})
			r.Message = fmt.Sprintf("%s\n\n🧠 АНАЛИЗ:\n%s", r.Message, analysis)
		}
	}
}
