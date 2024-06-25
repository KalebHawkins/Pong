package pong

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gtihub.com/KalebHawkins/pong/assets"
)

const (
	fontSize      = 48
	titleFontSize = fontSize * 1.5
	paddleWidth   = 20
	paddleHeight  = 80
	maxScore      = 10
	accelerator   = 1
)

// Game is a structure containing the game data and configuration.
type Game struct {
	*Cfg
	state

	playerOne paddle
	playerTwo paddle
	ball
	isPaused bool
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
	// paddleImage is the image used for the paddle
	paddleImage *ebiten.Image
	// ballImage is the image used for the ball
	ballImage *ebiten.Image

	// verticalLine is used as the central delimiter of the screen vertically..
	verticalLine *ebiten.Image
	// horizontalLine is used as the central delimiter of the screen horizontally.
	horizontalLine *ebiten.Image
}

// NewGame returns a Game instance to be ran.
func NewGame(cfg *Cfg) *Game {
	fs, err := text.NewGoTextFaceSource(bytes.NewReader(assets.EmpireStateNF_ttf))
	if err != nil {
		panic(fmt.Sprintf("failed to load font file: %v", err))
	}

	bgImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Background_png))
	if err != nil {
		panic(fmt.Sprintf("failed to background texture: %v", err))
	}

	paddleImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Paddle_png))
	if err != nil {
		panic(fmt.Sprintf("failed to load paddle texture: %v", err))
	}

	ballImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Ball_png))
	if err != nil {
		panic(fmt.Sprintf("failed to load ball texture: %v", err))
	}

	cfg.faceSource = fs
	cfg.backgroundImage = bgImg
	cfg.paddleImage = paddleImg
	cfg.ballImage = ballImg
	cfg.verticalLine = ebiten.NewImage(1, cfg.ScreenHeight)
	cfg.horizontalLine = ebiten.NewImage(cfg.ScreenWidth, 1)

	return &Game{
		Cfg:   cfg,
		state: mainMenu,
		playerOne: paddle{
			x:      10 + paddleWidth/2,
			y:      cfg.ScreenHeight / 2,
			dy:     10,
			score:  0,
			sprite: ebiten.NewImageFromImage(cfg.paddleImage),
		},
		playerTwo: paddle{
			x:      cfg.ScreenWidth - paddleWidth/2 - 10,
			y:      cfg.ScreenHeight / 2,
			dy:     10,
			score:  0,
			sprite: ebiten.NewImageFromImage(cfg.paddleImage),
		},
		ball: ball{
			x:      cfg.ScreenWidth / 2,
			y:      cfg.ScreenHeight / 2,
			dx:     3,
			dy:     -2,
			sprite: ebiten.NewImageFromImage(cfg.ballImage),
		},
	}
}

// Update manages user input, handles physics processes and updates game states.
func (g *Game) Update() error {
	switch g.state {
	case mainMenu:
		// Start the game.
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = gameLoop
		}

	case gameLoop:
		// Pause the game.
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.isPaused = !g.isPaused

			// If the game is paused the ball should stop moving and no input for the players should
			// be procesesd.
			if g.isPaused {
				g.ball.prevDx = g.ball.dx
				g.ball.prevDy = g.ball.dy
				g.ball.dx = 0
				g.ball.dy = 0
			} else {
				g.ball.dx = g.ball.prevDx
				g.ball.dy = g.ball.prevDy
			}
		}

		// Check if player has reached the win limit.
		if g.playerOne.score >= maxScore || g.playerTwo.score >= maxScore {
			g.state = gameOver
		}

		if !g.isPaused {
			g.ball.x += g.ball.dx
			g.ball.y += g.ball.dy

			// Handle player Input
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				if g.playerOne.y >= 0+(paddleHeight/2) {
					g.playerOne.y -= g.playerOne.dy
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				if g.playerOne.y <= g.ScreenHeight-(paddleHeight/2) {
					g.playerOne.y += g.playerOne.dy
				}
			}

			if g.ball.y < g.playerTwo.y {
				if g.playerTwo.y >= 0+(paddleHeight/2) {
					g.playerTwo.y -= g.playerTwo.dy
				}
			}
			if g.ball.y > g.playerTwo.y {
				if g.playerTwo.y <= g.ScreenHeight-(paddleHeight/2) {
					g.playerTwo.y += g.playerTwo.dy
				}
			}

			// Handle ball collision with top and bottom of the screen
			if g.ball.y <= 0 || g.ball.y >= g.ScreenHeight {
				g.ball.dy *= -1
			}

			// Handle ball collision with playerOne
			if g.ball.x <= g.playerOne.x+(paddleWidth/2) && g.ball.x >= g.playerOne.x-(paddleWidth/2) &&
				g.ball.y >= g.playerOne.y-(paddleHeight/2) && g.ball.y <= g.playerOne.y+(paddleHeight/2) {
				g.hitBall()
			}

			// Handle ball collision with playerTwo
			if g.ball.x >= g.playerTwo.x-paddleWidth && g.ball.x <= g.playerTwo.x+(paddleWidth/2) &&
				g.ball.y >= g.playerTwo.y-(paddleHeight/2) && g.ball.y <= g.playerTwo.y+(paddleHeight/2) {
				g.hitBall()
			}

			// Handle player getting a score.
			if g.ball.x > g.ScreenWidth {
				g.resetGameLoop()
				g.playerOne.score++
			}
			if g.ball.x < 0 {
				g.resetGameLoop()
				g.playerTwo.score++
			}
		}

	case gameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = mainMenu
		}
	}
	return nil
}

