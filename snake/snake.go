package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	snake []int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsidewith, outsideheight int) (int, int) {
	return screenWidth * 2, screenHeight * 2
}

func main() {
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
