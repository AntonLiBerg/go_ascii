package scenario

import (
	"fmt"
	"strings"
)

func normalizeInputProfileType(value string) string {
	switch value {
	case InputProfileTypeNone, InputProfileTypeTerminal, InputProfileTypeControl:
		return InputProfileTypeName + value
	default:
		return value
	}
}

type FileContent struct {
	entities      map[rune]string
	components    map[string]map[string][]string
	inputProfiles map[string]map[string]string
	groups        map[string]map[rune]struct{}
}

func ParseContentFile(contentText string) (FileContent, error) {
	contentMap := FileContent{
		entities:      make(map[rune]string),
		components:    make(map[string]map[string][]string),
		inputProfiles: make(map[string]map[string]string),
		groups:        make(map[string]map[rune]struct{}),
	}

	contentText = strings.ReplaceAll(contentText, "\r\n", "\n")
	contentText = strings.ReplaceAll(contentText, "\r", "\n")
	lines := strings.Split(strings.TrimRight(contentText, "\n"), "\n")
	if err := IsValidContentFile(lines); err != nil {
		return contentMap, err
	}
	groupByName := group(lines, func(s string) bool { return isSectionHeader(s) }, sectionName)
	for name, ls := range groupByName {
		if name == SectionNameInputProfile {
			if nInpProfiles, err := MakeInputProfiles(ls); err != nil {
				return contentMap, err
			} else {
				contentMap.inputProfiles = nInpProfiles
			}
			continue
		}
		if nEntities, err := AppendEntities(ls, contentMap.entities); err != nil {
			return contentMap, err
		} else {
			contentMap.entities = nEntities
		}
		if nComps, err := AppendComponents(contentMap, ls); err != nil {
			return contentMap, err
		} else {
			contentMap.components = nComps
		}

		if nGroups, err := AppendGroups(contentMap, name, ls); err != nil {
			return contentMap, err
		} else {
			contentMap.groups = nGroups
		}
	}
	return contentMap, nil
}

func MakeInputProfiles(lines []string) (map[string]map[string]string, error) {
	profiles := make(map[string]map[string]string)
	currentProfile := ""
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if !strings.HasPrefix(line, "- ") {
			currentProfile = line
			profiles[currentProfile] = make(map[string]string)
			continue
		}
		binding := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		action, button, _ := strings.Cut(binding, ":")
		if !strings.Contains(binding, ":") {
			action, button, _ = strings.Cut(binding, "=")
		}
		action = strings.TrimSpace(action)
		button = strings.TrimSpace(button)
		if action == InputProfileTypeName {
			button = normalizeInputProfileType(button)
		}
		profiles[currentProfile][action] = button
	}
	return profiles, nil
}

func MakeEntities(lines []string) (map[rune]string, error) {
	entities := make(map[rune]string)
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "- ") {
			continue
		}
		keyText, name, _ := strings.Cut(line, ":")
		key := []rune(strings.TrimSpace(keyText))
		entities[key[0]] = strings.TrimSpace(name)
	}
	return entities, nil
}

func AppendEntities(lines []string, existing map[rune]string) (map[rune]string, error) {
	entities := make(map[rune]string, len(existing))
	for key, name := range existing {
		entities[key] = name
	}
	parsed, _ := MakeEntities(lines)
	for key, name := range parsed {
		entities[key] = name
	}
	return entities, nil
}

func AppendComponents(contentMap FileContent, lines []string) (map[string]map[string][]string, error) {
	components := make(map[string]map[string][]string, len(contentMap.components))
	for entity, values := range contentMap.components {
		components[entity] = make(map[string][]string, len(values))
		for name, values := range values {
			components[entity][name] = append([]string(nil), values...)
		}
	}
	currentEntity := ""
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if !strings.HasPrefix(line, "- ") {
			_, name, _ := strings.Cut(line, ":")
			currentEntity = strings.TrimSpace(name)
			if _, exists := components[currentEntity]; !exists {
				components[currentEntity] = make(map[string][]string)
			}
			continue
		}
		componentText := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		name, valuesText, ok := strings.Cut(componentText, ":")
		if !ok {
			name, valuesText, _ = strings.Cut(componentText, "=")
		}
		values := make([]string, 0)
		for _, value := range strings.Split(valuesText, ",") {
			if value = strings.TrimSpace(value); value != "" {
				if value == "SPACE" {
					value = " "
				}
				values = append(values, value)
			}
		}
		components[currentEntity][strings.TrimSpace(name)] = values
	}
	return components, nil
}

