package worldmanager

import "github.com/df-mc/dragonfly/server/world"

type World struct {
	*world.World
	name string
}

func NewWorld(w *world.World, name string) *World {
	return &World{World: w, name: name}
}

func (w *World) Name() string { return w.name }
