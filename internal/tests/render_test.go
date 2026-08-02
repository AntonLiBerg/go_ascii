package tests

import (
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPlayerCoversEntityAtSamePosition(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	position := [2]int{1, 1}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"D"}, "door": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	output := captureTerminalOutput(t, gameWorld)
	if !strings.Contains(output, "o") {
		t.Fatalf("expected player in output %q", output)
	}
	if strings.Contains(output, "D") {
		t.Fatalf("expected door to be covered by player in output %q", output)
	}
}

func TestUpdateTerminalSelectsActiveRoom(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{0, 0}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	}); err != nil {
		t.Fatalf("add player: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{1, 0}, map[string][]string{
		"pos": {}, "ascii": {"B"},
	}); err != nil {
		t.Fatalf("add bridge entity: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("comms", [2]int{1, 0}, map[string][]string{
		"pos": {}, "ascii": {"C"},
	}); err != nil {
		t.Fatalf("add comms entity: %v", err)
	}

	tests := []struct {
		name   string
		view   string
		shown  string
		hidden string
	}{
		{name: "player room", shown: "B", hidden: "C"},
		{name: "terminal view", view: "comms", shown: "C", hidden: "B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld.ViewRoom = tt.view
			output := captureTerminalOutput(t, gameWorld)
			if !strings.Contains(output, tt.shown) {
				t.Fatalf("expected active room entity %q in output %q", tt.shown, output)
			}
			if strings.Contains(output, tt.hidden) {
				t.Fatalf("expected inactive room entity %q to be hidden in output %q", tt.hidden, output)
			}
		})
	}
}

func captureTerminalOutput(t *testing.T, gameWorld world.World) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create output pipe: %v", err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer
	render.UpdateTerminal(gameWorld)
	os.Stdout = originalStdout
	if err := writer.Close(); err != nil {
		t.Fatalf("close output pipe: %v", err)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read terminal output: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close output reader: %v", err)
	}
	return string(output)
}

func TestHigherLayerCoversLowerLayer(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	position := [2]int{1, 1}
	if err := gameWorld.AddEntityAtLayer(position, 0, map[string][]string{
		"pos": {}, "ascii": {"."},
	}); err != nil {
		t.Fatalf("AddEntityAtLayer returned error: %v", err)
	}
	if err := gameWorld.AddEntityAtLayer(position, 1, map[string][]string{
		"pos": {}, "ascii": {"D"}, "door": {},
	}); err != nil {
		t.Fatalf("AddEntityAtLayer returned error: %v", err)
	}

	output := captureTerminalOutput(t, gameWorld)
	if !strings.Contains(output, "D") {
		t.Fatalf("expected higher-layer door in output %q", output)
	}
	if strings.Contains(output, ".") {
		t.Fatalf("expected lower-layer entity to be covered in output %q", output)
	}
}
