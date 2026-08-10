package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"strings"
)

func normalizeInputProfileType(value string) string {
	switch value {
	case "none", "terminal", "control":
		return inputProfileTypePrefix + value
	default:
		return value
	}
}

func getEntitiesAndInputProfiles(contentText string) (map[rune]string, map[string]map[string][]string, map[string]map[string]string, map[string]map[rune]struct{}, error) {
	entities := make(map[rune]string)
	components := make(map[string]map[string][]string)
	inputProfiles := make(map[string]map[string]string)
	groups := make(map[string]map[rune]struct{})
	text := strings.ReplaceAll(contentText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	inputProfileStart := -1
	for i, line := range lines {
		trimmed := strings.TrimLeft(strings.TrimSpace(line), SectionNameDivider)
		if trimmed == SectionNameInputProfile && inputProfileStart == -1 {
			inputProfileStart = i + 1
		}
	}

	currentGroup := ""
	currentEntity := ""
	definedEntities := make(map[string]struct{})
	entityEnd := len(lines)
	if inputProfileStart != -1 {
		entityEnd = inputProfileStart - 1
	}
	for lineNumber, line := range lines[:entityEnd] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "===") {
			groupName := strings.TrimSpace(strings.TrimLeft(line, SectionNameDivider))
			if groupName == "" {
				return nil, nil, nil, nil, fmt.Errorf("group header on line %d has no name", lineNumber+1)
			}
			if _, exists := groups[groupName]; exists {
				return nil, nil, nil, nil, fmt.Errorf("duplicate entity group %q", groupName)
			}
			groups[groupName] = make(map[rune]struct{})
			currentGroup = groupName
			currentEntity = ""
			continue
		}
		if currentGroup == "" {
			return nil, nil, nil, nil, fmt.Errorf("entity %q has no group", line)
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
				if componentName == component.NameASCII && valueText == "SPACE" {
					values = append(values, " ")
				} else if valueText != "" {
					for _, value := range strings.Split(valueText, ",") {
						value = strings.TrimSpace(value)
						if value != "" {
							values = append(values, value)
						}
					}
				}
			}
			if currentEntity != "" {
				if _, exists := components[currentEntity][componentName]; exists {
					return nil, nil, nil, nil, fmt.Errorf("duplicate component %q for entity %q", componentName, currentEntity)
				}
				components[currentEntity][componentName] = values
			}
			continue
		}

		keyText, name, ok := strings.Cut(line, ":")
		key := []rune(strings.TrimSpace(keyText))
		name = strings.TrimSpace(name)
		if !ok || len(key) != 1 || name == "" {
			return nil, nil, nil, nil, fmt.Errorf("invalid entity header %q: expected key:name", line)
		}
		if existingName, exists := entities[key[0]]; exists {
			return nil, nil, nil, nil, fmt.Errorf("duplicate entity key %q for %q and %q", key[0], existingName, name)
		}
		if _, exists := definedEntities[name]; exists {
			return nil, nil, nil, nil, fmt.Errorf("duplicate entity name %q", name)
		}
		entities[key[0]] = name
		definedEntities[name] = struct{}{}
		components[name] = make(map[string][]string)
		groups[currentGroup][key[0]] = struct{}{}
		currentEntity = name
	}

	if inputProfileStart != -1 {
		inputProfileEnd := len(lines)
		for i := inputProfileStart; i < len(lines); i++ {
			trimmed := strings.TrimLeft(strings.TrimSpace(lines[i]), SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameInputProfile {
				inputProfileEnd = i
				break
			}
		}
		currentProfile := ""
		for _, line := range lines[inputProfileStart:inputProfileEnd] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "- ") {
				if _, exists := inputProfiles[line]; exists {
					return nil, nil, nil, nil, fmt.Errorf("duplicate input profile %q", line)
				}
				inputProfiles[line] = make(map[string]string)
				currentProfile = line
				continue
			}
			if currentProfile == "" {
				return nil, nil, nil, nil, fmt.Errorf("input binding %q has no profile", line)
			}
			binding := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			separator := strings.IndexAny(binding, ":=")
			if separator == -1 {
				return nil, nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}
			action := strings.TrimSpace(binding[:separator])
			button := strings.TrimSpace(binding[separator+1:])
			if action == "" || button == "" {
				return nil, nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}
			if _, exists := inputProfiles[currentProfile][action]; exists {
				return nil, nil, nil, nil, fmt.Errorf("duplicate input binding %q for profile %q", action, currentProfile)
			}
			if action == inputProfileTypeName {
				button = normalizeInputProfileType(button)
			}
			inputProfiles[currentProfile][action] = button
		}
	}

	return entities, components, inputProfiles, groups, nil
}
