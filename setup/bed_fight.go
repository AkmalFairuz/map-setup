package setup

import (
	"strings"

	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type BedFightConfig struct {
	HeightLimit int                  `yaml:"height_limit"`
	Mid         cube.Pos             `yaml:"mid"`
	Spawn       BedWarsTeamPosConfig `yaml:"spawn"`
	Bed         BedWarsTeamPosConfig `yaml:"bed"`
}

func NewBedFightSetup(name string) ISetup {
	var v BedFightConfig
	v.Spawn = make(BedWarsTeamPosConfig)
	v.Bed = make(BedWarsTeamPosConfig)
	return &Setup[BedFightConfig]{
		Name:  name,
		Steps: BedFightSteps(),
		Value: &v,
	}
}

func BedFightSteps() []step.Step[BedFightConfig] {
	steps := []step.Step[BedFightConfig]{
		step.New("Mid Pos", func(t *BedFightConfig, pos cube.Pos) {
			t.Mid = pos
		}),
	}
	teams := []string{"red", "blue"}

	for _, team := range teams {
		teamName := strings.Title(team)
		steps = append(steps, step.New(teamName+" Spawn Pos", func(t *BedFightConfig, pos cube.Pos) {
			pos[1] += 1
			t.Spawn[team] = pos
		}))
		steps = append(steps, step.New(teamName+" Bed Pos (the pillow)", func(t *BedFightConfig, pos cube.Pos) {
			t.Bed[team] = pos
		}))
	}

	return steps
}
