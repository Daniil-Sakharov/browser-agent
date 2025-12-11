package rules

import (
	"strings"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// HighRiskRules возвращает правила высокого риска
func HighRiskRules() []Rule {
	return []Rule{
		{
			Pattern: "email_deletion", Level: RiskLevelHigh, Reason: "attempting to delete email",
			Suggestions: []string{"Email will be moved to trash", "Data may be lost"},
			Matcher: func(a domain.Action, c *domain.PageContext) bool {
				if !IsClickAction(a) || c == nil {
					return false
				}
				url, sel := strings.ToLower(c.URL), strings.ToLower(a.Selector)
				return ContainsAny(url, []string{"mail.yandex", "mail.google", "outlook", "mail.ru"}) &&
					ContainsAny(sel, []string{"delete", "trash", "удалить", "корзин"})
			},
		},
		{
			Pattern: "job_application", Level: RiskLevelHigh, Reason: "sending job application",
			Suggestions: []string{"Application will be sent", "Verify info is correct"},
			Matcher: func(a domain.Action, c *domain.PageContext) bool {
				if a.Type != domain.ActionTypeClick || c == nil {
					return false
				}
				url, text := strings.ToLower(c.URL), GetActionText(a)
				return ContainsAny(url, []string{"hh.ru", "headhunter", "superjob", "rabota"}) &&
					ContainsAny(text, []string{"откликнуться", "отправить", "apply", "submit", "отклик"})
			},
		},
		{
			Pattern: "order_placement", Level: RiskLevelHigh, Reason: "placing order",
			Suggestions: []string{"Order will be placed", "May involve real money"},
			Matcher: func(a domain.Action, c *domain.PageContext) bool {
				if a.Type != domain.ActionTypeClick || c == nil {
					return false
				}
				url, text := strings.ToLower(c.URL), GetActionText(a)
				return ContainsAny(url, []string{"lavka.yandex", "eda.yandex", "checkout", "cart", "ozon", "wildberries"}) &&
					ContainsAny(text, []string{"оформить", "заказать", "оплатить", "купить", "checkout", "pay", "order"})
			},
		},
		{
			Pattern: "data_deletion", Level: RiskLevelHigh, Reason: "deleting data",
			Suggestions: []string{"Data may be permanently lost", "Consider backup"},
			Matcher: func(a domain.Action, c *domain.PageContext) bool {
				if !IsClickAction(a) {
					return false
				}
				return ContainsAny(GetActionText(a), []string{"delete", "remove", "trash", "clear", "удалить", "очистить", "permanently", "навсегда"})
			},
		},
		{
			Pattern: "sensitive_form", Level: RiskLevelHigh, Reason: "submitting sensitive info",
			Suggestions: []string{"Verify info is correct", "Check website is legitimate"},
			Matcher: func(a domain.Action, c *domain.PageContext) bool {
				if a.Type != domain.ActionTypeClick && a.Type != domain.ActionTypeType {
					return false
				}
				return ContainsAny(GetActionText(a), []string{"password", "credit card", "ssn", "bank account"})
			},
		},
	}
}
