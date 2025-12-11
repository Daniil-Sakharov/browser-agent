package tools

import "github.com/anthropics/anthropic-sdk-go"

// TabTools - инструменты для работы с вкладками браузера
func TabTools() []anthropic.ToolParam {
	return []anthropic.ToolParam{
		{
			Name:        "list_tabs",
			Description: anthropic.String("List all open browser tabs. Use when click succeeded but page didn't change - new tab might have opened!"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
				Required:   []string{},
			},
		},
		{
			Name:        "switch_tab",
			Description: anthropic.String("Switch to a specific browser tab by index (1-based). Use after list_tabs to switch to new tab"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"tab_index": map[string]interface{}{
						"type":        "integer",
						"description": "Tab index (1 = first tab, 2 = second tab, etc.)",
					},
				},
				Required: []string{"tab_index"},
			},
		},
		{
			Name:        "close_tab",
			Description: anthropic.String("Close current tab and switch to previous one"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
				Required:   []string{},
			},
		},
	}
}
