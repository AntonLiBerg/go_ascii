package main

import (
	"os"
	"strings"
)

func main() {
	res, err := os.ReadFile("./scenarios/demo/map.txt")
	if err != nil {
		panic(err)
	}
	world := World{}

	content := string(res)
	splitContent := strings.Split(content, "===")
	scMap := []rune(mastrings.Split(splitContent[1], "MAP")[1])
	scEntities := []rune(strings.Split(splitContent[2], "ENTITY")[1])

	rowNum := 0
	for i := 0; i < len(scMap); i++ {
		world.AddEntity()
		eId := world.nextEnt - 1
		world.Ascii[eId] = cAscii{Ascii: scMap[i]}
		world.Pos[eId] = cPosition{X:}
	}
}

type World struct {
	nextEnt  int
	Entities []int
	Pos      map[int]cPosition
	Ascii    map[int]cAscii
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
