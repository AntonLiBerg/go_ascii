package main

import (
	"go_ascii/internal/scenario"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func TestSelectedScenarioServicesAreRegistered(t *testing.T) {
	name, err := scenario.ReadCurrentScenario(filepath.Join(scenariosDirectory, "current_scenario.txt"))
	if err != nil {
		t.Fatalf("ReadCurrentScenario returned error: %v", err)
	}
	serviceNames, err := scenario.ReadStartupServices(filepath.Join(scenariosDirectory, name, "startup_config.txt"))
	if err != nil {
		t.Fatalf("ReadStartupServices returned error: %v", err)
	}
	services, err := servicesFromNames(serviceNames)
	if err != nil {
		t.Fatalf("servicesFromNames returned error: %v", err)
	}

	resolvedNames := make([]string, len(services))
	for i, service := range services {
		resolvedNames[i] = reflect.TypeOf(service).Name()
	}
	if !slices.Equal(resolvedNames, serviceNames) {
		t.Fatalf("expected services %v, got %v", serviceNames, resolvedNames)
	}
}

func TestServicesFromNamesRejectsUnknownService(t *testing.T) {
	if _, err := servicesFromNames([]string{"UnknownService"}); err == nil {
		t.Fatal("expected unknown service error")
	}
}
