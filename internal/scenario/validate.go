package scenario

import (
	"fmt"
	hlp "go_ascii/internal/helpers"
	"slices"
	"strings"
)

func isValidMapFile(roomLines []string) error {
	if len(roomLines) == 0 {
		return fmt.Errorf("roomLines is empty!")
	}
	if !isSectionHeader(roomLines[0]){
		return fmt.Errorf("incorrect header!")
	}
	if helpers.Any(roomLines,func(s string)bool{return s == ""}){
		return fmt.Errorf("room contains empty lines!")
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
	sectionHeaders := hlp.Filter(uiNonLayoutSections, isSectionHeader)
	sectionHeaders = hlp.Transform(sectionHeaders, func(header string) string {
		return strings.TrimSpace(strings.TrimPrefix(header, SectionDivider))
	})

	if !slices.Contains(uiLayout, "room") {
		return fmt.Errorf("layout does not contain required keyword room")
	}
	if !hlp.IsUnique(uiLayout) {
		return fmt.Errorf("duplicate headers in layout!")
	}
	if !hlp.IsUnique(sectionHeaders) {
		return fmt.Errorf("duplicate headers in UI file!")
	}
	if !hlp.IsAllS1InS2(sectionHeaders, uiLayout) {
		return fmt.Errorf("some sections are not mentioned in the layout!")
	}
	uiSectionsInLayout := hlp.Filter(uiLayout, func(name string) bool {
		return name != "room"
	})
	if !hlp.IsAllS1InS2(uiSectionsInLayout, sectionHeaders) {
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
