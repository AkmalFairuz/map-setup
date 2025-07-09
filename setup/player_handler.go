package setup

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log/slog"
	"strings"
	"time"
)

type PlayerHandler struct {
	log *slog.Logger

	player.NopHandler

	setup ISetup

	lastExec time.Time
}

func NewPlayerHandler(log *slog.Logger) *PlayerHandler {
	return &PlayerHandler{
		log:   log,
		setup: NopSetup{},
	}
}

func (ph *PlayerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	ctx.Cancel()
	ctx.Val().Tx().SetBlock(pos, b, &world.SetOpts{DisableBlockUpdates: true})
	ctx.Val().Tx().PlaySound(pos.Vec3Centre(), sound.BlockPlace{Block: b})
}

func (ph *PlayerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if mainHandVal(ctx.Val()) == "break" {
		ctx.Cancel()
		if time.Since(ph.lastExec) < 500*time.Millisecond {
			ctx.Val().Message(text.Red + "You are breaking blocks too fast! TRY AGAIN.")
			return
		}
		ph.lastExec = time.Now()
		ph.setup.Execute(pos)
		ph.next(ctx.Val())
		return
	}

	ctx.Cancel()

	ctx.Val().Tx().AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: ctx.Val().Tx().Block(pos)})
	ctx.Val().Tx().SetBlock(pos, block.Air{}, &world.SetOpts{DisableBlockUpdates: true})
}

func (ph *PlayerHandler) HandleItemUse(ctx *player.Context) {
	if mainHandVal(ctx.Val()) == "forward" {
		ctx.Cancel()
		mul := 50.0
		if ctx.Val().OnGround() {
			mul = 10.0
		}
		if ctx.Val().Sprinting() {
			mul *= 0.7
		}
		ctx.Val().SetVelocity(ctx.Val().Rotation().Vec3().Mul(mul))
	}
}

func (ph *PlayerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	ctx.Cancel()
	if mainHandVal(ctx.Val()) == "block" {
		if time.Since(ph.lastExec) < 500*time.Millisecond {
			ctx.Val().Message(text.Red + "You are interacting blocks too fast! TRY AGAIN.")
			return
		}
		ph.lastExec = time.Now()
		ph.setup.Execute(pos)
		ph.next(ctx.Val())
	}
}

func (ph *PlayerHandler) HandleChat(ctx *player.Context, message *string) {
	ctx.Cancel()
	parts := strings.Split(*message, " ")
	if len(parts) < 1 {
		return
	}
	switch parts[0] {
	case "start":
		if len(parts) < 3 {
			ctx.Val().Message(text.Colourf("<red>Usage: start [game] [name]</red>"))
			return
		}
		game := parts[1]
		name := strings.Join(parts[2:], " ")
		switch game {
		case "bw4":
			ph.setup = NewBedWarsSetup(name, true)
		case "bw8":
			ph.setup = NewBedWarsSetup(name, false)
		default:
			ctx.Val().Message(text.Colourf("<red>Unknown game: %s, available: <green>bw4, bw8</green></red>", game))
			return
		}
		ph.log.Info("setup started", "game", game, "name", name)
		ph.setup.SetLog(ph.log)
		ctx.Val().Message(text.Colourf("<green>Setup started for <yellow>%s</yellow>: <yellow>%s</yellow>. Please follow the instructions!</green>", game, name))
		ctx.Val().Message(text.Colourf("<red>in case you want back, use <yellow>prev</yellow> command</red>"))
		ctx.Val().Message(text.Colourf("<aqua>Initial Instruction: <grey>%s</grey>.</aqua>", ph.setup.CurrentDescription()))
	case "next":
		ph.next(ctx.Val())
	case "prev":
		ph.prev(ctx.Val())
	case "tp":
		if len(parts) < 4 {
			ctx.Val().Messagef("Usage: tp <x> <y> <z>")
			return
		}
		var pos cube.Pos
		if _, err := fmt.Sscanf(parts[0]+" "+parts[1]+" "+parts[2], "%d %d %d", &pos[0], &pos[1], &pos[2]); err != nil {
			ctx.Val().Messagef("invalid coordinates")
			return
		}
		ctx.Val().Teleport(pos.Vec3Centre())
	case "exit":
		ctx.Val().Message(text.Colourf("<red>Exiting setup mode.</red>"))
		ph.setup = NopSetup{}
	default:
		ctx.Val().Message(text.Colourf("<red>Unknown command: %s, available: <green>start, next, prev, tp, exit</green></red>", parts[0]))
	}
}

func (ph *PlayerHandler) next(p *player.Player) {
	if !ph.setup.HasNext() {
		p.Message(text.Colourf("<red>No more steps available. Finished? Type <green>exit</green>.</red>"))
		return
	}
	ph.setup.Next()
	p.Message(text.Colourf("<green>Next Instruction: <grey>%s</grey>.</green>", ph.setup.CurrentDescription()))
}

func (ph *PlayerHandler) prev(p *player.Player) {
	if !ph.setup.HasPrev() {
		p.Message(text.Colourf("<red>No previous steps available.</red>"))
		return
	}
	ph.setup.Prev()
	p.Message(text.Colourf("<green>Previous Instruction: <grey>%s</grey>.</green>", ph.setup.CurrentDescription()))
}

func (ph *PlayerHandler) HandleItemDrop(ctx *player.Context, s item.Stack) {
	ctx.Cancel()
}

func SendSetupItems(p *player.Player) {
	_ = p.Inventory().SetItem(0, item.NewStack(item.Axe{Tier: item.ToolTierGold}, 1).WithValue("setup", "break").WithCustomName("[BREAK]"))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Axe{Tier: item.ToolTierStone}, 1).WithValue("setup", "block").WithCustomName("[BLOCK]"))
	_ = p.Inventory().SetItem(2, item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithValue("setup", "forward").WithCustomName("PUSH YOUR CHARACTER TO FRONT"))
	_ = p.Inventory().SetItem(3, item.NewStack(block.InvisibleBedrock{}, 1))
}

func mainHandVal(p *player.Player) string {
	v, _ := p.HeldItems()
	x, ok := v.Value("setup")
	if ok {
		return x.(string)
	}
	return ""
}
