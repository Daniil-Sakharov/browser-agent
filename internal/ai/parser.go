package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

func (c *Client) parseDecision(ctx context.Context, resp *anthropic.Message) (*domain.Decision, error) {
	d := &domain.Decision{Complete: resp.StopReason == anthropic.StopReasonEndTurn}

	for _, block := range resp.Content {
		switch content := block.AsAny().(type) {
		case anthropic.TextBlock:
			d.Reasoning = content.Text
			logger.Debug(ctx, "💭 Reasoning", zap.String("text", content.Text))
		case anthropic.ToolUseBlock:
			logger.Info(ctx, "🔧 Tool", zap.String("tool", content.Name))
			d.ToolUseID = content.ID
			action, err := parseToolUse(content)
			if err != nil { return nil, err }
			d.Action = *action
			if content.Name == "complete_task" {
				d.Complete = true
				var in struct{ Result string `json:"result"` }
				json.Unmarshal([]byte(content.JSON.Input.Raw()), &in)
				d.Result = in.Result
			}
		}
	}
	return d, nil
}

func parseToolUse(t anthropic.ToolUseBlock) (*domain.Action, error) {
	raw, a := []byte(t.JSON.Input.Raw()), &domain.Action{}
	switch t.Name {
	case "navigate":
		var in struct{ URL string `json:"url"` }; json.Unmarshal(raw, &in); a.Type, a.URL = domain.ActionTypeNavigate, in.URL
	case "click":
		var in struct{ Selector string `json:"selector"` }; json.Unmarshal(raw, &in); a.Type, a.Selector = domain.ActionTypeClick, in.Selector
	case "click_at_position":
		var in struct{ X, Y int }; json.Unmarshal(raw, &in); a.Type, a.X, a.Y = domain.ActionTypeClickAtPosition, in.X, in.Y
	case "type_text":
		var in struct{ Selector, Text string }; json.Unmarshal(raw, &in); a.Type, a.Selector, a.Value = domain.ActionTypeType, in.Selector, in.Text
	case "scroll":
		var in struct{ Direction string }; json.Unmarshal(raw, &in); a.Type, a.Direction = domain.ActionTypeScroll, in.Direction
	case "wait":
		var in struct{ Selector string }; json.Unmarshal(raw, &in); a.Type, a.Selector = domain.ActionTypeWait, in.Selector
	case "press_enter":
		a.Type = domain.ActionTypePressEnter
	case "complete_task":
		var in struct{ Result string }; json.Unmarshal(raw, &in); a.Type, a.Value = domain.ActionTypeCompleteTask, in.Result
	case "take_screenshot":
		var in struct{ FullPage bool }; json.Unmarshal(raw, &in); a.Type, a.FullPage = domain.ActionTypeTakeScreenshot, in.FullPage
	case "query_dom":
		var in struct{ Query string }; json.Unmarshal(raw, &in); a.Type, a.Query = domain.ActionTypeQueryDOM, in.Query
	case "analyze_page":
		var in struct{ Question string }; json.Unmarshal(raw, &in); a.Type, a.Question = domain.ActionTypeAnalyzePage, in.Question
	default:
		return nil, fmt.Errorf("unknown: %s", t.Name)
	}
	return a, nil
}
