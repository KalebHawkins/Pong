package pong

import "github.com/hajimehoshi/ebiten/v2"

type ball struct {
	x, y   int
	dx, dy int
	prevDx int
	prevDy int
	sprite *ebiten.Image
}
