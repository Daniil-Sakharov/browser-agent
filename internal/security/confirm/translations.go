package confirm

import "github.com/Daniil-Sakharov/BrowserAgent/internal/domain"

// GetActionName возвращает название действия
func GetActionName(actionType domain.ActionType) string {
	names := map[domain.ActionType]string{
		domain.ActionTypeClick:           "Клик",
		domain.ActionTypeType:            "Ввод текста",
		domain.ActionTypeNavigate:        "Переход",
		domain.ActionTypeScroll:          "Прокрутка",
		domain.ActionTypeWait:            "Ожидание",
		domain.ActionTypePressEnter:      "Нажатие Enter",
		domain.ActionTypeCompleteTask:    "Завершение",
		domain.ActionTypeTakeScreenshot:  "Скриншот",
		domain.ActionTypeQueryDOM:        "Запрос DOM",
		domain.ActionTypeClickAtPosition: "Клик по координатам",
	}
	if name, ok := names[actionType]; ok {
		return name
	}
	return string(actionType)
}

// TranslateReason переводит причину на русский
func TranslateReason(reason string) string {
	translations := map[string]string{
		"attempting to delete data":        "Попытка удаления данных",
		"attempting to delete account":     "Попытка удаления аккаунта",
		"attempting financial transaction": "Финансовая транзакция",
		"submitting sensitive info":        "Отправка конфиденциальной информации",
		"attempting to send email":         "Отправка письма",
		"attempting to change settings":    "Изменение настроек",
		"sending job application":          "Отправка отклика на вакансию",
		"placing order":                    "Оформление заказа",
		"attempting to delete email":       "Удаление письма",
		"deleting data":                    "Удаление данных",
	}
	if translated, ok := translations[reason]; ok {
		return translated
	}
	return reason
}

// TranslateSuggestion переводит предупреждение на русский
func TranslateSuggestion(suggestion string) string {
	translations := map[string]string{
		"Data may be permanently lost":      "Данные могут быть безвозвратно потеряны",
		"Consider backup":                   "Рекомендуется сначала сделать резервную копию",
		"This action cannot be undone":      "Это действие нельзя отменить",
		"All data will be permanently lost": "Все данные будут безвозвратно потеряны",
		"This may involve real money":       "Это может затронуть реальные деньги",
		"Verify payment details carefully":  "Внимательно проверьте данные платежа",
		"Verify info is correct":            "Убедитесь, что информация корректна",
		"Check website is legitimate":       "Убедитесь, что вы на правильном сайте",
		"Email will be sent to recipients":  "Письмо будет отправлено получателям",
		"Application will be sent":          "Отклик будет отправлен работодателю",
		"Order will be placed":              "Заказ будет оформлен",
		"Email will be moved to trash":      "Письмо будет перемещено в корзину",
		"May involve real money":            "Может затронуть реальные деньги",
		"Data may be lost":                  "Данные могут быть потеряны",
	}
	if translated, ok := translations[suggestion]; ok {
		return translated
	}
	return suggestion
}
