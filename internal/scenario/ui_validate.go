package scenario

import (
	"fmt"
	"go_ascii/internal/helpers"
	"slices"
	"strings"
)

func isAllS1InS2(slice1 []string, slice2 []string) bool {
	return helpers.All(slice1, func(value string) bool {
		return slices.Contains(slice2, value)
	})
}

func isUnique(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
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
		if isUISectionHeader(line) {
			break
		}
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			uiLayout = append(uiLayout, trimmedLine)
		}
	}

	uiNonLayoutSections := lines[layoutEnd:]
	sectionHeaders := helpers.Filter(uiNonLayoutSections, isUISectionHeader)
	sectionHeaders = helpers.Transform(sectionHeaders, func(header string) string {
		return strings.TrimSpace(strings.TrimPrefix(header, SectionDivider))
	})

	if !slices.Contains(uiLayout, "room") {
		return fmt.Errorf("layout does not contain required keyword room")
	}
	if !isUnique(uiLayout) {
		return fmt.Errorf("duplicate headers in layout!")
	}
	if !isUnique(sectionHeaders) {
		return fmt.Errorf("duplicate headers in UI file!")
	}
	if !isAllS1InS2(sectionHeaders, uiLayout) {
		return fmt.Errorf("some sections are not mentioned in the layout!")
	}
	uiSectionsInLayout := helpers.Filter(uiLayout, func(name string) bool {
		return name != "room"
	})
	if !isAllS1InS2(uiSectionsInLayout, sectionHeaders) {
		return fmt.Errorf("some sections mentioned in layout do not exist in the UI file")
	}
	return nil
}

func hasAtMostOneUIFilePath(paths []string) bool {
	return len(paths) <= 1
}

func isUISectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, SectionDivider) && strings.Count(line, SectionNameDivider) == 3
}
