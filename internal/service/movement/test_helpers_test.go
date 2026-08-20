package movement_test

import (
	"testing"

	"go_ascii/internal/world"
)

func addTestEntity(t *testing.T, gameWorld *world.World, position [2]int, components map[string][]string) int {
	t.Helper()
	eID := len(gameWorld.Entities)
	if err := gameWorld.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	return eID
}