func AppendGroups(contentMap FileContent, name string, lines []string) (map[string]map[rune]struct{}, error) {
	groups := make(map[string]map[rune]struct{}, len(contentMap.groups)+1)
	for groupName, members := range contentMap.groups {
		groups[groupName] = make(map[rune]struct{}, len(members))
		for member := range members {
			groups[groupName][member] = struct{}{}
		}
	}
	groups[name] = make(map[rune]struct{})
	parsed, _ := MakeEntities(lines)
	for member := range parsed {
		groups[name][member] = struct{}{}
	}
	return groups, nil
}

func IsValidContentFile(lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("content file is empty")
	}

	sections := make(map[string]struct{})
	entities := make(map[rune]string)
	entityNames := make(map[string]struct{})
	currentSection, currentEntity, currentProfile := "", "", ""
	components := make(map[string]map[string]struct{})
	profiles := make(map[string]struct{})
	bindings := make(map[string]map[string]struct{})

	for lineIndex, rawLine := range lines {
		lineNumber := lineIndex + 1
		line := strings.TrimSpace(rawLine)
		if line == "" {
			return fmt.Errorf("content file contains an empty line at line %d", lineNumber)
		}
		if isSectionHeader(line) {
			name := sectionName(line)
			if name == "" {
				return fmt.Errorf("section header on line %d has no name", lineNumber)
			}
			if _, exists := sections[name]; exists {
				return fmt.Errorf("duplicate content section %q", name)
			}
			sections[name] = struct{}{}
			currentSection, currentEntity, currentProfile = name, "", ""
			continue
		}
		if currentSection == "" {
			return fmt.Errorf("content on line %d has no section", lineNumber)
		}

		if currentSection == SectionNameInputProfile {
			if !strings.HasPrefix(line, "- ") {
				if _, exists := profiles[line]; exists {
					return fmt.Errorf("duplicate input profile %q", line)
				}
				profiles[line] = struct{}{}
				bindings[line] = make(map[string]struct{})
				currentProfile = line
				continue
			}
			binding := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			action, button, ok := strings.Cut(binding, ":")
			if !ok {
				action, button, ok = strings.Cut(binding, "=")
			}
			action, button = strings.TrimSpace(action), strings.TrimSpace(button)
			if currentProfile == "" || !ok || action == "" || button == "" {
				return fmt.Errorf("invalid input binding on line %d", lineNumber)
			}
			if _, exists := bindings[currentProfile][action]; exists {
				return fmt.Errorf("duplicate input binding %q for profile %q", action, currentProfile)
			}
			bindings[currentProfile][action] = struct{}{}
			continue
		}

		if strings.HasPrefix(line, "- ") {
			if currentEntity == "" {
				return fmt.Errorf("component on line %d has no entity", lineNumber)
			}
			componentText := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			componentName := componentText
			if separator := strings.IndexAny(componentText, ":="); separator >= 0 {
				componentName = strings.TrimSpace(componentText[:separator])
			}
			if componentName == "" {
				return fmt.Errorf("invalid component on line %d", lineNumber)
			}
			if _, exists := components[currentEntity][componentName]; exists {
				return fmt.Errorf("duplicate component %q for entity %q", componentName, currentEntity)
			}
			components[currentEntity][componentName] = struct{}{}
			continue
		}

		keyText, name, ok := strings.Cut(line, ":")
		key := []rune(strings.TrimSpace(keyText))
		name = strings.TrimSpace(name)
		if !ok || len(key) != 1 || name == "" {
			return fmt.Errorf("invalid entity header on line %d", lineNumber)
		}
		if existing, exists := entities[key[0]]; exists {
			return fmt.Errorf("duplicate entity key %q for %q and %q", key[0], existing, name)
		}
		if _, exists := entityNames[name]; exists {
			return fmt.Errorf("duplicate entity name %q", name)
		}
		entities[key[0]], entityNames[name] = name, struct{}{}
		currentEntity = name
		components[name] = make(map[string]struct{})
	}
	return nil
}

func isSectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, SectionDivider) {
		return false
	}
	remaining := strings.TrimPrefix(line, SectionDivider)
	if strings.HasPrefix(remaining, SectionNameDivider) {
		remaining = strings.TrimPrefix(remaining, SectionNameDivider)
		if strings.HasPrefix(remaining, SectionNameDivider) {
			return false
		}
	}
	return strings.TrimSpace(remaining) != ""
}
