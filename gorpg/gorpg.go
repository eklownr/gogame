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
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth   = 1920 / 3
	screenHeight  = 1080 / 3
	imgSize       = 48
	SPEED         = time.Second / 4
	houseTileSize = 64
)

var (
	skyBlue       = color.RGBA{120, 180, 255, 255}
	red           = color.RGBA{255, 0, 0, 255}
	red_rect      = color.RGBA{255, 0, 0, 40}
	blue          = color.RGBA{0, 20, 120, 255}
	blue_rect     = color.RGBA{0, 20, 120, 40}
	yellow        = color.RGBA{220, 200, 0, 255}
	green         = color.RGBA{0, 220, 0, 255}
	purple        = color.RGBA{200, 0, 200, 255}
	orange        = color.RGBA{180, 160, 0, 255}
	white         = color.RGBA{255, 255, 255, 255}
	black         = color.RGBA{0, 0, 0, 255}
	gameSpeed     = SPEED
	PlayerSpeed   = 3.0
	diagonalSpeed = 0.8
)

type Game struct {
	Player     *Charakters
	costomer   []*Charakters
	worker     []*Charakters
	lastUpdate time.Time
	tick       bool
	fullWindow bool
	bgImg      *ebiten.Image
	village    *ebiten.Image
	housePos   Point
	coins      []*Objects
	house      []*Objects
}
type Sprite struct {
	img     *ebiten.Image
	pos     Point
	prePos  Point
	rectPos image.Rectangle
	rectTop Point
	rectBot Point
}
type Charakters struct {
	*Sprite
	Dir
	speed  float64
	dest   Point
	coin   int
	wallet int
}
type Objects struct {
	*Sprite
	variety string
	dest    Point
	picked  bool
}
type Point struct {
	x, y float64
}
type Dir struct {
	down, up, right, left bool
}

// Idle faceing front animation
func (g *Game) idle() {
	// show animation subImage
	if g.tick {
		g.Player.rectTop.x = imgSize - imgSize // 0
		g.Player.rectTop.y = imgSize - imgSize // 0
		g.Player.rectBot.x = imgSize           // 48
		g.Player.rectBot.y = imgSize           // 48
	} else {
		g.Player.rectTop.x = imgSize
		g.Player.rectTop.y = imgSize - imgSize
		g.Player.rectBot.x = imgSize * 2
		g.Player.rectBot.y = imgSize
	}
}

