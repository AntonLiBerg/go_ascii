package aiapi

import "strings"

// GetAsciiMap returns every rune in mapText keyed by its zero-based x,y coordinate.
//
// mapText may be either raw map rows or the full map.txt content with MAP and
// ENTITY sections.
func (api AiAPI) GetAsciiMap(mapText string) map[[2]int]rune {
	asciiMap := make(map[[2]int]rune)
	mapText = strings.Trim(extractMapSection(mapText), "\r\n")
	if mapText == "" {
		return asciiMap
	}

	for y, line := range strings.Split(normalizeLineEndings(mapText), "\n") {
		for x, char := range []rune(line) {
			asciiMap[[2]int{x, y}] = char
		}
	}

	return asciiMap
}

func extractMapSection(text string) string {
	text = normalizeLineEndings(text)
	lines := strings.Split(text, "\n")

	start := -1
	for i, line := range lines {
		if isSectionLine(line, "MAP") {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return text
	}

	end := len(lines)
	for i := start; i < len(lines); i++ {
		if isSectionLine(lines[i], "ENTITY") {
			end = i
			break
		}
	}

	return strings.Join(lines[start:end], "\n")
}

func isSectionLine(line string, name string) bool {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, "=")
	return line == name
}

func normalizeLineEndings(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\r", "\n")
}
