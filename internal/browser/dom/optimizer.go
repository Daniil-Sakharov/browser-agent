package dom

import "github.com/Daniil-Sakharov/BrowserAgent/internal/domain"

// optimizeElements сортирует элементы по приоритету
func (e *Extractor) optimizeElements(elements []domain.Element) []domain.Element {
	filtered := make([]domain.Element, 0, len(elements))
	for _, elem := range elements {
		if elem.Selector != "" && (elem.Text != "" || elem.Type == "input") {
			filtered = append(filtered, elem)
		}
	}
	return prioritizeElements(filtered)
}

func prioritizeElements(elements []domain.Element) []domain.Element {
	var buttons, inputs, links, others []domain.Element

	for _, elem := range elements {
		switch elem.Type {
		case "button":
			buttons = append(buttons, elem)
		case "input":
			inputs = append(inputs, elem)
		case "link":
			links = append(links, elem)
		default:
			others = append(others, elem)
		}
	}

	result := make([]domain.Element, 0, len(elements))
	result = append(result, buttons...)
	result = append(result, inputs...)
	result = append(result, links...)
	result = append(result, others...)
	return result
}
