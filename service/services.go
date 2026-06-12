package service

import (
	wrld "go_ascii/world"
	"os"
)

type IService interface {
	Update(world *wrld.World) error
}

type ServiceUserInput struct{}

func (ServiceUserInput) Update(world *wrld.World) error {
	var key [1]byte
	_, err := os.Stdin.Read(key[:])
	if err != nil {
		return err
	}

	world.ClearUserInput()
	world.SetKeyDown(string(key[:]))

	return nil
}
