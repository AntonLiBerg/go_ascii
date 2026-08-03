package tests

import (
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"
	"strings"
	"testing"
)

func TestPlayerCoversEntityAtSamePosition(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	position := [2]int{1, 1}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"D"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	output := render.TerminalFrame(gameWorld)
	if !strings.Contains(output, "o") {
		t.Fatalf("expected player in output %q", output)
	}
	if strings.Contains(output, "D") {
		t.Fatalf("expected door to be covered by player in output %q", output)
	}
}

func TestTerminalFrameSelectsActiveRoom(t *testing.T) {
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
			output := render.TerminalFrame(gameWorld)
			if !strings.Contains(output, tt.shown) {
				t.Fatalf("expected active room entity %q in output %q", tt.shown, output)
			}
			if strings.Contains(output, tt.hidden) {
				t.Fatalf("expected inactive room entity %q to be hidden in output %q", tt.hidden, output)
			}
		})
	}
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
		"pos": {}, "ascii": {"D"},
	}); err != nil {
		t.Fatalf("AddEntityAtLayer returned error: %v", err)
	}

	output := render.TerminalFrame(gameWorld)
	if !strings.Contains(output, "D") {
		t.Fatalf("expected higher-layer door in output %q", output)
	}
	if strings.Contains(output, ".") {
		t.Fatalf("expected lower-layer entity to be covered in output %q", output)
	}
}

func TestTerminalFrameDrawsLayoutInOrder(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UILayout = []string{"room", "infobox"}
	gameWorld.UIContent["infobox"] = []string{"+---+", "| i |", "+---+"}
	if err := gameWorld.AddEntity([2]int{0, 0}, map[string][]string{
		"pos": {}, "ascii": {"R"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	output := render.TerminalFrame(gameWorld)
	roomIndex := strings.Index(output, "R")
	uiIndex := strings.Index(output, "+---+")
	if roomIndex == -1 || uiIndex == -1 {
		t.Fatalf("expected room and infobox in output %q", output)
	}
	if roomIndex > uiIndex {
		t.Fatalf("expected room before infobox, got output %q", output)
	}
}

func TestTerminalFrameCentersRoomWithinUIWidth(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UILayout = []string{"room", "infobox"}
	gameWorld.UIContent["infobox"] = []string{"----------"}
	if err := gameWorld.AddEntity([2]int{0, 0}, map[string][]string{
		"pos": {}, "ascii": {"R"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	if err := gameWorld.AddEntity([2]int{3, 0}, map[string][]string{
		"pos": {}, "ascii": {"R"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	output := render.TerminalFrame(gameWorld)
	if !strings.Contains(output, "\033[1;4HR") {
		t.Fatalf("expected four-column room centering offset, got output %q", output)
	}
}

func TestTerminalFrameDoesNotCenterRoomWiderThanUI(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UILayout = []string{"room", "infobox"}
	gameWorld.UIContent["infobox"] = []string{"----------"}
	if err := gameWorld.AddEntity([2]int{0, 0}, map[string][]string{
		"pos": {}, "ascii": {"R"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	if err := gameWorld.AddEntity([2]int{10, 0}, map[string][]string{
		"pos": {}, "ascii": {"R"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	output := render.TerminalFrame(gameWorld)
	if !strings.Contains(output, "\033[1;1HR") {
		t.Fatalf("expected wide room to remain left-aligned, got output %q", output)
	}
}
