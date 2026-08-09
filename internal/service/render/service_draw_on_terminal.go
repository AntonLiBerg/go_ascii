package render

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ServiceDrawOnTerminal struct {
	Writer io.Writer
}

func (s ServiceDrawOnTerminal) GetUpdateFunc(_ world.World) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 100,
		UpdateFunc: func(w world.World) (world.World, error) {
			if w.IterationNr != 1 && !w.HasChanged {
				return w, nil
			}

			writer := s.Writer
			if writer == nil {
				writer = os.Stdout
			}
			if _, err := fmt.Fprint(writer, TerminalFrame(w)); err != nil {
				return w, err
			}
			w.HasChanged = false
			return w, nil
		},
	}
}

func TerminalFrame(gameWorld world.World) string {
	var frame strings.Builder
	frame.WriteString("\033[2J\033[H")

	nextY := 0
	roomXOffset := roomHorizontalOffset(gameWorld, gameWorld.Room, uiWidth(gameWorld))
	if len(gameWorld.UILayout) == 0 {
		nextY = drawRoom(&frame, gameWorld, gameWorld.Room, nextY, 0)
	} else {
		// Layout order controls the vertical composition of the frame.
		for _, name := range gameWorld.UILayout {
			if name == "room" {
				nextY = drawRoom(&frame, gameWorld, gameWorld.Room, nextY, roomXOffset)
				continue
			}
			nextY = drawUI(&frame, gameWorld.UIContent[name], nextY)
		}
	}

	writeCursor(&frame, nextY+1, 1)
	return frame.String()
}

func drawRoom(frame *strings.Builder, gameWorld world.World, activeRoom string, yOffset int, xOffset int) int {
	nextY := yOffset
	for _, eID := range gameWorld.Entities {
		pos, okPos := gameWorld.Pos[eID]
		ascii, okASCII := gameWorld.Ascii[eID]
		if !okPos || !okASCII || (activeRoom != "" && pos.Room != activeRoom) || pos.X < 0 || pos.Y < 0 || isCovered(gameWorld, eID, pos) {
			continue
		}

		screenY := pos.Y + yOffset + 1
		writeCursor(frame, screenY, pos.X+xOffset+1)
		char := ascii.Ascii
		if gameWorld.EditingControl && gameWorld.ActiveControl != nil && gameWorld.ActiveControl.SelectableEntityID == eID {
			char = gameWorld.ActiveControl.SelectedASCII
		}
		frame.WriteRune(char)
		if screenY > nextY {
			nextY = screenY
		}
	}
	return nextY
}

func uiWidth(gameWorld world.World) int {
	width := 0
	for _, name := range gameWorld.UILayout {
		if name == "room" {
			continue
		}
		for _, line := range gameWorld.UIContent[name] {
			if lineWidth := utf8.RuneCountInString(line); lineWidth > width {
				width = lineWidth
			}
		}
	}
	return width
}

func roomWidth(gameWorld world.World, activeRoom string) int {
	width := 0
	for _, eID := range gameWorld.Entities {
		pos, okPos := gameWorld.Pos[eID]
		_, okASCII := gameWorld.Ascii[eID]
		if !okPos || !okASCII || (activeRoom != "" && pos.Room != activeRoom) || pos.X < 0 {
			continue
		}
		if pos.X+1 > width {
			width = pos.X + 1
		}
	}
	return width
}

func roomHorizontalOffset(gameWorld world.World, activeRoom string, width int) int {
	roomWidth := roomWidth(gameWorld, activeRoom)
	if roomWidth == 0 || roomWidth >= width {
		return 0
	}
	return (width - roomWidth) / 2
}

func drawUI(frame *strings.Builder, lines []string, yOffset int) int {
	for _, line := range lines {
		writeCursor(frame, yOffset+1, 1)
		frame.WriteString(line)
		yOffset++
	}
	return yOffset
}

func writeCursor(frame *strings.Builder, y, x int) {
	frame.WriteString("\033[")
	frame.WriteString(strconv.Itoa(y))
	frame.WriteByte(';')
	frame.WriteString(strconv.Itoa(x))
	frame.WriteByte('H')
}

func isCovered(world world.World, eID int, pos component.Position) bool {
	layer := world.Layer[eID].Nr
	_, isPlayer := world.Player[eID]
	for _, otherID := range world.Entities {
		if otherID == eID || world.Pos[otherID] != pos {
			continue
		}
		if _, hasASCII := world.Ascii[otherID]; !hasASCII {
			continue
		}

		otherLayer := world.Layer[otherID].Nr
		_, otherIsPlayer := world.Player[otherID]
		if otherLayer > layer || (otherLayer == layer && otherIsPlayer && !isPlayer) {
			return true
		}
		if otherLayer == layer && otherIsPlayer == isPlayer && otherID > eID {
			return true
		}
	}
	return false
}
