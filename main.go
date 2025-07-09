package main

import (
	"fmt"
	"github.com/akmalfairuz/map-setup/setup"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/lmittmann/tint"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log/slog"
	"os"
	"time"
)

func main() {
	log := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.DateTime,
	}))
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(log)
	if err != nil {
		panic(err)
	}
	conf.RandomTickSpeed = -1
	conf.PlayerProvider = player.NopProvider{}
	conf.Generator = func(dim world.Dimension) world.Generator {
		return world.NopGenerator{}
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()
	srv.World().StopTime()
	srv.World().SetTime(3000)
	srv.World().StopWeatherCycle()
	srv.World().StopRaining()
	srv.World().StopThundering()

	srv.Listen()
	for p := range srv.Accept() {
		_ = p
		p.Message(text.Colourf("<grey>welcome to <green>venity's map-setup</green> server! type <green>start [game] [name]</green> to start configure a map.</grey>"))
		p.SetGameMode(world.GameModeCreative)
		p.ShowCoordinates()
		p.Inventory().Clear()
		setup.SendSetupItems(p)
		p.Handle(setup.NewPlayerHandler(log))
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
