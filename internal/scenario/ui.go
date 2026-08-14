package scenario

import "strings"

func getUiLayoutAndUIs(uiFileText string) ([]string, []string, map[string][]string, error) {
	text := strings.ReplaceAll(uiFileText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	if err := isUiFileValid(lines); err != nil {
		return nil, nil, nil, err
	}
	uiLayout := getNextUiSection(lines[1:])
	uiSections := make(map[string][]string)
	remainingLines := lines[len(uiLayout)+1:]
	for len(remainingLines) > 0 {
		sName := strings.TrimSpace(strings.TrimPrefix(remainingLines[0], SectionDivider))
		uiSections[sName] = trimUIContent(getNextUiSection(remainingLines[1:]))
		remainingLines = remainingLines[len(uiSections[sName])+1:]
	}

	uis := make([]string, 0, len(uiLayout))
	for _, name := range uiLayout {
		if name != "room" {
			uis = append(uis, name)
		}
	}
	return uiLayout, uis, uiSections, nil
}

func trimUIContent(lines []string) []string {
	start := 0
	end := len(lines)
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return lines[start:end]
}

func getNextUiSection(lines []string) []string {
	section := []string{}
	for _, line := range lines {
		if isSectionHeader(line) {
			break
		}
		section = append(section, line)
	}
	return section
}
