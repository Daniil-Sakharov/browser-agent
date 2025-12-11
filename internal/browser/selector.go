package browser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// generateSelector генерирует уникальный CSS селектор для элемента
func (e *Extractor) generateSelector(s *goquery.Selection) string {
	// Приоритет 1: ID
	if id, exists := s.Attr("id"); exists && id != "" {
		return "#" + id
	}

	// Приоритет 2: data-* атрибуты
	for _, attr := range []string{"data-testid", "data-qa", "data-test", "data-cy"} {
		if val, exists := s.Attr(attr); exists && val != "" {
			return fmt.Sprintf("[%s='%s']", attr, val)
		}
	}

	// Приоритет 3: name атрибут
	if name, exists := s.Attr("name"); exists && name != "" {
		tag := goquery.NodeName(s)
		return fmt.Sprintf("%s[name='%s']", tag, name)
	}

	// Приоритет 4: комбинация тег + класс (первый класс)
	if class, exists := s.Attr("class"); exists && class != "" {
		classes := strings.Fields(class)
		if len(classes) > 0 {
			tag := goquery.NodeName(s)
			return fmt.Sprintf("%s.%s", tag, classes[0])
		}
	}

	// Приоритет 5: тег + индекс
	tag := goquery.NodeName(s)
	parent := s.Parent()
	if parent.Length() > 0 {
		// Получаем индекс среди соседних элементов
		index := 0
		parent.Children().EachWithBreak(func(i int, child *goquery.Selection) bool {
			if child.Get(0) == s.Get(0) {
				index = i
				return false // прерываем цикл
			}
			return true
		})
		return fmt.Sprintf("%s:nth-child(%d)", tag, index+1)
	}

	// Fallback: просто тег
	return tag
}
