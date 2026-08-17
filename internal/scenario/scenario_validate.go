package scenario

import (
	"fmt"
	"go_ascii/internal/helpers"
	"slices"
	"strings"
)

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
			name := strings.TrimSpace(strings.TrimPrefix(line, SectionDivider))
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

func isValidMapFile(roomLines []string) error {
	if len(roomLines) == 0 {
		return fmt.Errorf("roomLines is empty!")
	}

	rooms := make(map[string]struct{})
	currentRoom := ""
	featuresFound := false
	features := make(map[string]struct{})

	for lineNumber, rawLine := range roomLines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			return fmt.Errorf("room contains empty lines on line %d", lineNumber+1)
		}

		if isSectionHeader(line) {
			roomName, _ := makeRoomHeaderParts(line)
			if roomName == "" {
				return fmt.Errorf("room header on line %d has no name", lineNumber+1)
			}
			if _, exists := rooms[roomName]; exists {
				return fmt.Errorf("duplicate room %q", roomName)
			}
			rooms[roomName] = struct{}{}
			currentRoom = roomName
			featuresFound = false
			features = make(map[string]struct{})
			continue
		}
		if currentRoom == "" {
			return fmt.Errorf("map content on line %d has no room", lineNumber+1)
		}

		if line == SectionNameFeatures {
			if featuresFound {
				return fmt.Errorf("duplicate features section in room %q", currentRoom)
			}
			featuresFound = true
			continue
		}
		if !featuresFound {
			continue // ASCII room content
		}
		if !strings.HasPrefix(line, "- ") {
			return fmt.Errorf("invalid feature on line %d", lineNumber+1)
		}

		feature := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		name, values, ok := strings.Cut(feature, ":")
		if !ok || strings.TrimSpace(name) == "" {
			return fmt.Errorf("invalid feature on line %d", lineNumber+1)
		}
		name = strings.TrimSpace(name)
		if _, exists := features[name]; exists {
			return fmt.Errorf("duplicate feature %q in room %q", name, currentRoom)
		}
		features[name] = struct{}{}

		args := strings.Split(strings.TrimSpace(strings.SplitN(values, "//", 2)[0]), ",")
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
			if args[i] == "" {
				return fmt.Errorf("empty value for feature %q on line %d", name, lineNumber+1)
			}
		}
		switch name {
		case FeatureGround, FeatureInputProfile:
			if len(args) != 1 {
				return fmt.Errorf("feature %q on line %d expects 1 value", name, lineNumber+1)
			}
		case FeaturePortal:
			if len(args) != 3 {
				return fmt.Errorf("feature %q on line %d expects 3 values", name, lineNumber+1)
			}
		case FeatureTerminal:
			if len(args) != 2 {
				return fmt.Errorf("feature %q on line %d expects 2 values", name, lineNumber+1)
			}
		case FeatureSelectableOrder:
			if len(args) == 0 {
				return fmt.Errorf("feature %q on line %d requires a value", name, lineNumber+1)
			}
		default:
			return fmt.Errorf("unknown feature %q on line %d", name, lineNumber+1)
		}
	}
	return nil
}

func isUiFileValid(lines []string) error {
	if len(lines) <= 1 {
		return fmt.Errorf("empty file or just the header!")
	}
	if strings.TrimSpace(lines[0]) != SectionDivider+SectionNameUILayout {
		return fmt.Errorf("no layout section at start of file!")
	}

	uiLayout := make([]string, 0)
	layoutEnd := 1
	for ; layoutEnd < len(lines); layoutEnd++ {
		line := lines[layoutEnd]
		if isSectionHeader(line) {
			break
		}
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			uiLayout = append(uiLayout, trimmedLine)
		}
	}

	uiNonLayoutSections := lines[layoutEnd:]
	sectionHeaders := helpers.Filter(uiNonLayoutSections, isSectionHeader)
	sectionHeaders = helpers.Transform(sectionHeaders, func(header string) string {
		return strings.TrimSpace(strings.TrimPrefix(header, SectionDivider))
	})

	if !slices.Contains(uiLayout, UILayoutRoom) {
		return fmt.Errorf("layout does not contain required keyword room")
	}
	if !helpers.IsUnique(uiLayout) {
		return fmt.Errorf("duplicate headers in layout!")
	}
	if !helpers.IsUnique(sectionHeaders) {
		return fmt.Errorf("duplicate headers in UI file!")
	}
	if !helpers.IsAllS1InS2(sectionHeaders, uiLayout) {
		return fmt.Errorf("some sections are not mentioned in the layout!")
	}
	uiSectionsInLayout := helpers.Filter(uiLayout, func(name string) bool {
		return name != UILayoutRoom
	})
	if !helpers.IsAllS1InS2(uiSectionsInLayout, sectionHeaders) {
		return fmt.Errorf("some sections mentioned in layout do not exist in the UI file")
	}
	return nil
}

func hasAtMostOneUIFilePath(paths []string) bool {
	return len(paths) <= 1
}

func isSectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, SectionDivider) && strings.Count(line, SectionNameDivider) == 3
}
