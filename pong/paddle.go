package pong

import "github.com/hajimehoshi/ebiten/v2"

type paddle struct {
	x, y   int
	dx     int
	score  int
	sprite *ebiten.Image
}