// Draw draws the appropriate images to the screen based on game state.
func (g *Game) Draw(screen *ebiten.Image) {
	imgOpts := &ebiten.DrawImageOptions{}
	imgOpts.GeoM.Scale(3, 2)
	screen.DrawImage(g.backgroundImage, imgOpts)

	op := &ebiten.DrawImageOptions{}
	g.Cfg.verticalLine.Fill(color.White)
	op.GeoM.Translate(-float64(g.Cfg.verticalLine.Bounds().Dx())/2, 0)
	op.GeoM.Translate(float64(g.Cfg.ScreenWidth/2), 0)
	screen.DrawImage(g.Cfg.verticalLine, op)

	op = &ebiten.DrawImageOptions{}
	g.Cfg.horizontalLine.Fill(color.White)
	op.GeoM.Translate(0, -float64(g.Cfg.horizontalLine.Bounds().Dy()/2))
	op.GeoM.Translate(0, float64(g.Cfg.ScreenHeight)/2)
	screen.DrawImage(g.Cfg.horizontalLine, op)

	switch g.state {
	case mainMenu:
		g.drawMainMenu(screen)
	case gameLoop:
		g.drawGameLoop(screen)
		if g.isPaused {
			g.drawPauseMenu(screen)
		}
	case gameOver:
		g.drawGameLoop(screen)
		g.drawGameOverMenu(screen)
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
	text.Draw(screen, "Press Space\nto Play", menuTextFace, opts)
}

func (g *Game) drawGameLoop(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(g.playerOne.sprite.Bounds().Dx())/2, -float64(g.playerOne.sprite.Bounds().Dy())/2)
	op.GeoM.Translate(float64(g.playerOne.x), float64(g.playerOne.y))
	screen.DrawImage(g.playerOne.sprite, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(g.playerTwo.sprite.Bounds().Dx())/2, -float64(g.playerTwo.sprite.Bounds().Dy())/2)
	op.GeoM.Translate(float64(g.playerTwo.x), float64(g.playerTwo.y))
	screen.DrawImage(g.playerTwo.sprite, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(g.ball.sprite.Bounds().Dx())/2, -float64(g.ball.sprite.Bounds().Dy())/2)
	op.GeoM.Translate(float64(g.ball.x), float64(g.ball.y))
	screen.DrawImage(g.ball.sprite, op)

	face := &text.GoTextFace{
		Source: g.faceSource,
		Size:   fontSize,
	}

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(g.ScreenWidth/4), 0)
	textOp.LineSpacing = fontSize
	textOp.PrimaryAlign = text.AlignCenter
	textOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, fmt.Sprintf("%d", g.playerOne.score), face, textOp)

	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(float64(g.ScreenWidth)-float64(g.ScreenWidth/4), 0)
	textOp.LineSpacing = fontSize
	textOp.PrimaryAlign = text.AlignCenter
	textOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, fmt.Sprintf("%d", g.playerTwo.score), face, textOp)

}

func (g *Game) drawPauseMenu(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.Cfg.ScreenWidth)/2, fontSize)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, "GAME PAUSED", &text.GoTextFace{
		Source: g.Cfg.faceSource,
		Size:   fontSize,
	}, op)
}

func (g *Game) drawGameOverMenu(screen *ebiten.Image) {
	var winner string
	if g.playerOne.score >= maxScore {
		winner = "Player One"
	} else {
		winner = "Player Two"
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.Cfg.ScreenWidth)/2, fontSize)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, fmt.Sprintf("GAME OVER\n\n\n%s\nWins!", winner), &text.GoTextFace{
		Source: g.Cfg.faceSource,
		Size:   fontSize,
	}, op)
}

func (g *Game) resetGameLoop() {
	g.ball.dx = 3
	g.ball.dy = 2
	g.ball.x = g.ScreenWidth / 2
	g.ball.y = g.ScreenHeight / 2

	g.playerOne.y = g.ScreenHeight / 2
	g.playerTwo.y = g.ScreenHeight / 2
}

func (g *Game) hitBall() {
	g.ball.dx *= -1
	if g.ball.dx < 0 {
		g.ball.dx -= accelerator
	} else {
		g.ball.dx += accelerator
	}
}
