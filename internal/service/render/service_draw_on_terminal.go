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

	maxY := 0
	for _, eID := range gameWorld.Entities {
		pos, okPos := gameWorld.Pos[eID]
		ascii, okASCII := gameWorld.Ascii[eID]
		if !okPos || !okASCII || (hasActiveRoom && pos.Room != activeRoom) || pos.X < 0 || pos.Y < 0 || isCovered(gameWorld, eID, pos) {
			continue
		}

		frame.WriteString("\033[")
		frame.WriteString(strconv.Itoa(pos.Y + 1))
		frame.WriteByte(';')
		frame.WriteString(strconv.Itoa(pos.X + 1))
		frame.WriteByte('H')
		frame.WriteRune(ascii.Ascii)
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	frame.WriteString("\033[")
	frame.WriteString(strconv.Itoa(maxY + 2))
	frame.WriteString(";1H")
	return frame.String()
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
