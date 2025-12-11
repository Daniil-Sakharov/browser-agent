package rules

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// RiskLevel уровень риска
type RiskLevel int

const (
	RiskLevelSafe RiskLevel = iota
	RiskLevelMedium
	RiskLevelHigh
	RiskLevelCritical
)

// Rule правило проверки опасности
type Rule struct {
	Pattern     string
	Level       RiskLevel
	Reason      string
	Suggestions []string
	Matcher     func(action domain.Action, ctx *domain.PageContext) bool
}

// Matches проверяет, подходит ли действие под правило
func (r *Rule) Matches(action domain.Action, ctx *domain.PageContext) bool {
	if r.Matcher != nil {
		return r.Matcher(action, ctx)
	}
	return false
}

// BuildRules создает список всех правил
func BuildRules() []Rule {
	var rules []Rule
	rules = append(rules, CriticalRules()...)
	rules = append(rules, HighRiskRules()...)
	rules = append(rules, MediumRiskRules()...)
	return rules
}

// ContainsAny проверяет, содержит ли строка любой из паттернов
func ContainsAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	return false
}

// GetActionText возвращает текст для проверки из действия
func GetActionText(action domain.Action) string {
	return strings.ToLower(action.Selector + " " + action.Value)
}

// IsClickAction проверяет является ли действие кликом
func IsClickAction(action domain.Action) bool {
	return action.Type == domain.ActionTypeClick || action.Type == domain.ActionTypeClickAtPosition
}
