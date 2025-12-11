package security

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// Rule правило проверки опасности
type Rule struct {
	Pattern     string
	Level       RiskLevel
	Reason      string
	Suggestions []string
	Matcher     func(action domain.Action, pageContext *domain.PageContext) bool
}

// Matches проверяет, подходит ли действие под правило
func (r *Rule) Matches(action domain.Action, pageContext *domain.PageContext) bool {
	if r.Matcher != nil {
		return r.Matcher(action, pageContext)
	}
	return false
}

// buildDangerousRules создает список правил опасных действий
func buildDangerousRules() []Rule {
	var rules []Rule
	rules = append(rules, criticalRules()...)
	rules = append(rules, highRiskRules()...)
	rules = append(rules, mediumRiskRules()...)
	return rules
}

// containsAny проверяет, содержит ли строка любой из паттернов
func containsAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	return false
}

// getActionText возвращает текст для проверки из действия
func getActionText(action domain.Action) string {
	return strings.ToLower(action.Selector + " " + action.Value)
}

// isClickAction проверяет является ли действие кликом
func isClickAction(action domain.Action) bool {
	return action.Type == domain.ActionTypeClick || action.Type == domain.ActionTypeClickAtPosition
}
