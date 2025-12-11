package rules

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// CriticalRules возвращает правила критического уровня
func CriticalRules() []Rule {
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
				text := GetActionText(action)
				return ContainsAny(text, []string{
					"delete account", "close account", "deactivate account",
					"permanently delete", "remove account",
				})
			},
		},
		{
			Pattern: "financial_transaction",
			Level:   RiskLevelCritical,
			Reason:  "финансовая операция - требуется подтверждение",
			Suggestions: []string{
				"Это может включать реальные деньги",
				"Проверьте детали платежа внимательно",
			},
			Matcher: func(action domain.Action, ctx *domain.PageContext) bool {
				if action.Type != domain.ActionTypeClick {
					return false
				}
				text := GetActionText(action)
				
				// Ловим по тексту кнопки (независимо от URL)
				paymentWords := []string{
					"pay", "оплатить", "оплата", "купить", "purchase", "buy",
					"confirm payment", "place order", "оформить заказ",
					"заказать", "checkout", "подтвердить оплату",
					"proceed to payment", "complete order", "завершить заказ",
				}
				if ContainsAny(text, paymentWords) {
					return true
				}
				
				// Ловим по URL если есть контекст
				if ctx != nil {
					url := strings.ToLower(ctx.URL)
					if ContainsAny(url, []string{"payment", "checkout", "billing", "pay", "order"}) &&
						ContainsAny(text, []string{"confirm", "submit", "next", "continue", "далее", "подтвердить"}) {
						return true
					}
				}
				
				return false
			},
		},
	}
}
