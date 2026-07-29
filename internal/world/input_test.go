package world

import "testing"

func TestNewUserInputProfileUsesConfiguredKeys(t *testing.T) {
	profile := NewUserInputProfile(map[string]string{
		"quitgame":  "q",
		"moveup":    "w",
		"moveleft":  "a",
		"movedown":  "s",
		"moveright": "d",
		"interact":  "e",
	})
	want := UserInputProfile{
		KeyQuitGame:  "q",
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
