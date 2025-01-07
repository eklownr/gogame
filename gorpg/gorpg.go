package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1920 / 3
	screenHeight = 1080 / 3
	imgSize      = 48
	SPEED        = time.Second / 4
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
	rectTop     = Point{0, 0}
	rectBot     = Point{imgSize, imgSize}
	gameSpeed   = SPEED
)

type Game struct {
	Player     *Charakters
	costomer   *[]Charakters
	worker     *[]Charakters
	plants     *[]Plant
	lastUpdate time.Time
	tick       bool
}
type Sprite struct {
	img *ebiten.Image
	pos Point
}
type Charakters struct {
	*Sprite
	speed  float64
	dest   Point
	coin   int
	wallet int
}
type Plant struct {
	*Sprite
	variety string
}
type Point struct {
	x, y float64
}

// Idle faceing front
func (g *Game) idle() {
	if g.tick {
		rectTop.x = imgSize - imgSize // 0
		rectTop.y = imgSize - imgSize // 0
		rectBot.x = imgSize           // 48
		rectBot.y = imgSize           // 48
	} else {
		rectTop.x = imgSize
		rectTop.y = imgSize - imgSize
		rectBot.x = imgSize * 2
		rectBot.y = imgSize
	}
}

// set new position and player image
func (g *Game) dirDown() {
	g.Player.pos.y += g.Player.speed
	if g.tick {
		rectTop.x = imgSize * 2
		rectTop.y = imgSize - imgSize
		rectBot.x = imgSize * 3
		rectBot.y = imgSize
	} else {
		rectTop.x = imgSize * 3
		rectTop.y = imgSize - imgSize
		rectBot.x = imgSize * 4
		rectBot.y = imgSize
	}
}
func (g *Game) dirUp() {
	g.Player.pos.y -= g.Player.speed
	if g.tick {
		rectTop.x = imgSize
		rectTop.y = imgSize
		rectBot.x = imgSize * 2
		rectBot.y = imgSize * 2
	} else {
		rectTop.x = imgSize - imgSize
		rectTop.y = imgSize
		rectBot.x = imgSize
		rectBot.y = imgSize * 2
	}
}
func (g *Game) dirLeft() {
	g.Player.pos.x -= g.Player.speed
	if g.tick {
		rectTop.x = imgSize * 2
		rectTop.y = imgSize * 2
		rectBot.x = imgSize * 3
		rectBot.y = imgSize * 3
	} else {
		rectTop.x = imgSize * 3
		rectTop.y = imgSize * 2
		rectBot.x = imgSize * 4
		rectBot.y = imgSize * 3
	}
}
func (g *Game) dirRight() {
	g.Player.pos.x += g.Player.speed
	if g.tick {
		rectTop.x = imgSize - imgSize
		rectTop.y = imgSize * 3
		rectBot.x = imgSize
		rectBot.y = imgSize * 4
	} else {
		rectTop.x = imgSize * 2
		rectTop.y = imgSize * 3
		rectBot.x = imgSize * 3
		rectBot.y = imgSize * 4
	}
}

func (g *Game) Update() error {
	g.readKeys()

	// check Animation tick every 60 FPS
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	if g.tick {
		g.tick = false
	} else {
		g.tick = true
	}
	g.lastUpdate = time.Now() // update lastUpdate
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(skyBlue) // background collor

	///////// draw img player ///////////
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.pos.x, g.Player.pos.y)

	screen.DrawImage(
		g.Player.img.SubImage(
			image.Rect(int(rectTop.x), int(rectTop.y), int(rectBot.x), int(rectBot.y)),
		).(*ebiten.Image),
		opts,
	)
}

// vim-keys to move "hjkl" or Arrowkeys
func (g *Game) readKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.dirDown()
	} else {
		g.idle()
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
		Player: &Charakters{
			Sprite: &Sprite{
				img: playerImg,
				pos: Point{200, 200},
			},
			speed: 2,
		},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// check for errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
