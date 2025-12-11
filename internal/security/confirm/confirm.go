package confirm

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/rules"
)

// Risk информация о риске действия
type Risk struct {
	Level       rules.RiskLevel
	Reason      string
	Suggestions []string
}

type confirmModel struct {
	action    domain.Action
	risk      Risk
	selected  int
	confirmed bool
	cancelled bool
	quitting  bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k", "left", "h":
			m.selected = 0
		case "down", "j", "right", "l":
			m.selected = 1
		case "enter", " ":
			m.confirmed, m.quitting = m.selected == 0, true
			return m, tea.Quit
		case "y", "Y":
			m.confirmed, m.quitting = true, true
			return m, tea.Quit
		case "n", "N", "q", "esc", "ctrl+c":
			m.cancelled, m.quitting = true, true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.quitting {
		if m.confirmed {
			return selectedStyle.Render("✓ Выполняю...\n")
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("✗ Отменено\n")
	}

	var b strings.Builder
	
	// Разные заголовки для разных уровней риска
	title := "⚠️  ОПАСНОЕ ДЕЙСТВИЕ"
	if m.risk.Level == rules.RiskLevelCritical {
		title = "💳 ФИНАНСОВАЯ ОПЕРАЦИЯ - ПОДТВЕРДИТЕ ОПЛАТУ"
	}
	
	criticalTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000")).
		Background(lipgloss.Color("#330000")).
		Padding(0, 1).
		MarginBottom(1)
	
	if m.risk.Level == rules.RiskLevelCritical {
		b.WriteString(criticalTitleStyle.Render(title) + "\n")
	} else {
		b.WriteString(titleStyle.Render(title) + "\n")
	}
	b.WriteString(boxStyle.Render(m.details()) + "\n")

	if len(m.risk.Suggestions) > 0 {
		b.WriteString(warningStyle.Render("⚡ Предупреждения:\n"))
		for _, s := range m.risk.Suggestions {
			b.WriteString(warningStyle.Render("   • "+TranslateSuggestion(s)) + "\n")
		}
	}

	b.WriteString(m.options())
	b.WriteString(hintStyle.Render("[↑↓] выбор  [Enter] ОК  [Y] да  [N/Esc] нет"))
	return b.String()
}

func (m confirmModel) details() string {
	var d strings.Builder
	d.WriteString(labelStyle.Render("Действие: ") + valueStyle.Render(GetActionName(m.action.Type)) + "\n")
	if m.action.Selector != "" {
		sel := m.action.Selector
		if len(sel) > 50 {
			sel = sel[:50] + "..."
		}
		d.WriteString(labelStyle.Render("Селектор: ") + valueStyle.Render(sel) + "\n")
	}
	if m.action.Value != "" {
		val := m.action.Value
		if len(val) > 50 {
			val = val[:50] + "..."
		}
		d.WriteString(labelStyle.Render("Значение: ") + valueStyle.Render(val) + "\n")
	}
	riskStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(GetRiskColor(m.risk.Level))).Bold(true)
	d.WriteString(labelStyle.Render("Риск: ") + riskStyle.Render(GetRiskLevelName(m.risk.Level)) + "\n")
	d.WriteString(labelStyle.Render("Причина: ") + valueStyle.Render(TranslateReason(m.risk.Reason)))
	return d.String()
}

func (m confirmModel) options() string {
	yes, no := "  Выполнить", "  Отменить"
	if m.selected == 0 {
		yes = "▶ Выполнить"
	} else {
		no = "▶ Отменить"
	}
	return unselectedStyle.Render(yes) + "\n" + selectedStyle.Render(no) + "\n"
}

// Action запрашивает подтверждение действия
func Action(action domain.Action, risk Risk) (bool, error) {
	m := confirmModel{action: action, risk: risk, selected: 1}
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("dialog failed: %w", err)
	}
	return final.(confirmModel).confirmed, nil
}
