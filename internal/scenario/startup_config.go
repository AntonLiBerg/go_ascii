package scenario

import (
	"fmt"
	"strings"
)

const startupServicesHeader = "services"

func ParseCurrentScenario(content string) (string, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	name := strings.TrimSpace(content)

	if name == "" {
		return "", fmt.Errorf("current scenario is empty")
	}
	if strings.Contains(name, "\n") {
		return "", fmt.Errorf("current scenario must contain exactly one scenario name")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("invalid scenario name %q", name)
	}
	return name, nil
}

func ParseStartupServices(content string) ([]string, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	services := make([]string, 0)
	seen := make(map[string]struct{})
	foundHeader := false

	for lineIndex, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		if !foundHeader {
			if line != startupServicesHeader {
				return nil, fmt.Errorf("expected %q header on line %d", startupServicesHeader, lineIndex+1)
			}
			foundHeader = true
			continue
		}
		if !strings.HasPrefix(line, "- ") {
			return nil, fmt.Errorf("invalid service on line %d", lineIndex+1)
		}

		name := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		if name == "" {
			return nil, fmt.Errorf("service on line %d has no name", lineIndex+1)
		}
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf("duplicate service %q", name)
		}
		seen[name] = struct{}{}
		services = append(services, name)
	}

	if !foundHeader {
		return nil, fmt.Errorf("startup config has no %q header", startupServicesHeader)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("startup config has no services")
	}
	return services, nil
}
