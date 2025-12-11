package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/ai/tools"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// DecideNextAction отправляет запрос в Claude и получает решение
func (c *Client) DecideNextAction(ctx context.Context) (*domain.Decision, error) {
	logger.Info(ctx, "🤔 Asking Claude for next action")

	systemPrompt := c.buildSystemPrompt()

	response, err := c.anthropic.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: int64(c.maxTokens),
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: c.conversation.GetMessages(),
		Tools:    tools.BrowserTools(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get response from Claude: %w", err)
	}

	logger.Info(ctx, "📩 Received response from Claude",
		zap.String("stop_reason", string(response.StopReason)),
		zap.Int("content_blocks", len(response.Content)))

	// Сохраняем ответ в историю
	c.conversation.AddAssistantMessage(response)

	// Парсим решение
	decision, err := c.parseDecision(ctx, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decision: %w", err)
	}

	return decision, nil
}
