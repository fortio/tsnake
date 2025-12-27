//nolint:gosec //we don't need cryptographically secure randomness
package main

import (
	"math/rand/v2"
)

type coords struct{ X, Y int }

type direction = int32

const (
	U direction = iota
	D
	L
	R
)

var directionCoords = [4]coords{
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
}

type snake struct {
	m              map[coords]bool
	food           coords
	snake          []coords // last element is snake mouth
	maxX, maxY     int
	dir            direction
	verticalStreak int
	square         bool
	firstFrame     bool
	paused         bool
	circular       bool
}

func newSnake(mx, my int, square bool, circular bool) *snake {
	snak := make([]coords, 0, mx*my)
	snak = append(snak, coords{mx / 2, my / 2})
	m := make(map[coords]bool)
	m[coords{mx / 2, my / 2}] = true
	d := U
	s := snake{
		snake: snak, food: coords{
			rand.IntN(mx),
			rand.IntN(my),
		},
		maxX: mx,
		maxY: my, dir: d, m: m,
		square:   square,
		circular: circular,
	}
	return &s
}

func (s *snake) next() bool {
	if s.paused {
		return true
	}
	if (s.dir == U || s.dir == D) && !s.square {
		s.verticalStreak++
		if s.verticalStreak%2 != 0 {
			return true
		}
	}
	s.firstFrame = false
	length := len(s.snake)
	dir := s.dir
	changed := directionCoords[dir]
	x := (s.snake[length-1].X + changed.X)
	y := (s.snake[length-1].Y + changed.Y)
	if x >= s.maxX || y >= s.maxY || x < 0 || y < 0 && !s.circular {
		return false
	}
	next := coords{
		X: (s.snake[length-1].X + changed.X) % s.maxX,
		Y: (s.snake[length-1].Y + changed.Y) % s.maxY,
	}
	if next.X < 0 {
		next.X += s.maxX
	}
	if next.Y < 0 {
		next.Y += s.maxY
	}
	s.snake = append(s.snake, next)
	s.m[next] = true
	if s.food == next {
		s.food = coords{rand.IntN(s.maxX), rand.IntN(s.maxY)}
		return len(s.snake) == len(s.m)
	}
	delete(s.m, s.snake[0])
	s.snake = s.snake[1:]
	return len(s.snake) == len(s.m)
}
