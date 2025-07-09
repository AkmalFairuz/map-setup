package step

import "github.com/df-mc/dragonfly/server/block/cube"

type Step[T any] struct {
	Description string
	HandlePos   func(*T, cube.Pos)
}

func New[T any](description string, handlePos func(*T, cube.Pos)) Step[T] {
	return Step[T]{
		Description: description,
		HandlePos:   handlePos,
	}
}