// set new position and player images(animation)
func (g *Game) dirDown() {
	if g.Player.Dir.right || g.Player.Dir.left {
		g.Player.speed = PlayerSpeed * diagonalSpeed
	}
	g.Player.pos.y += g.Player.speed
	// show animation subImage
	if g.tick {
		g.Player.rectTop.x = imgSize * 2
		g.Player.rectTop.y = imgSize - imgSize
		g.Player.rectBot.x = imgSize * 3
		g.Player.rectBot.y = imgSize
	} else {
		g.Player.rectTop.x = imgSize * 3
		g.Player.rectTop.y = imgSize - imgSize
		g.Player.rectBot.x = imgSize * 4
		g.Player.rectBot.y = imgSize
	}
	g.Player.Dir.down = false
	g.Player.Dir.right = false
	g.Player.Dir.left = false
	g.Player.speed = PlayerSpeed
}
func (g *Game) dirUp() {
	if g.Player.Dir.right || g.Player.Dir.left {
		g.Player.speed = PlayerSpeed * diagonalSpeed
	}
	g.Player.pos.y -= g.Player.speed
	// show animation subImage
	if g.tick {
		g.Player.rectTop.x = imgSize * 2
		g.Player.rectTop.y = imgSize
		g.Player.rectBot.x = imgSize * 3
		g.Player.rectBot.y = imgSize * 2
	} else {
		g.Player.rectTop.x = imgSize * 3
		g.Player.rectTop.y = imgSize
		g.Player.rectBot.x = imgSize * 4
		g.Player.rectBot.y = imgSize * 2
	}
	g.Player.Dir.up = false
	g.Player.Dir.right = false
	g.Player.Dir.left = false
	g.Player.speed = PlayerSpeed
}
func (g *Game) dirLeft() {
	if g.Player.Dir.up || g.Player.Dir.down {
		g.Player.speed = PlayerSpeed * diagonalSpeed
	}
	g.Player.pos.x -= g.Player.speed
	// show animation subImage
	if g.tick {
		g.Player.rectTop.x = imgSize * 2
		g.Player.rectTop.y = imgSize * 2
		g.Player.rectBot.x = imgSize * 3
		g.Player.rectBot.y = imgSize * 3
	} else {
		g.Player.rectTop.x = imgSize * 3
		g.Player.rectTop.y = imgSize * 2
		g.Player.rectBot.x = imgSize * 4
		g.Player.rectBot.y = imgSize * 3
	}
	g.Player.Dir.left = false
	g.Player.Dir.up = false
	g.Player.Dir.down = false
	g.Player.speed = PlayerSpeed
}
func (g *Game) dirRight() {
	if g.Player.Dir.up || g.Player.Dir.down {
		g.Player.speed = PlayerSpeed * diagonalSpeed
	}
	g.Player.pos.x += g.Player.speed
	// show animation subImage
	if g.tick {
		g.Player.rectTop.x = imgSize - imgSize
		g.Player.rectTop.y = imgSize * 3
		g.Player.rectBot.x = imgSize
		g.Player.rectBot.y = imgSize * 4
	} else {
		g.Player.rectTop.x = imgSize * 2
		g.Player.rectTop.y = imgSize * 3
		g.Player.rectBot.x = imgSize * 3
		g.Player.rectBot.y = imgSize * 4
	}
	g.Player.Dir.right = false
	g.Player.Dir.up = false
	g.Player.Dir.down = false
	g.Player.speed = PlayerSpeed
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

// check collision Objects to Point
func (g *Game) Collision_Object_Point(obj []*Objects, p2 Point) bool {
	for i := range obj {
		if obj[i].pos == p2 {
			obj[i].picked = true
			obj[i].pos = Point{
				x: -100,
				y: -100,
			}
			g.Player.coin++
			return true
		}
	}
	return false
}

// check buildings collision
func (g *Game) checkCollision(p1 Point, p2 Point) bool {
	if p1.x >= p2.x-imgSize/2+10 &&
		p1.x <= p2.x+imgSize/2 &&
		p1.y >= p2.y-imgSize/2+10 &&
		p1.y <= p2.y+imgSize/2 {
		return true
	}
	return false
}

// update function
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
	} else if g.checkCollision(g.Player.pos, g.house[0].pos) { //collision with house
		println("You are at home")
		g.Player.pos = g.Player.prePos
	}
	for i := range g.coins {
		if g.checkCollision(g.coins[i].pos, g.Player.pos) && g.coins[i].picked == false {
			g.Player.coin++
			println("You have NOW: ", g.Player.coin, "coins")
			g.coins[i].picked = true
			g.coins[i].pos = Point{
				x: -100,
				y: -100,
			}
		}
	}

	if g.Player.pos == g.coins[0].pos {
		println("player.pos == coins[0].pos")
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

// Draw function
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
	opt.GeoM.Translate(300, houseTileSize) // house position x, y

	screen.DrawImage(
		g.village.SubImage(
			image.Rect(0, 0, houseTileSize, imgSize),
		).(*ebiten.Image),
		opt,
	)

	///////// draw cooin player caring on head ////////////
	optst := &ebiten.DrawImageOptions{}
	for i := 3; i < 3+g.Player.coin; i++ {
		optst.GeoM.Translate(g.Player.pos.x+imgSize/2-3, g.Player.pos.y+float64(2.0*i)-10.0)

		screen.DrawImage(
			g.coins[0].img.SubImage(
				image.Rect(0, 0, imgSize, imgSize),
			).(*ebiten.Image),
			optst,
		)
		optst.GeoM.Reset()
	}
	/// Draw coin at same as Game constructor g.coins.pos ///
	for i := 0; i < 10; i++ {
		g.drawCoin(screen, g.coins[i].pos.x, g.coins[i].pos.y, *g.coins[i], i)
	}

	///////// draw img player ///////////
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.pos.x, g.Player.pos.y)

	screen.DrawImage(
		g.Player.img.SubImage(
			image.Rect(0, 0, imgSize, imgSize),
		).(*ebiten.Image),
		opts,
	)
	// TEST draw Player collision rect
	vector.DrawFilledRect(
		screen,
		float32(g.Player.pos.x),
		float32(g.Player.pos.y),
		float32(imgSize),
		float32(imgSize),
		blue_rect,
		true,
	)
	// TEST draw Player collision rect
	vector.DrawFilledRect(
		screen,
		float32(g.Player.pos.x+imgSize/4),
		float32(g.Player.pos.y+imgSize/4),
		float32(imgSize/2),
		float32(imgSize/2),
		red_rect,
		true,
	)
}

func (g *Game) drawCoin(screen *ebiten.Image, x, y float64, coin Objects, index int) {
	if coin.picked {
		g.coins[index].pos = Point{-100, -100} // outside of screen
		return
	}
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(x, y) // coin position x, y
	screen.DrawImage(
		g.coins[index].img.SubImage(
			image.Rect(0, 0, imgSize, imgSize),
		).(*ebiten.Image),
		option,
	)
	g.coins[index].rectPos = image.Rect(int(x-imgSize/4+10), int(y-imgSize/4+10), int(imgSize/4), int(imgSize/4))
	option.GeoM.Reset()

	// TEST draw coin rect
	vector.DrawFilledRect(
		screen,
		float32(x-imgSize/4+10),
		float32(y-imgSize/4+10),
		float32(imgSize/4),
		float32(imgSize/4),
		red_rect,
		true,
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

// check for errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	// load coin image
	coinImg, _, err := ebitenutil.NewImageFromFile("assets/images/coin.png")
	checkErr(err)

	// Game constructor
	g := &Game{
		Player: &Charakters{
			Sprite: &Sprite{
				img: playerImg,
				pos: Point{screenWidth/2 - (imgSize / 2), screenHeight/2 - (imgSize / 2)},
			},
			speed: PlayerSpeed,
			coin:  2,
		},
	}
	g.Player.rectTop = Point{g.Player.pos.y + imgSize/4, g.Player.pos.y + imgSize/4}
	g.Player.rectBot = Point{imgSize / 2, imgSize / 2}

	// add 10 coins
	for i := 1; i < 11; i++ {
		g.coins = append(g.coins, &Objects{
			Sprite: &Sprite{
				img:     coinImg,
				pos:     Point{200, 20*float64(i) + 60},
				rectPos: image.Rect(0-imgSize/4, 0-imgSize/4, imgSize/2, imgSize/2),
			},
		})
	}
	g.bgImg = bgImg
	g.village = village
	g.housePos = Point{300, houseTileSize}
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     g.housePos,
			rectPos: image.Rect(0, 0, houseTileSize, imgSize),
		},
	})

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
