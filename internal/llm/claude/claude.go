package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/llm"
)

// Provider реализация LLM для Claude
type Provider struct {
	client *anthropic.Client
}

// Config конфигурация Claude
type Config struct {
	APIKey  string
	BaseURL string
}

// New создаёт нового Claude провайдера
func New(cfg Config) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	opts := []option.RequestOption{option.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}

	client := anthropic.NewClient(opts...)
	return &Provider{client: &client}, nil
}

// Chat отправляет сообщения и получает ответ
func (p *Provider) Chat(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(req.Model),
		MaxTokens: int64(req.MaxTokens),
		Messages:  convertMessages(req.Messages),
	}

	if req.System != "" {
		params.System = []anthropic.TextBlockParam{{Text: req.System}}
	}

	if len(req.Tools) > 0 {
		params.Tools = convertTools(req.Tools)
	}

	resp, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("claude chat: %w", err)
	}

	return convertResponse(resp), nil
}

// ChatWithVision отправляет сообщение с изображением
func (p *Provider) ChatWithVision(ctx context.Context, req *llm.VisionRequest) (*llm.ChatResponse, error) {
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(req.Model),
		MaxTokens: int64(req.MaxTokens),
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				anthropic.NewImageBlockBase64(req.ImageType, req.ImageBase64),
				anthropic.NewTextBlock(req.Query),
			},
		}},
	}

	if req.System != "" {
		params.System = []anthropic.TextBlockParam{{Text: req.System}}
	}

	resp, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("claude vision: %w", err)
	}

	return convertResponse(resp), nil
}

func convertMessages(msgs []llm.Message) []anthropic.MessageParam {
	result := make([]anthropic.MessageParam, 0, len(msgs))
	for _, m := range msgs {
		var content []anthropic.ContentBlockParamUnion
		for _, c := range m.Content {
			if c.IsText() {
				content = append(content, anthropic.NewTextBlock(c.Text))
			}
		}
		if len(content) > 0 {
			result = append(result, anthropic.MessageParam{
				Role:    anthropic.MessageParamRole(m.Role),
				Content: content,
			})
		}
	}
	return result
}

func convertTools(tools []llm.Tool) []anthropic.ToolUnionParam {
	result := make([]anthropic.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		tool := anthropic.ToolParam{
			Name:        t.Name,
			Description: anthropic.String(t.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: t.InputSchema,
			},
		}
		result = append(result, anthropic.ToolUnionParam{OfTool: &tool})
	}
	return result
}

func convertResponse(resp *anthropic.Message) *llm.ChatResponse {
	result := &llm.ChatResponse{
		StopReason: string(resp.StopReason),
		Content:    make([]llm.ContentBlock, 0, len(resp.Content)),
	}

	for _, block := range resp.Content {
		cb := llm.ContentBlock{Type: string(block.Type)}
		switch block.Type {
		case "text":
			cb.Text = block.Text
		case "tool_use":
			cb.ToolUseID = block.ID
			cb.ToolName = block.Name
			// Input это json.RawMessage, парсим в map
			if len(block.Input) > 0 {
				var input map[string]interface{}
				if err := json.Unmarshal(block.Input, &input); err == nil {
					cb.ToolInput = input
				}
			}
		}
		result.Content = append(result.Content, cb)
	}

	return result
}
