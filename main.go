package main

import (
	"os"
	"strings"
	"go_ascii/aiAPI"
)

func main() {
	world := NewWorld()
	ai := aiapi.New()
	aMap,aEnt := ai.GetAsciiMapAndEntitiesFromFile("./scenarios/demo/map.txt")

	for pos,ch := range aMap{
		world.AddEntity()
		world.Pos
	}

}

type World struct {
	nextEnt  int
	Entities []int
	Pos      map[int]cPosition
	Ascii    map[int]cAscii
}

func NewWorld() World {
	return World{
		nextEnt:0
		Entities: []int{},
		Pos:      map[int]cPosition{},
		Ascii:    map[int]cAscii{},
	}
}
func (w *World) AddEntity() {
	w.Entities = append(w.Entities, w.nextEnt)
	w.nextEnt++
}

type cPosition struct {
	X int
	Y int
}
type cAscii struct {
	Ascii rune
}
