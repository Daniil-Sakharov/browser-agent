package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

type Conversation struct {
	messages []anthropic.MessageParam
}

func NewConversation() *Conversation {
	return &Conversation{messages: make([]anthropic.MessageParam, 0)}
}

func (c *Conversation) AddUserMessage(task string, ctx *domain.PageContext) error {
	text := fmt.Sprintf("Task: %s\n\nCurrent page context:\n%s", task, c.formatContext(ctx))
	c.messages = append(c.messages, anthropic.NewUserMessage(anthropic.NewTextBlock(text)))
	return nil
}

func (c *Conversation) AddAssistantMessage(msg *anthropic.Message) {
	c.messages = append(c.messages, msg.ToParam())
}

func (c *Conversation) AddToolResult(toolID, result string, isError bool) {
	c.messages = append(c.messages, anthropic.NewUserMessage(anthropic.NewToolResultBlock(toolID, result, isError)))
}

func (c *Conversation) AddToolResultWithImage(toolID, result, imgB64 string, isError bool) {
	if imgB64 == "" {
		c.AddToolResult(toolID, result, isError)
		return
	}

	img := anthropic.ImageBlockParam{Type: "image", Source: anthropic.ImageBlockParamSourceUnion{
		OfBase64: &anthropic.Base64ImageSourceParam{Type: "base64", MediaType: "image/png", Data: imgB64},
	}}
	txt := anthropic.TextBlockParam{Type: "text", Text: result}

	block := anthropic.ToolResultBlockParam{
		ToolUseID: toolID, IsError: anthropic.Bool(isError),
		Content: []anthropic.ToolResultBlockParamContentUnion{{OfImage: &img}, {OfText: &txt}},
	}

	c.messages = append(c.messages, anthropic.MessageParam{
		Role:    anthropic.MessageParamRoleUser,
		Content: []anthropic.ContentBlockParamUnion{{OfToolResult: &block}},
	})
}

func (c *Conversation) GetMessages() []anthropic.MessageParam { return c.messages }
func (c *Conversation) Clear()                                { c.messages = make([]anthropic.MessageParam, 0) }

func (c *Conversation) formatContext(pctx *domain.PageContext) string {
	if pctx == nil {
		return "No page context"
	}

	out := fmt.Sprintf("URL: %s\nTitle: %s\n\n", pctx.URL, pctx.Title)

	// Логируем детали контекста
	ctx := context.Background()
	logger.Info(ctx, "📄 Page context",
		zap.String("url", pctx.URL),
		zap.Int("elements", len(pctx.InteractiveElems)),
		zap.Int("text_len", len(pctx.VisibleText)))

	if len(pctx.InteractiveElems) > 0 {
		out += "Interactive elements:\n"
		for _, e := range pctx.InteractiveElems {
			out += fmt.Sprintf("- %s [%s]: \"%s\" (%s)\n", e.Type, e.Tag, truncateText(e.Text, 50), e.Selector)
		}
	}

	if pctx.VisibleText != "" {
		text := pctx.VisibleText
		if len(text) > 2000 {
			text = text[:2000] + "..."
		}
		out += fmt.Sprintf("Visible: %s\n", text)
	}
	return out
}

func truncateText(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func ParseToolInput(input, target interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
