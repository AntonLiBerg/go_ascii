package scenario

import (
	"os"
	"strings"
)

const (
	SectionNameEntity           string = "ENTITY"
	SectionNameMap              string = "MAP"
	SectionNameUserInputProfile string = "USERINPUTPROFILE"
	SectionNameDivider          string = "="
)

func GetAsciiMap(mapText string) map[[2]int]rune {
	asciiMap := make(map[[2]int]rune)
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(mapText, "\n")

	mapStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		if trimmed == SectionNameMap {
			mapStart = i + 1
			break
		}
	}

	mapEnd := len(lines)
	if mapStart == -1 {
		mapStart = 0
	} else {
		for i := mapStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameUserInputProfile {
				mapEnd = i
				break
			}
		}
	}

	mapSection := strings.Trim(strings.Join(lines[mapStart:mapEnd], "\n"), "\n")
	if mapSection == "" {
		return asciiMap
	}

	for y, line := range strings.Split(mapSection, "\n") {
		for x, char := range []rune(line) {
			asciiMap[[2]int{x, y}] = char
		}
	}

	return asciiMap
}

func GetAsciiMapAndEntitiesFromFile(filePath string) (map[[2]int]rune, map[rune]string, map[string]map[string][]string, map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	asciiMap := make(map[[2]int]rune)
	entities := make(map[rune]string)
	components := make(map[string]map[string][]string)
	userInputProfile := make(map[string]string)
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	mapStart := -1
	entityStart := -1
	userInputStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		switch trimmed {
		case SectionNameMap:
			if mapStart == -1 {
				mapStart = i + 1
			}
		case SectionNameEntity:
			if entityStart == -1 {
				entityStart = i + 1
			}
		case SectionNameUserInputProfile:
			if userInputStart == -1 {
				userInputStart = i + 1
			}
		}
	}

	mapEnd := len(lines)
	if mapStart == -1 {
		mapStart = 0
	} else {
		for i := mapStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameUserInputProfile {
				mapEnd = i
				break
			}
		}
	}

	mapSection := strings.Trim(strings.Join(lines[mapStart:mapEnd], "\n"), "\n")
	if mapSection != "" {
		for y, line := range strings.Split(mapSection, "\n") {
			for x, char := range []rune(line) {
				asciiMap[[2]int{x, y}] = char
			}
		}
	}

	currentEntity := ""
	if entityStart != -1 {
		entityEnd := len(lines)
		for i := entityStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameUserInputProfile {
				entityEnd = i
				break
			}
		}

		for _, line := range lines[entityStart:entityEnd] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "-") {
				componentText := strings.TrimSpace(strings.TrimPrefix(line, "-"))
				if componentText == "" {
					continue
				}

				componentName := componentText
				values := []string{}
				separator := strings.IndexAny(componentText, ":=")
				if separator != -1 {
					name := strings.TrimSpace(componentText[:separator])
					if name == "" {
						continue
					}

					componentName = name
					valueText := strings.TrimSpace(componentText[separator+1:])
					if valueText != "" {
						for _, value := range strings.Split(valueText, ",") {
							value = strings.TrimSpace(value)
							if value != "" {
								values = append(values, value)
							}
						}
					}
				}

				if currentEntity != "" {
					components[currentEntity][componentName] = values
					if componentName == "ascii" {
						for _, value := range values {
							symbol := []rune(value)
							if len(symbol) == 1 {
								entities[symbol[0]] = currentEntity
							}
						}
					}
				}
				continue
			}

			name := strings.TrimSpace(line)
			if _, exists := components[name]; !exists {
				components[name] = make(map[string][]string)
			}
			currentEntity = name
		}
	}

	if userInputStart != -1 {
		userInputEnd := len(lines)
		for i := userInputStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameUserInputProfile {
				userInputEnd = i
				break
			}
		}

		for _, line := range lines[userInputStart:userInputEnd] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			separator := strings.IndexAny(line, ":=")
			if separator == -1 {
				continue
			}

			action := strings.TrimSpace(line[:separator])
			button := strings.TrimSpace(line[separator+1:])
			if action == "" || button == "" {
				continue
			}

			userInputProfile[action] = button
		}
	}

	return asciiMap, entities, components, userInputProfile, nil
}
