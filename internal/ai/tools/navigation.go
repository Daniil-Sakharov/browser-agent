package tools

import "github.com/anthropics/anthropic-sdk-go"

// NavigationTools - навигация и клики
func NavigationTools() []anthropic.ToolParam {
	return []anthropic.ToolParam{
		{
			Name:        "navigate",
			Description: anthropic.String("Navigate to a URL"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"url": map[string]interface{}{"type": "string", "description": "URL to navigate to"},
				},
				Required: []string{"url"},
			},
		},
		{
			Name:        "click",
			Description: anthropic.String("Click element by selector. Supports text:ButtonText syntax"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"selector": map[string]interface{}{"type": "string", "description": "CSS selector or text:Text"},
				},
				Required: []string{"selector"},
			},
		},
		{
			Name:        "click_at_position",
			Description: anthropic.String("Click at coordinates. Use when selectors fail"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"x": map[string]interface{}{"type": "integer", "description": "X coordinate"},
					"y": map[string]interface{}{"type": "integer", "description": "Y coordinate"},
				},
				Required: []string{"x", "y"},
			},
		},
		{
			Name:        "scroll",
			Description: anthropic.String("Scroll the page"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"direction": map[string]interface{}{"type": "string", "enum": []string{"up", "down"}},
				},
				Required: []string{"direction"},
			},
		},
	}
}
