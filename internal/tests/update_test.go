package tests

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
	"testing"
)

func applyUpdate(t *testing.T, update game.UpdateFunc, gameWorld world.World) world.World {
	t.Helper()
	if update.UpdateFunc == nil {
		return gameWorld
	}
	next, err := update.UpdateFunc(gameWorld)
	if err != nil {
		t.Fatalf("UpdateFunc returned error: %v", err)
	}
	return next
}
