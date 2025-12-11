package domain

// Decision представляет решение AI о следующем действии
type Decision struct {
	Action     Action
	Reasoning  string
	Confidence float64
	Complete   bool
	Result     string
	ToolUseID  string
}

// DecisionRequest запрос для принятия решения
type DecisionRequest struct {
	Task        string
	PageContext *PageContext
	History     []ActionResult
	Step        int
}
