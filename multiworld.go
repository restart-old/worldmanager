package worldmanager

import (
	"fmt"
	"sync"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/sirupsen/logrus"
)

type WorldManager struct {
	server *server.Server
	logger *logrus.Logger

	worlds   map[string]*world.World
	worldsMu sync.RWMutex
}

func New(server *server.Server, logger *logrus.Logger) *WorldManager {
	defaultWorld := server.World()
	return &WorldManager{
		logger: logger,
		worlds: map[string]*world.World{
			defaultWorld.Name(): defaultWorld,
		},
	}
}
func (mw *WorldManager) Worlds() []*world.World {
	mw.worldsMu.RLock()
	worlds := make([]*world.World, 0, len(mw.worlds))
	for _, w := range mw.worlds {
		worlds = append(worlds, w)
	}
	mw.worldsMu.RUnlock()
	return worlds
}
func (mw *WorldManager) World(name string) (*world.World, bool) {
	mw.worldsMu.RLock()
	w, ok := mw.worlds[name]
	mw.worldsMu.RUnlock()
	return w, ok
}

func (mw *WorldManager) DefaultWorld() *world.World { return mw.server.World() }

func (mw *WorldManager) LoadWorld(worldPath string, settings *world.Settings, dimension world.Dimension) error {
	if _, ok := mw.World(settings.Name); ok {
		return fmt.Errorf("world manager: world with name '%s' is already loaded", settings.Name)
	}
	w := world.New(mw.logger, dimension, settings)
	p, err := mcdb.New(worldPath, dimension)
	if err != nil {
		return fmt.Errorf("world manager: %s", err)
	}
	w.Provider(p)

	mw.worldsMu.Lock()
	defer mw.worldsMu.Unlock()
	mw.worlds[settings.Name] = w
	return nil
}

func (mw *WorldManager) Close() error {
	mw.worldsMu.Lock()
	for _, w := range mw.worlds {
		if w != mw.DefaultWorld() {
			mw.logger.Debugf("Closing world '%s'\n", w.Name())
			if err := w.Close(); err != nil {
				return err
			}
		}
	}
	mw.worlds = nil
	mw.worldsMu.Unlock()
	return nil
}