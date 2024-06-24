package assets

import (
	_ "embed"
)

var (
	//go:embed EmpireStateNF.ttf
	EmpireStateNF_ttf []byte

	//go:embed background.png
	Background_png []byte

	//go:embed paddle.png
	Paddle_png []byte

	//go:embed ball.png
	Ball_png []byte
)
