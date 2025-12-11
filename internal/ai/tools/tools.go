package tools

import "github.com/anthropics/anthropic-sdk-go"

// BrowserTools возвращает все инструменты для Claude
func BrowserTools() []anthropic.ToolUnionParam {
	var params []anthropic.ToolParam
	params = append(params, NavigationTools()...)
	params = append(params, InputTools()...)
	params = append(params, QueryTools()...)
	params = append(params, TabTools()...)

	tools := make([]anthropic.ToolUnionParam, len(params))
	for i, p := range params {
		tools[i] = anthropic.ToolUnionParam{OfTool: &p}
	}
	return tools
}
