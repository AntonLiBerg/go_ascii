package service

import (
	wrld "go_ascii/world"
)

type UpdateFunc struct {
	Order      int
	UpdateFunc func(*wrld.World)
	Err        error
}
type IService interface {
	GetUpdateFunc(world wrld.World) UpdateFunc
}
