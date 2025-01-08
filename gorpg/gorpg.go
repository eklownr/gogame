package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth   = 1920 / 3
	screenHeight  = 1080 / 3
	imgSize       = 48
	SPEED         = time.Second / 4
	houseTileSize = 64
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

type Dir struct {
	down, up, right, left bool
}
type Game struct {
	Player     *Charakters
	costomer   *[]Charakters
	worker     *[]Charakters
	plants     *[]Plant
	lastUpdate time.Time
	tick       bool
	fullWindow bool
	bgImg      *ebiten.Image
	village    *ebiten.Image
	housePos   Point
	coins      []Point
}
type Sprite struct {
	img    *ebiten.Image
	pos    Point
	prePos Point
}
type Charakters struct {
	*Sprite
	Dir
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

// Idle faceing front animation
func (g *Game) idle() {
	// show animation subImage
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

// set new position and player images(animation)
func (g *Game) dirDown() {
	//	if g.Player.Dir.right || g.Player.Dir.left {
	//		g.Player.pos.y += g.Player.speed * 0.5
	//	} else {
	//		g.Player.pos.y += g.Player.speed
	//	}
	g.Player.pos.y += g.Player.speed
	// show animation subImage
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
	g.Player.Dir.down = false
}
func (g *Game) dirUp() {
	g.Player.pos.y -= g.Player.speed
	// show animation subImage
	if g.tick {
		rectTop.x = imgSize * 2
		rectTop.y = imgSize
		rectBot.x = imgSize * 3
		rectBot.y = imgSize * 2
	} else {
		rectTop.x = imgSize * 3
		rectTop.y = imgSize
		rectBot.x = imgSize * 4
		rectBot.y = imgSize * 2
	}
	g.Player.Dir.up = false
}
func (g *Game) dirLeft() {
	g.Player.pos.x -= g.Player.speed
	// show animation subImage
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
	g.Player.Dir.left = false
}
func (g *Game) dirRight() {
	g.Player.pos.x += g.Player.speed
	// show animation subImage
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
	g.Player.Dir.right = false
}

// F-key for full screen
func (g *Game) fullScreen() {
	if !g.fullWindow {
		ebiten.SetFullscreen(true)
		g.fullWindow = true
	} else {
		g.fullWindow = false
		ebiten.SetFullscreen(false)
	}
}

func (g *Game) Update() error {
	g.Player.prePos = g.Player.pos // save old position
	g.readKeys()                   // read keys and move player
	// Player collision
	if g.Player.pos.x < 0 {
		g.Player.pos = g.Player.prePos
	} else if g.Player.pos.x > screenWidth-imgSize {
		g.Player.pos = g.Player.prePos
	} else if g.Player.pos.y < 0 {
		g.Player.pos = g.Player.prePos
	} else if g.Player.pos.y > screenHeight-imgSize-5 {
		g.Player.pos = g.Player.prePos
	} else if g.Player.pos.x >= g.housePos.x-32 && //collision with house
		g.Player.pos.x <= g.housePos.x+imgSize &&
		g.Player.pos.y >= g.housePos.y-imgSize/2 &&
		g.Player.pos.y <= g.housePos.y+imgSize/2 {
		println("You are at home")
		g.Player.pos = g.Player.prePos
	}

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

	///////// draw background ///////////
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(20, 20)

	screen.DrawImage(
		g.bgImg.SubImage(
			image.Rect(0, 0, 600, 370),
		).(*ebiten.Image),
		op,
	)

	///////// draw hoouse 0 ////////////
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(300, houseTileSize)

	screen.DrawImage(
		g.village.SubImage(
			image.Rect(0, 0, houseTileSize, imgSize),
		).(*ebiten.Image),
		opt,
	)

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
		g.Player.Dir.down = true
	} else {
		g.idle()
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.dirUp()
		g.Player.Dir.up = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dirLeft()
		g.Player.Dir.left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dirRight()
		g.Player.Dir.right = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) { // Full screen
		g.fullScreen()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Window properties
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gopher Mart")

	// load village image
	village, _, err := ebitenutil.NewImageFromFile("assets/images/village.png")
	checkErr(err)

	// load background image
	bgImg, _, err := ebitenutil.NewImageFromFile("assets/images/grass.png")
	checkErr(err)

	// load player image
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/playerBlue.png")
	checkErr(err)

	// Game constructor
	g := &Game{
		Player: &Charakters{
			Sprite: &Sprite{
				img: playerImg,
				pos: Point{screenWidth/2 - (imgSize / 2), screenHeight/2 - (imgSize / 2)},
			},
			speed: 3,
		},
	}
	g.bgImg = bgImg
	g.village = village
	g.housePos = Point{300, houseTileSize}
	g.Player.Dir.down = false
	g.Player.Dir.up = false
	g.Player.Dir.right = false
	g.Player.Dir.left = false

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
