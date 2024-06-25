package pong

import "github.com/hajimehoshi/ebiten/v2"

type paddle struct {
	x, y   int
	dy     int
	score  int
	sprite *ebiten.Image
}
