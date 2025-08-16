package setup

import (
	"fmt"

	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type CrystalPVPConfig struct {
	ArenaPos1 cube.Pos `yaml:"arena_pos_1"`
	ArenaPos2 cube.Pos `yaml:"arena_pos_2"`
	Spawns    []cube.Pos
}

func NewCrystalPVPSetup(name string) ISetup {
	var v CrystalPVPConfig
	v.Spawns = make([]cube.Pos, 8)
	return &Setup[CrystalPVPConfig]{
		Name:  name,
		Steps: CrystalPVPSteps(),
		Value: &v,
	}
}

func CrystalPVPSteps() []step.Step[CrystalPVPConfig] {
	var steps []step.Step[CrystalPVPConfig]
	steps = append(steps, step.New("Arena Pos 1", func(t *CrystalPVPConfig, pos cube.Pos) {
		t.ArenaPos1 = pos
	}))
	steps = append(steps, step.New("Arena Pos 2", func(t *CrystalPVPConfig, pos cube.Pos) {
		t.ArenaPos2 = pos
	}))

	for i := 0; i < 8; i++ {
		steps = append(steps, step.New(fmt.Sprintf("Spawn #%d", i+1), func(t *CrystalPVPConfig, pos cube.Pos) {
			pos = pos.Add(cube.Pos{0, 1})
			t.Spawns[i] = pos
		}))
	}

	return steps
}
