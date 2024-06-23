package pong

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gtihub.com/KalebHawkins/pong/assets"
)

const (
	fontSize      = 48
	titleFontSize = fontSize * 1.5
)

// Game is a structure containing the game data and configuration.
type Game struct {
	*Cfg
	state
}

// Cfg contains the Game's configuration data.
type Cfg struct {
	// ScreenWidth represents the width of the window
	ScreenWidth int
	// ScreenHeight represents the height of the window
	ScreenHeight int
	// WindowTitle is the title displayed in the window's title bar
	WindowTitle string
	// faceSource is the font face used for the menu, scoreboard, etc
	faceSource *text.GoTextFaceSource
	// backgroundImage is the background image
	backgroundImage *ebiten.Image
}

// NewGame returns a Game instance to be ran.
func NewGame(cfg *Cfg) *Game {
	fs, err := text.NewGoTextFaceSource(bytes.NewReader(assets.EmpireStateNF_ttf))
	if err != nil {
		panic(fmt.Sprintf("failed to load font file: %v", err))
	}

	bgImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Background))
	if err != nil {
		panic(fmt.Sprintf("failed to background texture: %v", err))
	}

	cfg.faceSource = fs
	cfg.backgroundImage = bgImg

	return &Game{
		Cfg:   cfg,
		state: mainMenu,
	}
}

// Update manages user input, handles physics processes and updates game states.
func (g *Game) Update() error {
	switch g.state {
	case mainMenu:

	}

	return nil
}

// Draw draws the appropriate images to the screen based on game state.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case mainMenu:
		g.drawMainMenu(screen)
	}
}

// Layout returns the screen's logical width and height.
func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

// Run runs begins the game loop.
func (g *Game) Run() error {
	ebiten.SetWindowSize(g.Cfg.ScreenWidth, g.Cfg.ScreenHeight)
	ebiten.SetWindowTitle(g.Cfg.WindowTitle)

	return ebiten.RunGame(g)
}

func (g *Game) drawMainMenu(screen *ebiten.Image) {
	imgOpts := &ebiten.DrawImageOptions{}
	imgOpts.GeoM.Scale(3, 2)
	screen.DrawImage(g.backgroundImage, imgOpts)

	titleTextFace := &text.GoTextFace{
		Source: g.Cfg.faceSource,
		Size:   titleFontSize,
	}
	menuTextFace := &text.GoTextFace{
		Source: g.Cfg.faceSource,
		Size:   fontSize,
	}

	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(g.ScreenWidth)/2, titleFontSize)
	opts.PrimaryAlign = text.AlignCenter
	opts.LineSpacing = titleFontSize
	opts.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, "Pong", titleTextFace, opts)

	opts = &text.DrawOptions{}
	opts.GeoM.Translate(float64(g.ScreenWidth)/2, 3*titleFontSize)
	opts.LineSpacing = fontSize
	opts.PrimaryAlign = text.AlignCenter
	opts.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, "Single Player\nMultiplayer", menuTextFace, opts)
}
