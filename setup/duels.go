package setup

import (
	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type DuelsConfig struct {
	Mid        *YamlNode   `yaml:"mid"`
	Spawns     []*YamlNode `yaml:"spawns"`
	SumoCenter *YamlNode   `yaml:"sumo_mid"`
	SumoSpawns []*YamlNode `yaml:"sumo_spawns"`
}

func NewDuelsSetup(name string) ISetup {
	var v DuelsConfig
	return &Setup[DuelsConfig]{
		Name:  name,
		Steps: DuelsSteps(),
		Value: &v,
	}
}

func DuelsSteps() []step.Step[DuelsConfig] {
	var steps []step.Step[DuelsConfig]

	for i := 1; i <= 2; i++ {
		index := i
		steps = append(steps, step.New(
			"Player Position-"+string(rune('0'+index)),
			func(c *DuelsConfig, pos cube.Pos) {
				pos[1] += 1
				c.Spawns = append(c.Spawns, posToNode(pos))
			},
		))
	}

	steps = append(steps, step.New("Set Center Position (Normal)", func(c *DuelsConfig, pos cube.Pos) {
		c.Mid = posToNode(pos)
	}))

	for i := 1; i <= 2; i++ {
		index := i
		steps = append(steps, step.New(
			"Sumo Player Position-"+string(rune('0'+index)),
			func(c *DuelsConfig, pos cube.Pos) {
				pos[1] += 1
				c.SumoSpawns = append(c.SumoSpawns, posToNode(pos))
			},
		))
	}

	steps = append(steps, step.New("Set Center Position (Sumo)", func(c *DuelsConfig, pos cube.Pos) {
		c.SumoCenter = posToNode(pos)
	}))

	return steps
}
