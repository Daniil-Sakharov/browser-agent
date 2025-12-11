package tools

import "github.com/anthropics/anthropic-sdk-go"

// QueryTools - анализ и скриншоты
func QueryTools() []anthropic.ToolParam {
	return []anthropic.ToolParam{
		{
			Name:        "take_screenshot",
			Description: anthropic.String("Take screenshot of current page"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"full_page": map[string]interface{}{"type": "boolean", "description": "Capture full page"},
				},
				Required: []string{},
			},
		},
		{
			Name:        "query_dom",
			Description: anthropic.String("Find clickable elements. Returns text: selectors"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"query": map[string]interface{}{"type": "string", "description": "Optional filter"},
				},
				Required: []string{},
			},
		},
		{
			Name:        "analyze_page",
			Description: anthropic.String("Deep AI analysis of page structure and elements"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"question": map[string]interface{}{"type": "string", "description": "Question about the page"},
				},
				Required: []string{"question"},
			},
		},
	}
}
