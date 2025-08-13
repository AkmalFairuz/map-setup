package setup

import (
	"strconv"

	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type SkywarsConfig struct {
	Mid    *YamlNode   `yaml:"mid"`
	Spawns []*YamlNode `yaml:"spawns"`
}

func NewSkywarsSetup(name string) ISetup {
	var v SkywarsConfig
	return &Setup[SkywarsConfig]{
		Name:  name,
		Steps: SkywarsSteps(),
		Value: &v,
	}
}

func SkywarsSteps() []step.Step[SkywarsConfig] {
	var steps []step.Step[SkywarsConfig]

	for i := 1; i <= 12; i++ {
		label := "Player Cage Position-" + strconv.Itoa(i)
		steps = append(steps, step.New(
			label,
			func(c *SkywarsConfig, pos cube.Pos) {
				pos[1] += 1
				c.Spawns = append(c.Spawns, posToNode(pos))
			},
		))
	}

	steps = append(steps, step.New("Set Mid Position", func(c *SkywarsConfig, pos cube.Pos) {
		c.Mid = posToNode(pos)
	}))
	return steps
}
