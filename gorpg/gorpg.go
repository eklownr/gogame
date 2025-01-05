package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1920 / 3
	screenHeight = 1080 / 3
)

var (
	runnerImage *ebiten.Image
	skyBlue     = color.RGBA{120, 180, 255, 255}
	red         = color.RGBA{255, 0, 0, 255}
	yellow      = color.RGBA{220, 200, 0, 255}
	green       = color.RGBA{0, 220, 0, 255}
	blue        = color.RGBA{0, 20, 120, 255}
	purple      = color.RGBA{200, 0, 200, 255}
	orange      = color.RGBA{180, 160, 0, 255}
	white       = color.RGBA{255, 255, 255, 255}
	black       = color.RGBA{0, 0, 0, 255}
)

type Game struct {
	Player   *Sprite
	costomer *Sprite
	worker   *Sprite
}
type Sprite struct {
	img       *ebiten.Image
	direction Point
	speed     float64
}

type plant struct {
	Sprite
	types string
}

func (g *Game) dirRight() {
	g.Player.direction.x += g.Player.speed
}
func (g *Game) dirLeft() {
	g.Player.direction.x -= g.Player.speed
}
func (g *Game) dirUp() {
	g.Player.direction.y -= g.Player.speed
}
func (g *Game) dirDown() {
	g.Player.direction.y += g.Player.speed
}

func (g *Game) Update() error {
	g.readKeys()
	return nil
}

type Point struct {
	x, y float64
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(skyBlue) // background collor

	///////// draw img player ///////////
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.direction.x, g.Player.direction.y)

	screen.DrawImage(
		g.Player.img.SubImage(
			image.Rect(0, 0, 40, 40),
		).(*ebiten.Image),
		opts,
	)
}

// vim-keys to move "hjkl" or Arrowkeys
// if g.direction is Up you can´t move Down. Same for all direction
func (g *Game) readKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.dirDown()
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.dirUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dirLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dirRight()
	}
	//	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) && g.gameOver == true {
	//		g.restartGame(&g.snake) // move snake to start position and remove body
	//		g.gameOver = false      // Start the game with Enter-key if GameOver.
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // Pause the game
	//
	//		g.pauseGame()
	//	} else if ebiten.IsKeyPressed(ebiten.KeyQ) { // Quit the Game!
	//
	//		g.quitGame()
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) { // Full screen
	//
	//		g.fullScreen()
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) { // set Food point randomly
	//
	//		g.randMultiFood() // TEST
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyY) { // buy yellow snake
	//
	//		if g.score >= 20 {
	//			g.snakeColor = yellow
	//			g.score -= 20
	//		}
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyP) { // buy yellow snake
	//
	//		if g.score >= 30 {
	//			g.snakeColor = purple
	//			g.score -= 30
	//		}
	//	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) { // buy yellow snake
	//
	//		if g.score >= 40 {
	//			g.snakeColor = red
	//			g.score -= 40
	//		}
	//	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	///////// Animation ///////////
	// Decode an image from the image file's byte slice.
	//	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	runnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gopher Mart")

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/player.png")
	checkErr(err)

	g := &Game{
		Player: &Sprite{
			img:       playerImg,
			direction: Point{200, 200},
			speed:     2,
		},
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
