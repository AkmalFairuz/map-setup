package setup

import (
	"strconv"

	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type PresentHuntConfig struct {
	Pos []*YamlNode `yaml:"pos"`
}

func NewPresentHuntSetup(name string) ISetup {
	var v PresentHuntConfig
	return &Setup[PresentHuntConfig]{
		Name:  name,
		Steps: PresentHuntSteps(),
		Value: &v,
	}
}

func PresentHuntSteps() []step.Step[PresentHuntConfig] {
	var steps []step.Step[PresentHuntConfig]

	for i := 1; i <= 20; i++ {
		label := "Present Position-" + strconv.Itoa(i)
		steps = append(steps, step.New(
			label,
			func(c *PresentHuntConfig, pos cube.Pos) {
				pos[1] += 1
				c.Pos = append(c.Pos, posToNode(pos))
			},
		))
	}
	return steps
}
