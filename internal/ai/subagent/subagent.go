package subagent

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/llm"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// DOMSubAgent - AI-эксперт по анализу DOM
type DOMSubAgent struct {
	provider  llm.Provider
	model     string
	maxTokens int
}

// New создаёт Sub-Agent с LLM провайдером
func New(provider llm.Provider, model string, maxTokens int) *DOMSubAgent {
	return &DOMSubAgent{provider: provider, model: model, maxTokens: maxTokens}
}

// Analyze анализирует страницу и отвечает на вопрос
func (d *DOMSubAgent) Analyze(ctx context.Context, html, liveElements, question string) (string, error) {
	logger.Info(ctx, "🧠 Sub-Agent: Analyzing", zap.String("question", question))
	html = truncate(html, 80000)
	msg := fmt.Sprintf("ЭЛЕМЕНТЫ:\n%s\n\nHTML:\n%s\n\nВОПРОС: %s", liveElements, html, question)

	result, err := d.send(ctx, AnalyzePrompt, msg)
	if err != nil {
		logger.Error(ctx, "❌ Sub-Agent error", zap.Error(err))
		return "", err
	}
	logger.Info(ctx, "🧠 Sub-Agent result", zap.String("analysis", truncateStr(result, 200)))
	return result, nil
}

// AnalyzeError анализирует ошибку и предлагает альтернативу
func (d *DOMSubAgent) AnalyzeError(ctx context.Context, html, liveElements, failedAction, errorMsg string) (string, error) {
	logger.Info(ctx, "🔍 Sub-Agent: Error analysis", zap.String("action", failedAction), zap.String("error", errorMsg))
	html = truncate(html, 60000)
	msg := fmt.Sprintf("ДЕЙСТВИЕ: %s\nОШИБКА: %s\n\nЭЛЕМЕНТЫ:\n%s\n\nHTML:\n%s", failedAction, errorMsg, liveElements, html)

	result, err := d.send(ctx, ErrorAnalysisPrompt, msg)
	if err != nil {
		logger.Error(ctx, "❌ Sub-Agent error", zap.Error(err))
		return "", err
	}
	logger.Info(ctx, "🔍 Sub-Agent diagnosis", zap.String("result", truncateStr(result, 300)))
	return result, nil
}

// Query выполняет запрос к DOM
func (d *DOMSubAgent) Query(ctx context.Context, html, query string) (string, error) {
	return d.Analyze(ctx, html, "", query)
}

// QueryWithScreenshot анализирует скриншот
func (d *DOMSubAgent) QueryWithScreenshot(ctx context.Context, screenshotB64, query string) (string, error) {
	logger.Info(ctx, "🔍 Sub-Agent: Visual query", zap.String("query", query))

	resp, err := d.provider.ChatWithVision(ctx, &llm.VisionRequest{
		Model:       d.model,
		MaxTokens:   d.maxTokens,
		System:      VisualAnalysisPrompt,
		ImageBase64: screenshotB64,
		ImageType:   "image/png",
		Query:       query,
	})
	if err != nil {
		return "", err
	}
	return extractText(resp), nil
}

func (d *DOMSubAgent) send(ctx context.Context, system, user string) (string, error) {
	resp, err := d.provider.Chat(ctx, &llm.ChatRequest{
		Model:     d.model,
		MaxTokens: d.maxTokens,
		System:    system,
		Messages: []llm.Message{{
			Role:    "user",
			Content: []llm.ContentBlock{{Type: "text", Text: user}},
		}},
	})
	if err != nil {
		return "", err
	}
	return extractText(resp), nil
}

func extractText(resp *llm.ChatResponse) string {
	for _, b := range resp.Content {
		if b.IsText() {
			return b.Text
		}
	}
	return ""
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "\n... [обрезано]"
	}
	return s
}

func truncateStr(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
