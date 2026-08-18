package tests

import (
	"go_ascii/internal/world"
	"testing"
)

func TestNewUserInputProfile(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  world.UserInputProfile
	}{
		{
			name: "configured keys",
			input: map[string]string{
				"quitgame": "q", "exit": "x", "moveup": "w", "moveleft": "a",
				"movedown": "s", "moveright": "d", "interact": "e",
			},
			want: world.UserInputProfile{
				KeyQuitGame: "q", KeyExit: "x", KeyMoveUp: "w", KeyMoveLeft: "a",
				KeyMoveDown: "s", KeyMoveRight: "d", KeyInteract: "e",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := world.NewUserInputProfile(tt.input); got != tt.want {
				t.Fatalf("expected profile %+v, got %+v", tt.want, got)
			}
		})
	}
}
