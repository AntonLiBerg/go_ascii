package scenario

import (
	"fmt"
	"os"
	"strconv"
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
		// Spaces are transparent so lower layers remain visible.
		for x, char := range []rune(line) {
			asciiMap[[2]int{x, y}] = char
		}
	}

	return asciiMap
}

func GetScenarioFromFiles(mapFilePath string, contentFilePath string) (map[int]map[[2]int]rune, map[rune]string, map[string]map[string][]string, map[string]string, error) {
	mapContent, err := os.ReadFile(mapFilePath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	content, err := os.ReadFile(contentFilePath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	asciiMap, err := GetLayeredAsciiMap(string(mapContent))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	entities, components, userInputProfile, err := getEntitiesAndUserInputProfile(string(content))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return asciiMap, entities, components, userInputProfile, nil
}

func GetLayeredAsciiMap(mapText string) (map[int]map[[2]int]rune, error) {
	layers := make(map[int]map[[2]int]rune)
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(strings.TrimRight(mapText, "\n"), "\n")
	currentLayer := -1
	expectedLayer := 0
	y := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "===layer[") && strings.HasSuffix(trimmed, "]") {
			layerText := strings.TrimSuffix(strings.TrimPrefix(trimmed, "===layer["), "]")
			layer, err := strconv.Atoi(layerText)
			if err != nil || layer != expectedLayer {
				return nil, fmt.Errorf("expected layer %d, got %q", expectedLayer, trimmed)
			}
			layers[layer] = make(map[[2]int]rune)
			currentLayer = layer
			expectedLayer++
			y = 0
			continue
		}
		if strings.HasPrefix(trimmed, "===layer") {
			return nil, fmt.Errorf("invalid layer header %q", trimmed)
		}
		if currentLayer == -1 {
			if trimmed == "" {
				continue
			}
			return nil, fmt.Errorf("map content must start with ===layer[0]")
		}

		for x, char := range []rune(line) {
			if char != ' ' {
				layers[currentLayer][[2]int{x, y}] = char
			}
		}
		y++
	}

	if _, ok := layers[0]; !ok {
		return nil, fmt.Errorf("map content must start with ===layer[0]")
	}
	return layers, nil
}

func getEntitiesAndUserInputProfile(contentText string) (map[rune]string, map[string]map[string][]string, map[string]string, error) {
	entities := make(map[rune]string)
	components := make(map[string]map[string][]string)
	userInputProfile := make(map[string]string)
	text := strings.ReplaceAll(contentText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	entityStart := -1
	userInputStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		switch trimmed {
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

	currentEntity := ""
	definedEntities := make(map[string]struct{})
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

			if strings.HasPrefix(line, "- ") {
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
				}
				continue
			}

			keyText, name, ok := strings.Cut(line, ":")
			key := []rune(strings.TrimSpace(keyText))
			name = strings.TrimSpace(name)
			if !ok || len(key) != 1 || name == "" {
				return nil, nil, nil, fmt.Errorf("invalid entity header %q: expected key:name", line)
			}
			if existingName, exists := entities[key[0]]; exists {
				return nil, nil, nil, fmt.Errorf("duplicate entity key %q for %q and %q", key[0], existingName, name)
			}
			if _, exists := definedEntities[name]; exists {
				return nil, nil, nil, fmt.Errorf("duplicate entity name %q", name)
			}

			entities[key[0]] = name
			definedEntities[name] = struct{}{}
			components[name] = make(map[string][]string)
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

	return entities, components, userInputProfile, nil
}
