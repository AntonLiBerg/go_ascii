package scenario_test

import (
	"go_ascii/internal/scenario"
	"slices"
	"testing"
)

func TestParseCurrentScenario(t *testing.T) {
	name, err := scenario.ParseCurrentScenario("  skyship\r\n")
	if err != nil {
		t.Fatalf("ParseCurrentScenario returned error: %v", err)
	}
	if name != "skyship" {
		t.Fatalf("expected skyship, got %q", name)
	}
}

func TestParseCurrentScenarioRejectsInvalidContent(t *testing.T) {
	for _, content := range []string{"", "first\nsecond", "../skyship", "folder/skyship"} {
		if _, err := scenario.ParseCurrentScenario(content); err == nil {
			t.Fatalf("expected %q to be rejected", content)
		}
	}
}

func TestParseStartupServices(t *testing.T) {
	content := "services\r\n- ServiceQuitGame\r\n- ServiceMovePlayer\r\n"
	services, err := scenario.ParseStartupServices(content)
	if err != nil {
		t.Fatalf("ParseStartupServices returned error: %v", err)
	}
	want := []string{"ServiceQuitGame", "ServiceMovePlayer"}
	if !slices.Equal(services, want) {
		t.Fatalf("expected services %v, got %v", want, services)
	}
}

func TestParseStartupServicesRejectsInvalidContent(t *testing.T) {
	tests := []string{
		"",
		"systems\n- ServiceQuitGame",
		"services",
		"services\nServiceQuitGame",
		"services\n- ServiceQuitGame\n- ServiceQuitGame",
	}
	for _, content := range tests {
		if _, err := scenario.ParseStartupServices(content); err == nil {
			t.Fatalf("expected %q to be rejected", content)
		}
	}
}
