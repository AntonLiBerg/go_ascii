package tests

import (
	"go_ascii/internal/world"
	"testing"
)

func TestNewUserInputProfileUsesConfiguredKeys(t *testing.T) {
	profile := world.NewUserInputProfile(map[string]string{
		"quitgame":  "q",
		"exit":      "x",
		"moveup":    "w",
		"moveleft":  "a",
		"movedown":  "s",
		"moveright": "d",
		"interact":  "e",
	})
	want := world.UserInputProfile{
		KeyQuitGame:  "q",
		KeyExit:      "x",
		KeyMoveUp:    "w",
		KeyMoveLeft:  "a",
		KeyMoveDown:  "s",
		KeyMoveRight: "d",
		KeyInteract:  "e",
	}

	if profile != want {
		t.Fatalf("expected profile %+v, got %+v", want, profile)
	}
}
