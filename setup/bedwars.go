package setup

import (
	"fmt"
	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
	"strings"
)

type BedWarsTeamPosConfig map[string]cube.Pos

type BedWarsConfig struct {
	HeightLimit int      `yaml:"height_limit"`
	Mid         cube.Pos `yaml:"mid"`
	Generator   struct {
		Team    BedWarsTeamPosConfig `yaml:"team"`
		Diamond []cube.Pos           `yaml:"diamond"`
		Emerald []cube.Pos           `yaml:"emerald"`
	} `yaml:"generator"`
	Spawn       BedWarsTeamPosConfig `yaml:"spawn"`
	Bed         BedWarsTeamPosConfig `yaml:"bed"`
	ItemShop    BedWarsTeamPosConfig `yaml:"item_shop"`
	UpgradeShop BedWarsTeamPosConfig `yaml:"upgrade_shop"`
}

func NewBedWarsSetup(name string, isSquads bool) ISetup {
	var v BedWarsConfig
	v.Generator.Team = make(BedWarsTeamPosConfig)
	v.Spawn = make(BedWarsTeamPosConfig)
	v.Bed = make(BedWarsTeamPosConfig)
	v.ItemShop = make(BedWarsTeamPosConfig)
	v.UpgradeShop = make(BedWarsTeamPosConfig)
	return &Setup[BedWarsConfig]{
		Name:  name,
		Steps: BedWarsSteps(isSquads),
		Value: &v,
	}
}

func BedWarsSteps(isSquads bool) []step.Step[BedWarsConfig] {
	steps := []step.Step[BedWarsConfig]{
		step.New("Mid Pos", func(t *BedWarsConfig, pos cube.Pos) {
			t.Mid = pos
		}),
	}

	//steps = append(steps, step.New("Height Limit", func(t *BedWarsConfig, pos cube.Pos) {
	//	t.HeightLimit = pos.Y()
	//}))

	teams := []string{"red", "blue", "green", "yellow"}
	if !isSquads {
		teams = append(teams, "aqua", "white", "pink", "gray")
	}
	for _, team := range teams {
		teamName := strings.Title(team)
		steps = append(steps, step.New(teamName+" Spawn Pos", func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.Spawn[team] = pos
		}))
		steps = append(steps, step.New(teamName+" Generator Pos", func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.Generator.Team[team] = pos
		}))
		steps = append(steps, step.New(teamName+" Item Shop Pos", func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.ItemShop[team] = pos
		}))
		steps = append(steps, step.New(teamName+" Upgrade Shop Pos", func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.UpgradeShop[team] = pos
		}))
		steps = append(steps, step.New(teamName+" Bed Pos (the pillow)", func(t *BedWarsConfig, pos cube.Pos) {
			t.Bed[team] = pos
		}))
	}

	maxDiamondGenerators := 4
	for i := 1; i <= maxDiamondGenerators; i++ {
		steps = append(steps, step.New(fmt.Sprintf("Diamond Generator #%d", i), func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.Generator.Diamond = append(t.Generator.Diamond[:i-1], pos)
		}))
	}

	maxEmeraldGenerators := 4
	if isSquads {
		maxEmeraldGenerators = 2
	}
	for i := 1; i <= maxEmeraldGenerators; i++ {
		steps = append(steps, step.New(fmt.Sprintf("Emerald Generator #%d", i), func(t *BedWarsConfig, pos cube.Pos) {
			pos[1] += 1
			t.Generator.Emerald = append(t.Generator.Emerald[:i-1], pos)
		}))
	}

	return steps
}
