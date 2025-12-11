package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// Client интеграция с Claude API
type Client struct {
	anthropic    anthropic.Client
	model        string
	maxTokens    int
	temperature  float64
	conversation *Conversation
}

// New создает новый AI клиент
func New(ctx context.Context, apiKey, baseURL, model string, maxTokens int, temperature float64) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}

	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := anthropic.NewClient(opts...)

	logger.Info(ctx, "✅ AI Client initialized",
		zap.String("model", model),
		zap.String("base_url", baseURL),
		zap.Int("max_tokens", maxTokens))

	return &Client{
		anthropic:   client,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}

// NewConversation создает новый диалог
func (c *Client) NewConversation() {
	c.conversation = NewConversation()
}

// AddUserMessage добавляет сообщение пользователя
func (c *Client) AddUserMessage(task string, pageContext *domain.PageContext) error {
	return c.conversation.AddUserMessage(task, pageContext)
}

// AddToolResult добавляет результат выполнения tool
func (c *Client) AddToolResult(toolUseID string, result string, isError bool) {
	c.conversation.AddToolResult(toolUseID, result, isError)
}

// AddToolResultWithImage добавляет результат tool с изображением
func (c *Client) AddToolResultWithImage(toolUseID string, result string, imageB64 string, isError bool) {
	c.conversation.AddToolResultWithImage(toolUseID, result, imageB64, isError)
}

// Close закрывает клиент
func (c *Client) Close(ctx context.Context) error {
	logger.Info(ctx, "🚫 Closing AI Client")
	return nil
}
