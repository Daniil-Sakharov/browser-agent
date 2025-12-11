package security

import "github.com/charmbracelet/lipgloss"

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

func getRiskColor(level RiskLevel) string {
	switch level {
	case RiskLevelCritical:
		return "#FF0000"
	case RiskLevelHigh:
		return "#FF6B6B"
	case RiskLevelMedium:
		return "#FFA500"
	case RiskLevelLow:
		return "#FFFF00"
	default:
		return "#00FF00"
	}
}

func getRiskLevelName(level RiskLevel) string {
	switch level {
	case RiskLevelCritical:
		return "КРИТИЧЕСКИЙ"
	case RiskLevelHigh:
		return "ВЫСОКИЙ"
	case RiskLevelMedium:
		return "СРЕДНИЙ"
	case RiskLevelLow:
		return "НИЗКИЙ"
	default:
		return "БЕЗОПАСНЫЙ"
	}
}
