package security

import "github.com/Daniil-Sakharov/BrowserAgent/internal/domain"

func getActionName(actionType domain.ActionType) string {
	names := map[domain.ActionType]string{
		domain.ActionTypeClick:          "Клик",
		domain.ActionTypeType:           "Ввод текста",
		domain.ActionTypeNavigate:       "Переход",
		domain.ActionTypeScroll:         "Прокрутка",
		domain.ActionTypeWait:           "Ожидание",
		domain.ActionTypePressEnter:     "Нажатие Enter",
		domain.ActionTypeCompleteTask:   "Завершение",
		domain.ActionTypeTakeScreenshot: "Скриншот",
		domain.ActionTypeQueryDOM:       "Запрос DOM",
		domain.ActionTypeClickAtPosition: "Клик по координатам",
	}
	if name, ok := names[actionType]; ok {
		return name
	}
	return string(actionType)
}

func translateReason(reason string) string {
	translations := map[string]string{
		"attempting to delete data":          "Попытка удаления данных",
		"attempting to delete account":       "Попытка удаления аккаунта",
		"attempting financial transaction":   "Финансовая транзакция",
		"submitting sensitive information":   "Отправка конфиденциальной информации",
		"attempting to send email":           "Отправка письма",
		"attempting to change settings":      "Изменение настроек",
		"attempting to send job application": "Отправка отклика на вакансию",
		"attempting to place order":          "Оформление заказа",
		"attempting to delete email":         "Удаление письма",
	}
	if translated, ok := translations[reason]; ok {
		return translated
	}
	return reason
}

func translateSuggestion(suggestion string) string {
	translations := map[string]string{
		"Data may be permanently lost":      "Данные могут быть безвозвратно потеряны",
		"Consider backing up first":         "Рекомендуется сначала сделать резервную копию",
		"This action cannot be undone":      "Это действие нельзя отменить",
		"All data will be permanently lost": "Все данные будут безвозвратно потеряны",
		"This may involve real money":       "Это может затронуть реальные деньги",
		"Verify payment details carefully":  "Внимательно проверьте данные платежа",
		"Verify the information is correct": "Убедитесь, что информация корректна",
		"Check you're on the right website": "Убедитесь, что вы на правильном сайте",
		"Email will be sent to recipients":  "Письмо будет отправлено получателям",
		"Application will be sent to employer": "Отклик будет отправлен работодателю",
		"Order will be placed":              "Заказ будет оформлен",
		"Email will be moved to trash":      "Письмо будет перемещено в корзину",
	}
	if translated, ok := translations[suggestion]; ok {
		return translated
	}
	return suggestion
}
