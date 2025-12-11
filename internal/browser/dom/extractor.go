package dom

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// Extractor извлекает контекст страницы
type Extractor struct {
	maxTextChars int
}

// NewExtractor создаёт новый экстрактор
func NewExtractor() *Extractor {
	return &Extractor{maxTextChars: 20000}
}

// ExtractContext извлекает контекст страницы
func (e *Extractor) ExtractContext(ctx context.Context, page *rod.Page) (*domain.PageContext, error) {
	info, err := page.Info()
	if err != nil {
		return e.emptyContext("Страница недоступна"), nil
	}

	if info.URL == "" || info.URL == "about:blank" {
		return e.emptyContext("Пустая страница"), nil
	}

	html, err := page.HTML()
	if err != nil {
		return &domain.PageContext{URL: info.URL, Title: info.Title, Metadata: map[string]string{}}, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	elements := e.extractElements(doc)
	elements = e.optimizeElements(elements)

	logger.Info(ctx, "✅ Page context extracted")
	return &domain.PageContext{
		URL: info.URL, Title: info.Title,
		InteractiveElems: elements,
		VisibleText:      e.extractVisibleText(doc),
		Metadata:         map[string]string{},
	}, nil
}

func (e *Extractor) emptyContext(title string) *domain.PageContext {
	return &domain.PageContext{
		Title: title, InteractiveElems: []domain.Element{},
		VisibleText: "Используйте navigate для перехода на сайт.",
		Metadata:    map[string]string{},
	}
}

func (e *Extractor) extractElements(doc *goquery.Document) []domain.Element {
	var elems []domain.Element
	selectors := map[string]string{
		"button, input[type=submit], [role=button]": "button",
		"a[href]": "link",
		"input[type=text], input[type=email], input[type=search], textarea": "input",
		"select": "select",
	}

	for sel, elemType := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			elem := e.createElement(s, elemType)
			if elemType == "link" {
				elem.Href, _ = s.Attr("href")
			}
			elems = append(elems, elem)
		})
	}
	return elems
}

func (e *Extractor) createElement(s *goquery.Selection, elemType string) domain.Element {
	id, _ := s.Attr("id")
	return domain.Element{
		Tag: goquery.NodeName(s), Text: strings.TrimSpace(s.Text()),
		Selector: generateSelector(s), Type: elemType, Visible: true,
		Clickable: elemType == "button" || elemType == "link",
		ID:        id, Classes: strings.Split(s.AttrOr("class", ""), " "),
	}
}

func (e *Extractor) extractVisibleText(doc *goquery.Document) string {
	doc.Find("script, style, noscript").Remove()
	bodyText := doc.Find("body").Text()

	lines := strings.Split(bodyText, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 2 {
			cleaned = append(cleaned, line)
		}
	}

	text := strings.Join(cleaned, " ")
	if len(text) > e.maxTextChars {
		text = text[:e.maxTextChars] + "..."
	}
	return text
}
