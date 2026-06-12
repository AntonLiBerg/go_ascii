package service

import (
	wrld "go_ascii/world"
)

type IService interface {
	Update(world wrld.World) (func(world *wrld.World), error)
}

type ServiceDrawOnTerminal struct{}
