package setup

import (
	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
)

type TheBridgeConfig struct {
	HeightLimit  int                 `yaml:"height_limit"`
	Void         int                 `yaml:"void"`
	Goals        map[string]cube.Pos `yaml:"goals"`
	Spawns       map[string]cube.Pos `yaml:"spawns"`
	WallStartPos cube.Pos            `yaml:"wall_start_pos"`
	WallEndPos   cube.Pos            `yaml:"wall_end_pos"`
}

func NewTheBridgeSetup(name string) ISetup {
	var v TheBridgeConfig
	v.Goals = make(map[string]cube.Pos)
	v.Spawns = make(map[string]cube.Pos)
	return &Setup[TheBridgeConfig]{
		Name:  name,
		Steps: TheBridgeSteps(),
		Value: &v,
	}
}

func TheBridgeSteps() []step.Step[TheBridgeConfig] {
	var steps []step.Step[TheBridgeConfig]
	teams := []string{"red", "blue"}
	for _, t := range teams {
		steps = append(steps, step.New("Cage Positions for "+t, func(c *TheBridgeConfig, pos cube.Pos) {
			pos[1] += 1
			c.Spawns[t] = pos
		}))
		steps = append(steps, step.New("Goal Position for "+t, func(c *TheBridgeConfig, pos cube.Pos) {
			c.Goals[t] = pos
		}))
	}
	steps = append(steps, step.New("Wall Start Position (pos1, from upper red team side)", func(c *TheBridgeConfig, pos cube.Pos) {
		c.WallStartPos = pos
	}))
	steps = append(steps, step.New("Wall End Position (pos2, from lower blue team side)", func(c *TheBridgeConfig, pos cube.Pos) {
		c.WallEndPos = pos
	}))
	//TODO: height limit
	return steps
}
