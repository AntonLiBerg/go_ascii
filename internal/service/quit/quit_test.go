package quit_test

import (
	"go_ascii/internal/service/quit"
	"go_ascii/internal/world"
	"testing"
)

func TestQuitUpdate(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantQuit bool
	}{
		{name: "quit key sets flag", key: "q", wantQuit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := world.NewWorldEmpty()
			gameWorld.UserInputProfile = world.UserInputProfile{KeyQuitGame: tt.key}
			gameWorld.KeyDown = tt.key

			next := applyUpdate(t, quit.ServiceQuitGame{}.GetUpdateFunc(gameWorld), gameWorld)
			if gameWorld.ShouldQuit {
				t.Fatal("expected quit update not to mutate its input")
			}
			if next.ShouldQuit != tt.wantQuit {
				t.Fatalf("expected ShouldQuit=%v, got %v", tt.wantQuit, next.ShouldQuit)
			}
		})
	}
}
