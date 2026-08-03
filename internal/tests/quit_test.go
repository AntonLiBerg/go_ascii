package tests

import (
	"go_ascii/internal/service/quit"
	"go_ascii/internal/world"
	"testing"
)

func TestQuitUpdateDoesNotMutateInput(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyQuitGame: "q"}
	gameWorld.KeyDown = "q"

	update := quit.ServiceQuitGame{}.GetUpdateFunc(gameWorld)
	next := applyUpdate(t, update, gameWorld)

	if gameWorld.ShouldQuit {
		t.Fatal("expected quit update not to mutate its input")
	}
	if !next.ShouldQuit {
		t.Fatal("expected quit update to set ShouldQuit")
	}
}
