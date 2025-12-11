package domain

// PageContext представляет контекст веб-страницы
type PageContext struct {
	URL              string
	Title            string
	InteractiveElems []Element
	VisibleText      string
	Metadata         map[string]string
}

// Element представляет интерактивный элемент на странице
type Element struct {
	Tag       string
	Text      string
	Selector  string
	Type      string
	Visible   bool
	Clickable bool
	Href      string
	ID        string
	Classes   []string
}
