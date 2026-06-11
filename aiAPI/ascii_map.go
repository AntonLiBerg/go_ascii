package aiapi

import (
	"os"
	"strings"
)

// GetAsciiMap returns every rune in mapText keyed by its zero-based x,y coordinate.
//
// mapText may be either raw map rows or the full map.txt content with MAP and
// ENTITY sections.
func (api AiAPI) GetAsciiMap(mapText string) map[[2]int]rune {
	asciiMap, _ := parseMapFileContent(mapText)
	return asciiMap
}

// GetAsciiMapAndEntitiesFromFile reads a map.txt-style file and returns the map
// runes keyed by zero-based col,row plus the entity names keyed to their rune.
//
// If filePath cannot be read, both returned maps are nil.
func (api AiAPI) GetAsciiMapAndEntitiesFromFile(filePath string) (map[[2]int]rune, map[string]rune) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil
	}

	return parseMapFileContent(string(content))
}

func parseMapFileContent(text string) (map[[2]int]rune, map[string]rune) {
	asciiMap := make(map[[2]int]rune)
	entityRunes := make(map[string]rune)

	mapText := strings.Trim(extractMapSection(text), "\r\n")
	if mapText != "" {
		for y, line := range strings.Split(normalizeLineEndings(mapText), "\n") {
			for x, char := range []rune(line) {
				asciiMap[[2]int{x, y}] = char
			}
		}
	}

	entityText := strings.Trim(extractEntitySection(text), "\r\n")
	for _, line := range strings.Split(normalizeLineEndings(entityText), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		symbol := []rune(strings.TrimSpace(parts[0]))
		name := strings.TrimSpace(parts[1])
		if len(symbol) != 1 || name == "" {
			continue
		}

		entityRunes[name] = symbol[0]
	}

	return asciiMap, entityRunes
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

func extractEntitySection(text string) string {
	text = normalizeLineEndings(text)
	lines := strings.Split(text, "\n")

	start := -1
	for i, line := range lines {
		if isSectionLine(line, "ENTITY") {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}

	return strings.Join(lines[start:], "\n")
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
