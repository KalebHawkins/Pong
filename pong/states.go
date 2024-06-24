package pong

// state represents the games state
type state int

const (
	mainMenu state = iota
	gameLoop
	gameOver
)
