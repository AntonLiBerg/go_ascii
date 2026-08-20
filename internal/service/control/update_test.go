package control_test

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
	"testing"
)

func applyUpdate(t *testing.T, update game.UpdateFunc, gameWorld world.World) world.World {
	t.Helper()
	if update.UpdateFunc == nil {
		t.Fatal("expected update function")
	}
	next, err := update.UpdateFunc(gameWorld)
	if err != nil {
		t.Fatalf("update returned error: %v", err)
	}
	return next
}
