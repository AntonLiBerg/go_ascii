package tests

import (
	"bytes"
	"errors"
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"
	"strings"
	"testing"
)

func TestDrawOnTerminalService(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.IterationNr = 1
	gameWorld.HasChanged = true
	var output bytes.Buffer

	update := render.ServiceDrawOnTerminal{Writer: &output}.GetUpdateFunc(gameWorld)
	next := applyUpdate(t, update, gameWorld)

	if update.Order != 100 {
		t.Fatalf("expected render order 100, got %d", update.Order)
	}
	if output.String() != render.TerminalFrame(gameWorld) {
		t.Fatalf("unexpected terminal output %q", output.String())
	}
	if next.HasChanged {
		t.Fatal("expected render service to clear HasChanged")
	}
}

func TestDrawOnTerminalServiceSkipsUnchangedWorld(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.IterationNr = 2
	var output bytes.Buffer

	next := applyUpdate(t, render.ServiceDrawOnTerminal{Writer: &output}.GetUpdateFunc(gameWorld), gameWorld)

	if output.Len() != 0 {
		t.Fatalf("expected no terminal output, got %q", output.String())
	}
	if next.HasChanged {
		t.Fatal("expected unchanged world")
	}
}

func TestDrawOnTerminalServiceReturnsWriterError(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.IterationNr = 1
	wantErr := errors.New("write failed")
	update := render.ServiceDrawOnTerminal{Writer: errorWriter{err: wantErr}}.GetUpdateFunc(gameWorld)

	_, err := update.UpdateFunc(gameWorld)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected writer error, got %v", err)
	}
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

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

func TestTerminalFrameDrawsControlStateASCII(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	if err := gameWorld.AddEntity([2]int{0, 0}, map[string][]string{
		"pos": {}, "ascii": {"x"}, "selectable": {"u", "f", "s", "target"},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	control := &world.ControlNode{
		SelectableEntityID: 0,
	}
	control.Next, control.Prev = control, control

	tests := []struct {
		name          string
		activeControl *world.ControlNode
		editing       bool
		shown         string
	}{
		{name: "unfocused", shown: "u"},
		{name: "focused", activeControl: control, shown: "f"},
		{name: "selected", activeControl: control, editing: true, shown: "s"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gameWorld.ActiveControl = test.activeControl
			gameWorld.EditingControl = test.editing
			output := render.TerminalFrame(gameWorld)
			if !strings.Contains(output, test.shown) {
				t.Fatalf("expected %q in output %q", test.shown, output)
			}
			for _, hidden := range []string{"x", "u", "f", "s"} {
				if hidden != test.shown && strings.Contains(output, hidden) {
					t.Fatalf("expected %q not to appear in output %q", hidden, output)
				}
			}
		})
	}
}

func TestTerminalFrameSelectsActiveRoom(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.Room = "bridge"
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
		shown  string
		hidden string
	}{
		{name: "player room", shown: "B", hidden: "C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
