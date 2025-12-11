package ai

import "github.com/anthropics/anthropic-sdk-go"

// inputTools - ввод текста и клавиатура
func inputTools() []anthropic.ToolParam {
	return []anthropic.ToolParam{
		{
			Name:        "type_text",
			Description: anthropic.String("Type text into input field"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"selector": map[string]interface{}{"type": "string", "description": "Input selector"},
					"text":     map[string]interface{}{"type": "string", "description": "Text to type"},
				},
				Required: []string{"selector", "text"},
			},
		},
		{
			Name:        "press_enter",
			Description: anthropic.String("Press Enter key"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
				Required:   []string{},
			},
		},
		{
			Name:        "wait",
			Description: anthropic.String("Wait for element to appear"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"selector": map[string]interface{}{"type": "string", "description": "Element selector"},
				},
				Required: []string{"selector"},
			},
		},
		{
			Name:        "complete_task",
			Description: anthropic.String("Mark task as completed"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"result": map[string]interface{}{"type": "string", "description": "Result summary"},
				},
				Required: []string{"result"},
			},
		},
	}
}
