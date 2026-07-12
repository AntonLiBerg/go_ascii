package user

import "testing"

func TestNewUserInputProfileUsesConfiguredButtons(t *testing.T) {
	profile := NewUserInputProfile(map[string]string{
		string(Key_quitGame):  "q",
		string(Key_moveUp):    "w",
		string(Key_moveLeft):  "a",
		string(Key_moveDown):  "s",
		string(Key_moveRight): "d",
		string(KEY_INTERACT):  "e",
	})

	if profile.KeyQuitGame != "q" {
		t.Fatalf("expected quit key q, got %q", profile.KeyQuitGame)
	}
	if profile.KeyMoveUp != "w" {
		t.Fatalf("expected up key w, got %q", profile.KeyMoveUp)
	}
	if profile.KeyMoveLeft != "a" {
		t.Fatalf("expected left key a, got %q", profile.KeyMoveLeft)
	}
	if profile.KeyMoveDown != "s" {
		t.Fatalf("expected down key s, got %q", profile.KeyMoveDown)
	}
	if profile.KeyMoveRight != "d" {
		t.Fatalf("expected right key d, got %q", profile.KeyMoveRight)
	}
	if profile.KeyInteract != "e" {
		t.Fatalf("expected interact key e, got %q", profile.KeyInteract)
	}
}
