package confirm

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/rules"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")).
			Italic(true).
			MarginTop(1)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			MarginTop(1)
)

// GetRiskColor возвращает цвет для уровня риска
func GetRiskColor(level rules.RiskLevel) string {
	switch level {
	case rules.RiskLevelCritical:
		return "#FF0000"
	case rules.RiskLevelHigh:
		return "#FF6B6B"
	case rules.RiskLevelMedium:
		return "#FFA500"
	default:
		return "#00FF00"
	}
}

// GetRiskLevelName возвращает название уровня риска
func GetRiskLevelName(level rules.RiskLevel) string {
	switch level {
	case rules.RiskLevelCritical:
		return "КРИТИЧЕСКИЙ"
	case rules.RiskLevelHigh:
		return "ВЫСОКИЙ"
	case rules.RiskLevelMedium:
		return "СРЕДНИЙ"
	default:
		return "БЕЗОПАСНЫЙ"
	}
}
