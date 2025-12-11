package rules

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// MediumRiskRules возвращает правила среднего риска
func MediumRiskRules() []Rule {
	return []Rule{
		{
			Pattern: "email_send",
			Level:   RiskLevelMedium,
			Reason:  "attempting to send email",
			Suggestions: []string{
				"Email will be sent to recipients",
			},
			Matcher: func(action domain.Action, ctx *domain.PageContext) bool {
				if action.Type != domain.ActionTypeClick {
					return false
				}
				text := GetActionText(action)
				return ContainsAny(text, []string{"send", "send email", "отправить"}) &&
					ctx != nil && ContainsAny(strings.ToLower(ctx.URL), []string{"mail", "почта"})
			},
		},
		{
			Pattern: "settings_change",
			Level:   RiskLevelMedium,
			Reason:  "attempting to change settings",
			Matcher: func(action domain.Action, ctx *domain.PageContext) bool {
				if ctx == nil {
					return false
				}
				url := strings.ToLower(ctx.URL)
				return ContainsAny(url, []string{"settings", "preferences", "config", "настройки"})
			},
		},
	}
}
