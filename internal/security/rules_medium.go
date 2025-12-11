package security

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// mediumRiskRules возвращает правила среднего риска
func mediumRiskRules() []Rule {
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
				text := getActionText(action)
				return containsAny(text, []string{"send", "send email", "отправить"}) &&
					ctx != nil && containsAny(strings.ToLower(ctx.URL), []string{"mail", "почта"})
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
				return containsAny(url, []string{"settings", "preferences", "config", "настройки"})
			},
		},
	}
}
