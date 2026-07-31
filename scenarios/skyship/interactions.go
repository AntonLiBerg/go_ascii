package skyship

import (
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/world"
)

func InteractionHandlers() interaction.Handlers {
	return interaction.Handlers{
		Helm: interactWithHelm,
	}
}

func interactWithHelm(w *world.World, entityID int) {
	//if they interact, show the helm scene
}
