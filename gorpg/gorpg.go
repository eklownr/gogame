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

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
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
	img  *ebiten.Image
	x, y float64
}

type plant struct {
	Sprite
	types string
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(skyBlue) // background collor

	///////// draw img player ///////////
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.x, g.Player.y)

	screen.DrawImage(
		g.Player.img.SubImage(
			image.Rect(0, 0, 40, 40),
		).(*ebiten.Image),
		opts,
	)
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
			img: playerImg,
			x:   250,
			y:   250,
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
