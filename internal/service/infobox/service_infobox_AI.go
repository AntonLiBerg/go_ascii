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
	contentLines := strings.Split(content, "\n")

	for i := 1; i < len(box)-1; i++ {
		line := []rune(box[i])
		if len(line) < 2 {
			return nil, fmt.Errorf("infobox content line %d is too short", i)
		}

		text := ""
		if i-1 < len(contentLines) {
			text = contentLines[i-1]
		}
		if len([]rune(text)) > len(line)-2 {
			return nil, fmt.Errorf("infobox content is too wide for line %d", i)
		}

		result[i] = string(line[0]) + text + strings.Repeat(" ", len(line)-2-len([]rune(text))) + string(line[len(line)-1])
	}

	if len(contentLines) > len(box)-2 {
		return nil, fmt.Errorf("infobox content has too many lines")
	}
	return result, nil
}
