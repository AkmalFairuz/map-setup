package setup

import (
	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type DuelsConfig struct {
	HeightLimit int                 `yaml:"height_limit"`
	Void        int                 `yaml:"void"`
	Center      cube.Pos            `yaml:"center"`
	Spawns      map[string]cube.Pos `yaml:"spawns"`
	SumoSpawns  map[string]cube.Pos `yaml:"sumo_spawns"`
	SumoCenter  cube.Pos            `yaml:"sumo_center"`
}

func NewDuelsSetup(name string) ISetup {
	var v DuelsConfig
	v.Spawns = make(map[string]cube.Pos)
	v.SumoSpawns = make(map[string]cube.Pos)
	return &Setup[DuelsConfig]{
		Name:  name,
		Steps: DuelsSteps(),
		Value: &v,
	}
}

func DuelsSteps() []step.Step[DuelsConfig] {
	var steps []step.Step[DuelsConfig]
	teams := []string{"red", "blue"}
	for _, t := range teams {
		steps = append(steps, step.New("Player's Position for "+t, func(c *DuelsConfig, pos cube.Pos) {
			pos[1] += 1
			c.Spawns[t] = pos
		}))
	}
	steps = append(steps, step.New("Set Center Position (Normal)", func(c *DuelsConfig, pos cube.Pos) {
		c.Center = pos
	}))
	for _, t := range teams {
		steps = append(steps, step.New("Player's Sumo Position for "+t, func(c *DuelsConfig, pos cube.Pos) {
			pos[1] += 1
			c.SumoSpawns[t] = pos
		}))
	}
	steps = append(steps, step.New("Set Center Position (Sumo)", func(c *DuelsConfig, pos cube.Pos) {
		c.SumoCenter = pos
	}))
	//TODO: height limit and sumo void step
	return steps
}
