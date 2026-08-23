package scenario

import (
	"fmt"
	"os"
)

func ReadCurrentScenario(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read current scenario file %q: %w", path, err)
	}
	name, err := ParseCurrentScenario(string(content))
	if err != nil {
		return "", fmt.Errorf("parse current scenario file %q: %w", path, err)
	}
	return name, nil
}

func ReadStartupServices(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read startup config %q: %w", path, err)
	}
	services, err := ParseStartupServices(string(content))
	if err != nil {
		return nil, fmt.Errorf("parse startup config %q: %w", path, err)
	}
	return services, nil
}
