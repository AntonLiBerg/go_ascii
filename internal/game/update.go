package game

import "go_ascii/internal/world"

type UpdateFunc struct {
	Order      int
	UpdateFunc func(world.World) (world.World, error)
}
type IService interface {
	GetUpdateFunc(world world.World) UpdateFunc
}
