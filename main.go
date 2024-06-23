package main

import "gtihub.com/KalebHawkins/pong/pong"

func main() {
	cfg := &pong.Cfg{
		ScreenWidth:  640,
		ScreenHeight: 480,
		WindowTitle:  "Pong",
	}

	g := pong.NewGame(cfg)
	g.Run()
}
