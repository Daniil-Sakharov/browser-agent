package llm

import "context"

// Provider интерфейс LLM провайдера
type Provider interface {
	// Chat отправляет сообщения и получает ответ
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	// ChatWithVision отправляет сообщение с изображением
	ChatWithVision(ctx context.Context, req *VisionRequest) (*ChatResponse, error)
}

// ChatRequest запрос к LLM
type ChatRequest struct {
	Model       string
	MaxTokens   int
	System      string
	Messages    []Message
	Tools       []Tool
	Temperature float64
}

// VisionRequest запрос с изображением
type VisionRequest struct {
	Model       string
	MaxTokens   int
	System      string
	ImageBase64 string
	ImageType   string // "image/png", "image/jpeg"
	Query       string
}

// ChatResponse ответ от LLM
type ChatResponse struct {
	Content    []ContentBlock
	StopReason string
}

// ContentBlock блок контента в ответе
type ContentBlock struct {
	Type      string // "text" или "tool_use"
	Text      string
	ToolUseID string
	ToolName  string
	ToolInput map[string]interface{}
}

// Message сообщение в диалоге
type Message struct {
	Role    string // "user" или "assistant"
	Content []ContentBlock
}

// Tool определение инструмента
type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// IsToolUse проверяет является ли блок вызовом инструмента
func (c *ContentBlock) IsToolUse() bool {
	return c.Type == "tool_use"
}

// IsText проверяет является ли блок текстом
func (c *ContentBlock) IsText() bool {
	return c.Type == "text"
}
