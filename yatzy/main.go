package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state. 60Hz
func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen. 60Hz
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Yatzy")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
