package domain

import "time"

// Action представляет действие браузера
type Action struct {
	Type      ActionType
	Selector  string
	Value     string
	URL       string
	Direction string
	Query     string // для query_dom
	Question  string // для analyze_page
	FullPage  bool   // для take_screenshot
	X         int    // для click_at_position
	Y         int    // для click_at_position
}

// ActionType тип действия браузера
type ActionType string

const (
	ActionTypeNavigate        ActionType = "navigate"
	ActionTypeClick           ActionType = "click"
	ActionTypeClickAtPosition ActionType = "click_at_position"
	ActionTypeType            ActionType = "type"
	ActionTypeScroll          ActionType = "scroll"
	ActionTypeSelect          ActionType = "select"
	ActionTypeWait            ActionType = "wait"
	ActionTypePressEnter      ActionType = "press_enter"
	ActionTypeCompleteTask    ActionType = "complete_task"
	ActionTypeTakeScreenshot  ActionType = "take_screenshot"
	ActionTypeQueryDOM        ActionType = "query_dom"
	ActionTypeAnalyzePage     ActionType = "analyze_page"
)

// ErrorContext контекст ошибки для адаптации агента
type ErrorContext struct {
	FailedSelector  string   // Селектор который не сработал
	SimilarElements []string // Похожие элементы на странице
	Suggestion      string   // Рекомендация что делать
}

// ActionResult результат выполнения действия
type ActionResult struct {
	Success       bool
	Action        string
	Message       string
	Error         error
	ErrorContext  *ErrorContext // Контекст ошибки для адаптации
	Screenshot    string        // путь к скриншоту
	ScreenshotB64 string        // base64 скриншота для передачи в Claude
	QueryResult   string        // результат query_dom
	Duration      time.Duration
	Timestamp     time.Time
}
