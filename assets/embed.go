package assets

import (
	_ "embed"
)

var (
	//go:embed EmpireStateNF.ttf
	EmpireStateNF_ttf []byte

	//go:embed background.png
	Background []byte
)
