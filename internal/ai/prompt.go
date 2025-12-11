package ai

import "github.com/Daniil-Sakharov/BrowserAgent/internal/ai/prompts"

// buildSystemPrompt возвращает системный промпт для Claude
func (c *Client) buildSystemPrompt() string {
	return prompts.System
}
