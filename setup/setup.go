package setup

import (
	"fmt"
	"github.com/akmalfairuz/map-setup/setup/step"
	"github.com/df-mc/dragonfly/server/block/cube"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

type ISetup interface {
	Execute(pos cube.Pos)
	HasNext() bool
	HasPrev() bool
	Next()
	Prev()
	CurrentDescription() string
	SetLog(log *slog.Logger)
}

type Setup[T any] struct {
	Log       *slog.Logger
	Name      string
	StepIndex int
	Steps     []step.Step[T]
	Value     *T
}

func (s *Setup[T]) SetLog(log *slog.Logger) {
	s.Log = log
	if s.Log == nil {
		s.Log = slog.Default()
	}
}

func (s *Setup[T]) Execute(pos cube.Pos) {
	if s.StepIndex < 0 || s.StepIndex >= len(s.Steps) {
		return
	}
	s.Steps[s.StepIndex].HandlePos(s.Value, pos)
	s.Save()
	s.Log.Info("executed step", "description", s.Steps[s.StepIndex].Description, "x", pos.X(), "y", pos.Y(), "z", pos.Z(), "name", s.Name)
}

func (s *Setup[T]) HasNext() bool {
	return s.StepIndex < len(s.Steps)-1
}

func (s *Setup[T]) Next() {
	if s.StepIndex < len(s.Steps)-1 {
		s.StepIndex++
	}
}

func (s *Setup[T]) HasPrev() bool {
	return s.StepIndex > 0
}

func (s *Setup[T]) Prev() {
	if s.StepIndex > 0 {
		s.StepIndex--
	}
}

func (s *Setup[T]) Save() {
	bytes, err := yaml.Marshal(s.Value)
	if err != nil {
		fmt.Printf("error marshalling setup value: %v\n", err)
	}
	if err := os.WriteFile(s.Name+"_config.yml", bytes, 0644); err != nil {
		fmt.Printf("error writing setup config file: %v\n", err)
	}
}

func (s *Setup[T]) CurrentDescription() string {
	if s.StepIndex < 0 || s.StepIndex >= len(s.Steps) {
		return ""
	}
	return s.Steps[s.StepIndex].Description
}

type NopSetup struct{}

func (NopSetup) Execute(pos cube.Pos)       {}
func (NopSetup) HasNext() bool              { return false }
func (NopSetup) HasPrev() bool              { return false }
func (NopSetup) Next()                      {}
func (NopSetup) Prev()                      {}
func (NopSetup) CurrentDescription() string { return "" }
func (NopSetup) SetLog(log *slog.Logger)    {}
