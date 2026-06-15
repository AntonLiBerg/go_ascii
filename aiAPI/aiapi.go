package aiapi

import (
	"fmt"
	cmp "go_ascii/component"
	wrld "go_ascii/world"
	"os"
	"strings"
)

func UpdateTerminal(world wrld.World) {
	fmt.Print("\033[2J\033[H")

	maxY := 0
	for _, eId := range world.Entities {
		pos, okPos := world.Pos[eId]
		ascii, okAscii := world.Ascii[eId]
		if !okPos || !okAscii || pos.X < 0 || pos.Y < 0 {
			continue
		}

		fmt.Printf("\033[%d;%dH%c", pos.Y+1, pos.X+1, ascii.Ascii)
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	fmt.Printf("\033[%d;1H", maxY+2)
}

const (
	SectionNameEntity           string = "ENTITY"
	SectionNameMap              string = "MAP"
	SectionNameUserInputProfile string = "USERINPUTPROFILE"
	SectionNameDivider          string = "="
)

func GetAsciiMap(mapText string) map[[2]int]rune {
	asciiMap, _, _, _, _ := parseMapFileContent(mapText)
	return asciiMap
}

func GetAsciiMapAndEntitiesFromFile(filePath string) (map[[2]int]rune, map[rune]string, map[string]map[cmp.ComponentName][]string, map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return parseMapFileContent(string(content))
}

func parseMapFileContent(text string) (map[[2]int]rune, map[rune]string, map[string]map[cmp.ComponentName][]string, map[string]string, error) {
	asciiMap := make(map[[2]int]rune)
	entities := make(map[rune]string)
	components := make(map[string]map[cmp.ComponentName][]string)
	userInputProfile := make(map[string]string)

	mapText := strings.Trim(extractMapSection(text), "\r\n")
	if mapText != "" {
		for y, line := range strings.Split(normalizeLineEndings(mapText), "\n") {
			for x, char := range []rune(line) {
				asciiMap[[2]int{x, y}] = char
			}
		}
	}

	entityText := strings.Trim(extractEntitySection(text), "\r\n")
	currentEntity := ""
	for _, line := range strings.Split(normalizeLineEndings(entityText), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if component, values, ok := parseEntityComponentLine(line); ok {
			if currentEntity != "" {
				components[currentEntity][component] = values
				if component == cmp.C_ASCII {
					addAsciiEntities(entities, currentEntity, values)
				}
			}
			continue
		}

		name, ok := parseEntityLine(line)
		if !ok {
			continue
		}

		if _, exists := components[name]; !exists {
			components[name] = make(map[cmp.ComponentName][]string)
		}
		currentEntity = name
	}

	userInputProfileText := strings.Trim(extractUserInputProfileSection(text), "\r\n")
	for _, line := range strings.Split(normalizeLineEndings(userInputProfileText), "\n") {
		action, button, ok := parseUserInputProfileLine(line)
		if ok {
			userInputProfile[action] = button
		}
	}

	return asciiMap, entities, components, userInputProfile, nil
}

func addAsciiEntities(entities map[rune]string, entityName string, values []string) {
	for _, value := range values {
		symbol := []rune(value)
		if len(symbol) == 1 {
			entities[symbol[0]] = entityName
		}
	}
}

func parseEntityLine(line string) (string, bool) {
	name := strings.TrimSpace(line)
	return name, name != ""
}

func parseEntityComponentLine(line string) (cmp.ComponentName, []string, bool) {
	if !strings.HasPrefix(line, "-") {
		return "", nil, false
	}

	component := strings.TrimSpace(strings.TrimPrefix(line, "-"))
	if component == "" {
		return "", nil, false
	}

	separator := strings.IndexAny(component, ":=")
	if separator == -1 {
		return cmp.ComponentName(component), []string{}, true
	}

	name := strings.TrimSpace(component[:separator])
	if name == "" {
		return "", nil, false
	}

	return cmp.ComponentName(name), parseComponentValues(component[separator+1:]), true
}

func parseComponentValues(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	values := []string{}
	for _, value := range strings.Split(text, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			values = append(values, value)
		}
	}

	return values
}

func parseUserInputProfileLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", false
	}

	separator := strings.IndexAny(line, ":=")
	if separator == -1 {
		return "", "", false
	}

	action := strings.TrimSpace(line[:separator])
	button := strings.TrimSpace(line[separator+1:])
	if action == "" || button == "" {
		return "", "", false
	}

	return action, button, true
}

func extractMapSection(text string) string {
	text = normalizeLineEndings(text)
	lines := strings.Split(text, "\n")

	start := -1
	for i, line := range lines {
		if isSectionLine(line, SectionNameMap) {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return text
	}

	end := len(lines)
	for i := start; i < len(lines); i++ {
		if isAnySectionLine(lines[i]) {
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
		if isSectionLine(line, SectionNameEntity) {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}

	end := len(lines)
	for i := start; i < len(lines); i++ {
		if isAnySectionLine(lines[i]) {
			end = i
			break
		}
	}

	return strings.Join(lines[start:end], "\n")
}

func extractUserInputProfileSection(text string) string {
	text = normalizeLineEndings(text)
	lines := strings.Split(text, "\n")

	start := -1
	for i, line := range lines {
		if isSectionLine(line, SectionNameUserInputProfile) {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}

	end := len(lines)
	for i := start; i < len(lines); i++ {
		if isAnySectionLine(lines[i]) {
			end = i
			break
		}
	}

	return strings.Join(lines[start:end], "\n")
}

func isSectionLine(line string, name string) bool {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, SectionNameDivider)
	return line == name
}

func isAnySectionLine(line string) bool {
	return isSectionLine(line, SectionNameMap) ||
		isSectionLine(line, SectionNameEntity) ||
		isSectionLine(line, SectionNameUserInputProfile)
}

func normalizeLineEndings(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\r", "\n")
}
