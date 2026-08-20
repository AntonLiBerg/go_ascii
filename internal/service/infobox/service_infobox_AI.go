package infobox

import (
	"fmt"
	"strings"
)

type ai struct{}

var AI ai

// SetInfoboxText replaces the text between an infobox's border lines.
func (ai) SetInfoboxText(box []string, content string) ([]string, error) {
	if len(box) < 3 {
		return nil, fmt.Errorf("infobox must contain borders and at least one content line")
	}

	result := make([]string, len(box))
	copy(result, box)

	width := 0
	for i := 1; i < len(box)-1; i++ {
		lineWidth := len([]rune(box[i])) - 2
		if lineWidth < 0 {
			return nil, fmt.Errorf("infobox content line %d is too short", i)
		}
		if i == 1 || lineWidth < width {
			width = lineWidth
		}
	}
	if width == 0 {
		return nil, fmt.Errorf("infobox content has no available width")
	}

	contentLines := wrapContent(content, width)
	if len(contentLines) > len(box)-2 {
		return nil, fmt.Errorf("infobox content has too many lines")
	}

	for i := 1; i < len(box)-1; i++ {
		line := []rune(box[i])
		text := ""
		if i-1 < len(contentLines) {
			text = contentLines[i-1]
		}
		result[i] = string(line[0]) + text + strings.Repeat(" ", len(line)-2-len([]rune(text))) + string(line[len(line)-1])
	}
	return result, nil
}

func wrapContent(content string, width int) []string {
	wrapped := make([]string, 0)
	for _, line := range strings.Split(content, "\n") {
		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		current := ""
		for _, word := range words {
			for len([]rune(word)) > width {
				part := string([]rune(word)[:width])
				word = string([]rune(word)[width:])
				if current != "" {
					wrapped = append(wrapped, current)
					current = ""
				}
				wrapped = append(wrapped, part)
			}
			if current == "" {
				current = word
			} else if len([]rune(current))+1+len([]rune(word)) <= width {
				current += " " + word
			} else {
				wrapped = append(wrapped, current)
				current = word
			}
		}
		if current != "" {
			wrapped = append(wrapped, current)
		}
	}
	return wrapped
}
