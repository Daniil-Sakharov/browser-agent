package agent

import (
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// truncateForProgress обрезает строку для вывода прогресса
func truncateForProgress(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// formatErrorContextMessage форматирует контекст ошибки
func formatErrorContextMessage(result *domain.ActionResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("❌ Действие не выполнено: %s\n", result.Message))
	sb.WriteString(fmt.Sprintf("Селектор: %s\n\n", result.ErrorContext.FailedSelector))

	if len(result.ErrorContext.SimilarElements) > 0 {
		sb.WriteString("📋 Похожие элементы на странице:\n")
		for i, elem := range result.ErrorContext.SimilarElements {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, elem))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("💡 Рекомендация: %s\n", result.ErrorContext.Suggestion))
	sb.WriteString("\nПопробуй использовать один из похожих селекторов или другой подход.")

	return sb.String()
}
