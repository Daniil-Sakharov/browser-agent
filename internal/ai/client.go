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

// Config интерфейс конфигурации AI
type Config interface {
	APIKey() string
	BaseURL() string
	Model() string
	MaxTokens() int
	Temperature() float64
}

// New создает новый AI клиент
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.APIKey() == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}

	opts := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey()),
	}

	// Если указан кастомный URL - используем его (прокси/альтернативный API)
	if cfg.BaseURL() != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL()))
	}

	client := anthropic.NewClient(opts...)

	logger.Info(ctx, "✅ AI Client initialized",
		zap.String("model", cfg.Model()),
		zap.String("base_url", cfg.BaseURL()),
		zap.Int("max_tokens", cfg.MaxTokens()))

	return &Client{
		anthropic:   client,
		model:       cfg.Model(),
		maxTokens:   cfg.MaxTokens(),
		temperature: cfg.Temperature(),
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
