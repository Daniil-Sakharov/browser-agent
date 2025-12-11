package security

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// criticalRules возвращает правила критического уровня
func criticalRules() []Rule {
	return []Rule{
		{
			Pattern: "delete_account",
			Level:   RiskLevelCritical,
			Reason:  "attempting to delete account",
			Suggestions: []string{
				"This action cannot be undone",
				"All data will be permanently lost",
			},
			Matcher: func(action domain.Action, ctx *domain.PageContext) bool {
				if action.Type != domain.ActionTypeClick {
					return false
				}
				text := getActionText(action)
				return containsAny(text, []string{
					"delete account", "close account", "deactivate account",
					"permanently delete", "remove account",
				})
			},
		},
		{
			Pattern: "financial_transaction",
			Level:   RiskLevelCritical,
			Reason:  "attempting financial transaction",
			Suggestions: []string{
				"This may involve real money",
				"Verify payment details carefully",
			},
			Matcher: func(action domain.Action, ctx *domain.PageContext) bool {
				if ctx == nil {
					return false
				}
				url := strings.ToLower(ctx.URL)
				text := getActionText(action)
				return containsAny(url, []string{"payment", "checkout", "billing"}) &&
					containsAny(text, []string{"pay", "purchase", "buy", "confirm payment", "place order"})
			},
		},
	}
}
