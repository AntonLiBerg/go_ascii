package render

import (
	component "go_ascii/internal"
	"go_ascii/internal/world"
	"strconv"
	"strings"
)

func TerminalFrame(gameWorld world.World) string {
	var frame strings.Builder
	frame.WriteString("\033[2J\033[H")

	nextY := 0
	if len(gameWorld.UILayout) == 0 {
		nextY = drawRoom(&frame, gameWorld, nextY)
	} else {
		for _, name := range gameWorld.UILayout {
			if name == "room" {
				nextY = drawRoom(&frame, gameWorld, nextY)
				continue
			}
			nextY = drawUI(&frame, gameWorld.UIContent[name], nextY)
		}
	}

	writeCursor(&frame, nextY+1, 1)
	return frame.String()
}

func drawRoom(frame *strings.Builder, gameWorld world.World, yOffset int) int {
	activeRoom := gameWorld.ViewRoom
	hasActiveRoom := activeRoom != ""
	if !hasActiveRoom {
		for _, playerID := range world.GetPlayerIDs(gameWorld) {
			if pos, ok := gameWorld.Pos[playerID]; ok {
				activeRoom = pos.Room
				hasActiveRoom = true
				break
			}
		}
	}

	nextY := yOffset
	for _, eID := range gameWorld.Entities {
		pos, okPos := gameWorld.Pos[eID]
		ascii, okASCII := gameWorld.Ascii[eID]
		if !okPos || !okASCII || (hasActiveRoom && pos.Room != activeRoom) || pos.X < 0 || pos.Y < 0 || isCovered(gameWorld, eID, pos) {
			continue
		}

		screenY := pos.Y + yOffset + 1
		writeCursor(frame, screenY, pos.X+1)
		frame.WriteRune(ascii.Ascii)
		if screenY > nextY {
			nextY = screenY
		}
	}
	return nextY
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
