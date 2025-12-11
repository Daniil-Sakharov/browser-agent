package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

type DOMSubAgent struct {
	client    *anthropic.Client
	model     string
	maxTokens int
}

func NewDOMSubAgent(apiKey, baseURL, model string, maxTokens int) *DOMSubAgent {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	client := anthropic.NewClient(opts...)
	return &DOMSubAgent{client: &client, model: model, maxTokens: maxTokens}
}

func (d *DOMSubAgent) Analyze(ctx context.Context, html, liveElements, question string) (string, error) {
	logger.Info(ctx, "🧠 Sub-Agent: Analyzing", zap.String("question", question))
	html = truncate(html, 80000)
	msg := fmt.Sprintf("ЭЛЕМЕНТЫ:\n%s\n\nHTML:\n%s\n\nВОПРОС: %s", liveElements, html, question)
	result, err := d.send(ctx, analyzeSystemPrompt, msg)
	if err != nil {
		logger.Error(ctx, "❌ Sub-Agent error", zap.Error(err))
		return "", err
	}
	logger.Info(ctx, "🧠 Sub-Agent result", zap.String("analysis", truncateString(result, 200)))
	return result, nil
}

func (d *DOMSubAgent) AnalyzeError(ctx context.Context, html, liveElements, failedAction, errorMsg string) (string, error) {
	logger.Info(ctx, "🔍 Sub-Agent: Error analysis", zap.String("action", failedAction), zap.String("error", errorMsg))
	html = truncate(html, 60000)
	msg := fmt.Sprintf("ДЕЙСТВИЕ: %s\nОШИБКА: %s\n\nЭЛЕМЕНТЫ:\n%s\n\nHTML:\n%s", failedAction, errorMsg, liveElements, html)
	result, err := d.send(ctx, errorAnalysisSystemPrompt, msg)
	if err != nil {
		logger.Error(ctx, "❌ Sub-Agent error", zap.Error(err))
		return "", err
	}
	logger.Info(ctx, "🔍 Sub-Agent diagnosis", zap.String("result", truncateString(result, 300)))
	return result, nil
}

func (d *DOMSubAgent) Query(ctx context.Context, html, query string) (string, error) {
	return d.Analyze(ctx, html, "", query)
}

func (d *DOMSubAgent) QueryWithScreenshot(ctx context.Context, screenshotB64, query string) (string, error) {
	logger.Info(ctx, "🔍 Sub-Agent: Visual query", zap.String("query", query))
	msg, err := d.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.Model(d.model), MaxTokens: int64(d.maxTokens),
		System: []anthropic.TextBlockParam{{Text: visualAnalysisSystemPrompt}},
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				anthropic.NewImageBlockBase64("image/png", screenshotB64),
				anthropic.NewTextBlock(query),
			},
		}},
	})
	if err != nil { return "", err }
	return extractText(msg), nil
}

func (d *DOMSubAgent) send(ctx context.Context, system, user string) (string, error) {
	msg, err := d.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.Model(d.model), MaxTokens: int64(d.maxTokens),
		System: []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(user)},
		}},
	})
	if err != nil { return "", err }
	return extractText(msg), nil
}

func extractText(m *anthropic.Message) string {
	for _, b := range m.Content { if b.Type == "text" { return b.Text } }
	return ""
}

func truncate(s string, max int) string {
	if len(s) > max { return s[:max] + "\n... [обрезано]" }
	return s
}

func truncateString(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max { return s[:max] + "..." }
	return s
}
