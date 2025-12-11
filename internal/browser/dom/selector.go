package dom

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// generateSelector генерирует CSS селектор для элемента
func generateSelector(s *goquery.Selection) string {
	if id, exists := s.Attr("id"); exists && id != "" {
		return "#" + id
	}

	for _, attr := range []string{"data-testid", "data-qa", "data-test", "data-cy"} {
		if val, exists := s.Attr(attr); exists && val != "" {
			return fmt.Sprintf("[%s='%s']", attr, val)
		}
	}

	if name, exists := s.Attr("name"); exists && name != "" {
		return fmt.Sprintf("%s[name='%s']", goquery.NodeName(s), name)
	}

	if class, exists := s.Attr("class"); exists && class != "" {
		classes := strings.Fields(class)
		if len(classes) > 0 {
			return fmt.Sprintf("%s.%s", goquery.NodeName(s), classes[0])
		}
	}

	tag := goquery.NodeName(s)
	parent := s.Parent()
	if parent.Length() > 0 {
		index := 0
		parent.Children().EachWithBreak(func(i int, child *goquery.Selection) bool {
			if child.Get(0) == s.Get(0) {
				index = i
				return false
			}
			return true
		})
		return fmt.Sprintf("%s:nth-child(%d)", tag, index+1)
	}
	return tag
}
