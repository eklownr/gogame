package main

import (
	"bytes"
	_ "embed"
	"gorpg/tilemaps"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth   = 1920 / 3
	screenHeight  = 1080 / 3
	imgSize       = 48
	SPEED         = time.Second / 4
	houseTileSize = 64
	SampleRate    = 44100
	wheat         = "wheat"
	tomato        = "tomato"
)

//go:embed assets/sound/LostVillage.ogg
var audioBG []byte

//go:embed assets/sound/Village.ogg
var audioVillage []byte

//go:embed assets/sound/Coin.ogg
var audioCoin []byte

//go:embed assets/sound/Fx.ogg
var audioFx []byte

var (
	skyBlue         = color.RGBA{120, 180, 255, 255}
	red             = color.RGBA{255, 0, 0, 255}
	red_rect        = color.RGBA{255, 0, 0, 40}
	blue            = color.RGBA{0, 20, 120, 255}
	blue_rect       = color.RGBA{0, 20, 120, 40}
	yellow          = color.RGBA{220, 200, 0, 255}
	green           = color.RGBA{0, 220, 0, 255}
	purple          = color.RGBA{200, 0, 200, 255}
	orange          = color.RGBA{180, 160, 0, 255}
	white           = color.RGBA{255, 255, 255, 255}
	black           = color.RGBA{0, 0, 0, 255}
	gameSpeed       = SPEED
	PlayerSpeed     = 3.0
	diagonalSpeed   = 0.8
	tileSize        = 16
	mplusFaceSource *text.GoTextFaceSource
	coin_anim       = 0
	plant_anim      = 0
)

type Game struct {
	Player          *Characters
	costomers       []*Characters
	workers         []*Characters
	coins           []*Objects
	house           []*Objects
	plants          []*Objects
	lastUpdate      time.Time
	tick            bool
	fullWindow      bool
	gameOver        bool
	gamePause       bool
	village         *ebiten.Image
	bgImg           *ebiten.Image
	tilemapImg      *ebiten.Image
	tilemapImgWater *ebiten.Image
	plantImg        *ebiten.Image
	workImg         *ebiten.Image
	workerIdleImg   *ebiten.Image
	coinImg         *ebiten.Image
	smokeSprite     Sprite
	tilemapJSON1    *tilemaps.TilemapJSON
	tilemapJSON2    *tilemaps.TilemapJSON
	tilemapJSON3    *tilemaps.TilemapJSON
	scene           int
}
type Sprite struct {
	img          *ebiten.Image
	pos          Point
	prePos       Point
	rectPos      image.Rectangle
	rectTop      Point // Sprite amination
	rectBot      Point // Sprite amination
	active       bool
	frameCounter int
	frame        int
}
type Characters struct {
	*Sprite
	Dir
	speed       float64
	dest        Point
	coin        int
	wallet      int
	plantBasket int
}
type Objects struct {
	*Sprite
	variety  string
	dest     Point
	picked   bool
	pickable bool
}
type Point struct {
	x, y float64
}
type Dir struct {
	down, up, right, left bool
}

// Idle faceing front animation
func (g *Game) idleWorkers(i int) {
	// show animation subImage
	if g.tick {
		g.workers[i].rectTop.x = imgSize - imgSize // 0
		g.workers[i].rectTop.y = imgSize - imgSize // 0
		g.workers[i].rectBot.x = imgSize           // 48
		g.workers[i].rectBot.y = imgSize           // 48
	} else {
		g.workers[i].rectTop.x = imgSize
		g.workers[i].rectTop.y = imgSize - imgSize
		g.workers[i].rectBot.x = imgSize * 2
		g.workers[i].rectBot.y = imgSize
	}
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

// check collision Objects with charakter: Player-worker
func (g *Game) Collision_Character_Character(obj Characters, char Characters) bool {
	// Player....
	charakter_position := image.Rect(
		int(char.pos.x+imgSize/4),
		int(char.pos.y+imgSize/4),
		int(char.pos.x+imgSize/2),
		int(char.pos.y+imgSize/2))

	// Worker
	object_position := image.Rect(
		int(obj.pos.x),
		int(obj.pos.y),
		int(obj.pos.x+imgSize/2),
		int(obj.pos.y+imgSize/2))

	if object_position.Overlaps(charakter_position) {
		return true
	}
	return false
}

// check collision obj-obj: worker-plant
func (g *Game) Collision_worker_plant(worker Objects, plant Objects) bool {
	// Player....
	worker_position := image.Rect(
		int(worker.pos.x+imgSize/4),
		int(worker.pos.y+imgSize/4),
		int(worker.pos.x+imgSize/2),
		int(worker.pos.y+imgSize/2))

	// Worker
	plant_position := image.Rect(
		int(plant.pos.x),
		int(plant.pos.y),
		int(plant.pos.x+imgSize/2),
		int(plant.pos.y+imgSize/2))

	if worker_position.Overlaps(plant_position) {
		g.smokeSprite.active = true
		return true
	}
	return false
}

// check collision Objects with charakter: any obj-Player
func (g *Game) Collision_Object_Caracter(obj Objects, char Characters) bool {
	// Player....
	charakter_position := image.Rect(
		int(char.pos.x+imgSize/4),
		int(char.pos.y+imgSize/4),
		int(char.pos.x+imgSize/2),
		int(char.pos.y+imgSize/2))

	object_position := image.Rect(
		int(obj.pos.x),
		int(obj.pos.y),
		int(obj.pos.x+imgSize/2),
		int(obj.pos.y+imgSize/2))

	if obj.variety == "coin" || obj.variety == "wheat" || obj.variety == "tomato" {
		object_position = image.Rect(
			int(obj.pos.x-imgSize/4+10),
			int(obj.pos.y-imgSize/4+10),
			int(obj.pos.x+imgSize/4),
			int(obj.pos.y+imgSize/4))
	}
	if obj.variety == "budda" {
		object_position = image.Rect(
			int(obj.pos.x),
			int(obj.pos.y),
			int(obj.pos.x+imgSize/2),
			int(obj.pos.y+imgSize/2))
	}
	if obj.variety == "house" {
		object_position = image.Rect(
			int(obj.pos.x),
			int(obj.pos.y),
			int(obj.pos.x+houseTileSize-10),
			int(obj.pos.y+imgSize-10))
	}
	if obj.variety == "small_house" {
		object_position = image.Rect(
			int(obj.pos.x),
			int(obj.pos.y),
			int(obj.pos.x+imgSize-10),
			int(obj.pos.y+imgSize-10))
	}

	if object_position.Overlaps(charakter_position) {
		return true
	}
	return false
}

// TEST check buildings collision: Point-Point
func (g *Game) checkCollision(p1 Point, p2 Point) bool {
	if p1.x >= p2.x-imgSize &&
		p1.x <= p2.x+imgSize &&
		p1.y >= p2.y-imgSize &&
		p1.y <= p2.y+imgSize {
		return true
	}
	return false
}

// collide with budda Action: set player pos, span coin,
func (g *Game) buddaCollision() {
	// Portal Player
	g.Player.pos.x = screenWidth / 2
	g.Player.pos.y = screenHeight / 2
	// playSound
	playSound(audioFx)
	//	// change scene
	//	if g.scene < 3 {
	//		g.scene++
	//	} else {
	//		g.scene = 0
	//	}

}

// Move Workers to dest pos
func (g *Game) moveCharacters(c *Characters) {
	if c.pos != c.dest {
		c.img = g.workerIdleImg
		if c.pos.x < c.dest.x {
			c.pos.x++
		}
		if c.pos.x > c.dest.x {
			c.pos.x--
		}
		if c.pos.y < c.dest.y {
			c.pos.y++
		}
		if c.pos.y > c.dest.y {
			c.pos.y--
		}
	} else {
		c.img = g.workImg
		g.plants[5].active = true // TEST
	}
}

// check Animation tick every 60 FPS
func (g *Game) plantFrameAnim(plant *Objects) {
	var speed = 120
	if plant.frame < speed*5 {
		plant.frameCounter++
	}
	if plant.frameCounter < speed {
		plant.frame = 1
	} else if plant.frameCounter < speed*2 {
		plant.frame = 2
	} else if plant.frameCounter < speed*3 {
		plant.frame = 3
	} else if plant.frameCounter < speed*4 {
		plant.frame = 4
	} else if plant.frameCounter < speed*5 {
		plant.frame = 5
		plant.pickable = true
	}
}

// ///// Update function
func (g *Game) Update() error {
	g.Player.prePos = g.Player.pos // save old position before readKeys()
	g.readKeys()                   // read keys and move player
	g.coin_animation()

	// plants animation
	for _, plant := range g.plants {
		if plant.active && plant.frame < 5 {
			g.plantFrameAnim(plant)
		}
	}
	// TEST Move workers to new dest pos for every new scene
	for i := range g.workers { // Idle animation for all workers
		g.idleWorkers(i)
		g.moveCharacters(g.workers[i])

		// TEST
		if g.scene == 2 {
			g.workers[i].dest = Point{180 + (float64(i) * 40), 300}
		} else if g.scene == 3 {
			g.workers[i].dest = Point{30 + (float64(i) * 30), 20}
		} else if g.scene == 1 {
			g.workers[i].dest = Point{50, 10 + (float64(i) * 30)}
		} else if g.scene == 0 {
			g.workers[i].dest = Point{200 + (float64(i) * 20), 90}
		}

	}

	// Player border collision
	if g.Player.pos.x < 0-imgSize/2 {
		g.Player.pos.x = screenWidth - imgSize/2
	} else if g.Player.pos.x > screenWidth-imgSize/2 {
		g.Player.pos.x = 0 - imgSize/2
	} else if g.Player.pos.y < 0-imgSize/2 {
		g.Player.pos.y = screenHeight - imgSize/2
	} else if g.Player.pos.y > screenHeight {
		g.Player.pos.y = 0 - imgSize/2
	}
	//Player collide with []workers
	for i := range g.workers {
		if g.Collision_Character_Character(*g.workers[i], *g.Player) {
			if g.workers[i].coin < 2 && g.Player.coin > 0 {
				g.workers[i].coin++
				g.Player.coin--
				g.smokeSprite.active = true
			}
			// playSound
			playSound(audioFx)
			// play smoke animation
		}
	}
	//Player collide with []house
	for i := range g.house {
		if g.Collision_Object_Caracter(*g.house[i], *g.Player) {
			g.Player.pos = g.Player.prePos
			g.smokeSprite.active = true
			if g.house[i].variety == "budda" {
				g.buddaCollision()
			}
		}
	}
	//Player collide with []plants if active
	for i := range g.plants {
		if g.Collision_Object_Caracter(*g.plants[i], *g.Player) {
			if g.plants[i].pickable {
				g.smokeSprite.active = true
				// pick plant
				g.plants[i].active = false
				g.plants[i].pickable = false
				g.Player.plantBasket++
			}
		}
	}
	// Player collide with []coin
	for i := range g.coins {
		if g.Collision_Object_Caracter(*g.coins[i], *g.Player) && g.coins[i].picked == false {
			if g.Player.coin < g.Player.wallet { // add coins to your wallet
				g.Player.coin++
				println("You have: ", g.Player.coin, "coins") // TEST
				playSound(audioCoin)
				g.coins[i].picked = true
				g.coins[i].pos = Point{
					x: -100,
					y: -100,
				}
			}
		}
	}

	/////////////////////////////////////
	// check Animation tick every 60 FPS. 2 values On or Off
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	if g.tick {
		g.tick = false
	} else {
		g.tick = true
	}
	g.lastUpdate = time.Now() // update lastUpdate

	// last in Update()
	return nil
}

// ////// Draw function
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(skyBlue) // background collor

	if g.gamePause {
		g.pause(screen)
		return
	}

	op := &ebiten.DrawImageOptions{}
	if g.scene == 0 {
		///////// draw background ///////////
		op.GeoM.Translate(20, 20)

		screen.DrawImage(
			g.bgImg.SubImage(
				image.Rect(0, 0, 600, 370),
			).(*ebiten.Image),
			op,
		)
		op.GeoM.Reset()
	}
	if g.scene == 1 {
		/////////// draw bg tile layers ////////////
		for _, layer := range g.tilemapJSON1.Layers {
			for index, id := range layer.Data {
				x := index % layer.Width
				y := index / layer.Width
				x *= tileSize
				y *= tileSize

				srcX := (id - 1) % 22
				srcY := (id - 1) / 22
				srcX *= tileSize
				srcY *= tileSize

				op.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(
					g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image),
					op,
				)
				op.GeoM.Reset()
			}
		}
	}
	if g.scene == 2 {
		/////////// draw bg tile layers ////////////
		for _, layer := range g.tilemapJSON2.Layers {
			for index, id := range layer.Data {
				x := index % layer.Width
				y := index / layer.Width
				x *= tileSize
				y *= tileSize

				srcX := (id - 1) % 22
				srcY := (id - 1) / 22
				srcX *= tileSize
				srcY *= tileSize

				op.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(
					g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image),
					op,
				)
				op.GeoM.Reset()
			}
		}
	}
	if g.scene == 3 {
		/////////// draw bg tile layers from waterTileImg ////////////
		for _, layer := range g.tilemapJSON3.Layers {
			for index, id := range layer.Data {
				x := index % layer.Width
				y := index / layer.Width
				x *= tileSize
				y *= tileSize

				srcX := (id - 1) % 28
				srcY := (id - 1) / 28
				srcX *= tileSize
				srcY *= tileSize

				op.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(
					g.tilemapImgWater.SubImage(image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image),
					op,
				)
				op.GeoM.Reset()
			}
		}
	}

	//	///////// draw all house big and small  ////////////
	for i := range g.house {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(g.house[i].pos.x, g.house[i].pos.y) // house position x, y
		screen.DrawImage(
			g.village.SubImage(
				g.house[i].rectPos,
			).(*ebiten.Image),
			opt,
		)
		opt.GeoM.Reset()
	}

	/// Draw coin at same pos as Game constructor g.coins.pos in main() ///
	for i := 0; i < 10; i++ {
		g.drawCoin(screen, g.coins[i].pos.x, g.coins[i].pos.y, *g.coins[i], i)
	}

	/// Draw Workers at same pos as Game constructor in main() ///
	for i := range g.workers {
		// draw coin carring on workers head
		g.carry_objects(screen, g.workers[i].pos.x, g.workers[i].pos.y, g.workers[i].coin, g.coinImg)
		// draw all workers
		g.drawWorker(screen, g.workers[i].pos.x, g.workers[i].pos.y, i)
	}

	///////// draw coin player caring on the head. SubImg 0,0,10,10 /////////
	g.carry_objects(screen, g.Player.pos.x, g.Player.pos.y, g.Player.coin, g.coinImg)
	// SubImg 0,0,16,16
	g.carry_plant(screen, g.Player.pos.x, g.Player.pos.y, g.Player.plantBasket, g.plantImg)

	///// Draw all plants  if active ///
	for i := range g.plants {
		if g.plants[i].active {
			g.drawPlanst(screen, g.plants[i].pos.x, g.plants[i].pos.y, g.plants[i].variety, g.plants[i].frame) // wheat and tomato
			if g.plants[i].picked {
				g.plants[i].active = false
			}
		}
	}

	///////// draw img player ///////////
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.pos.x, g.Player.pos.y)
	// amination position to Player.img.SubImage(image.Rect(0, 0, imgSize, imgSize))
	screen.DrawImage(
		g.Player.img.SubImage(
			image.Rect(int(g.Player.rectTop.x), int(g.Player.rectTop.y), int(g.Player.rectBot.x), int(g.Player.rectBot.y)),
		).(*ebiten.Image),
		opts,
	)
	/////// TEST Draw player and house collision rect
	// vector.StrokeRect(screen, float32(g.Player.pos.x+imgSize/4),float32(g.Player.pos.y+imgSize/4),imgSize/2,imgSize/2,3.0,color.RGBA{122, 222, 0, 100},false)
	// vector.StrokeRect(screen,float32(g.housePos.x)+float32(g.house[0].rectPos.Min.X),float32(g.housePos.y)+float32(g.house[0].rectPos.Min.Y),houseTileSize,imgSize,3.0,color.RGBA{222, 122, 0, 100},false)

	// if active
	g.drawSmoke(screen, g.Player.pos.x, g.Player.pos.y)
}

// /////// draw images caring on the head ////////////
func (g *Game) carry_objects(screen *ebiten.Image, x, y float64, amount int, img *ebiten.Image) {
	optst := &ebiten.DrawImageOptions{}
	for i := 3; i < 3+amount; i++ { // i=3 3 pix apart
		optst.GeoM.Translate(x+imgSize/2-3, y+float64(2.0*i)-10.0)

		screen.DrawImage(
			img.SubImage(
				image.Rect(0, 0, 10, 10),
			).(*ebiten.Image),
			optst,
		)
		optst.GeoM.Reset()
	}
}

// /////// draw images caring on the head ////////////
func (g *Game) carry_plant(screen *ebiten.Image, x, y float64, amount int, img *ebiten.Image) {
	optst := &ebiten.DrawImageOptions{}
	for i := 3; i < 3+amount; i++ { // i=3 3 pix apart
		optst.GeoM.Translate(x+imgSize/2-3, y+float64(2.0*i)-10.0)

		screen.DrawImage(
			img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			optst,
		)
		optst.GeoM.Reset()
	}
}

func (g *Game) smoke_animation() {
	if g.tick {
		plant_anim = 32
		if time.Since(g.lastUpdate) < gameSpeed/2 {
			plant_anim = 32 * 2
		}
	} else {
		plant_anim = 32 * 3
		if time.Since(g.lastUpdate) < gameSpeed/2 {
			plant_anim = 32 * 4
		}
	}
}
func (g *Game) drawSmoke(screen *ebiten.Image, x, y float64) {
	if g.smokeSprite.active {
		g.smoke_animation()
		option := &ebiten.DrawImageOptions{}
		option.GeoM.Translate(x, y) // position x, y
		screen.DrawImage(
			g.smokeSprite.img.SubImage(
				image.Rect(plant_anim, 0, plant_anim+32, 32),
			).(*ebiten.Image),
			option,
		)
		option.GeoM.Reset()
		g.smokeSprite.active = false
	}
}

// TEST plants animation
func (g *Game) plant_animation(frame int) {
	plant_anim = 16 * frame
}
func (g *Game) drawPlanst(screen *ebiten.Image, x, y float64, variety string, frame int) {
	g.plant_animation(frame) // activate animation
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(x, y) // coin position x, y
	if variety == "wheat" {
		screen.DrawImage(
			g.plantImg.SubImage(
				//image.Rect(plant_anim, 0, plant_anim+imgSize, imgSize),
				image.Rect(plant_anim, 0, plant_anim+16, 16),
			).(*ebiten.Image),
			option,
		)
		option.GeoM.Reset()
	} else if variety == "tomato" {
		screen.DrawImage(
			g.plantImg.SubImage(
				//image.Rect(plant_anim, 0, plant_anim+imgSize, imgSize),
				image.Rect(plant_anim, 16, plant_anim+16, 16+16),
			).(*ebiten.Image),
			option,
		)
		option.GeoM.Reset()

	}
}

func (g *Game) drawWorker(screen *ebiten.Image, x, y float64, i int) {
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(x, y) // worker position x, y
	screen.DrawImage(
		g.workers[i].img.SubImage(
			image.Rect(int(g.workers[i].rectTop.x), int(g.workers[i].rectTop.y), int(g.workers[i].rectBot.x), int(g.workers[i].rectBot.y)),
		).(*ebiten.Image),
		option,
	)
	option.GeoM.Reset()
}

func (g *Game) coin_animation() {
	if g.tick {
		coin_anim = 0
		if time.Since(g.lastUpdate) < gameSpeed/2 {
			coin_anim = 10
		}
	} else {
		coin_anim = 20
		if time.Since(g.lastUpdate) < gameSpeed/2 {
			coin_anim = 30
		}
	}
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
			image.Rect(coin_anim, 0, coin_anim+10, 10),
		).(*ebiten.Image),
		option,
	)
	option.GeoM.Reset()
}

// Arrowkeys to move or vim-keys "hjkl"
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeyQ) { // Quit the game
		g.quitGame()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // Pause the game
		g.pauseGame()
	} else if inpututil.IsKeyJustPressed(ebiten.Key0) { // Pause the game
		g.scene = 0
	} else if inpututil.IsKeyJustPressed(ebiten.Key1) { // Pause the game
		g.scene = 1
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) { // Pause the game
		g.scene = 2
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) { // Pause the game
		g.scene = 3
	}

}

// F key for full screen
func (g *Game) fullScreen() {
	if !g.fullWindow {
		ebiten.SetFullscreen(true)
		g.fullWindow = true
	} else {
		g.fullWindow = false
		ebiten.SetFullscreen(false)
	}
}

// Q key for quit
func (g *Game) quitGame() {
	println("Warning! Quit the game!")
	ebiten.SetRunnableOnUnfocused(false)
	os.Exit(1)
}

// Escape-key to Pause the game
func (g *Game) pauseGame() {
	if !g.gamePause {
		g.gamePause = true
	} else {
		g.gamePause = false
	}
}
func (g *Game) pause(screen *ebiten.Image) {
	// Pause the game
	vector.DrawFilledRect(
		screen,
		float32(20),     // x position
		float32(20),     // y position
		screenWidth-40,  // width size
		screenHeight-40, // Height size
		blue,
		true,
	)
	addText(screen, 32, "Pause", black, screenWidth+5, screenHeight/3+4)
	addText(screen, 32, "Pause", yellow, screenWidth, screenHeight/3)
	addText(screen, 20, "Pause the Game - Esc", yellow, screenWidth, screenHeight/3+100)
	addText(screen, 20, "Quit the game - q", yellow, screenWidth, screenHeight/3+200)
	addText(screen, 20, "Full screen - f", yellow, screenWidth, screenHeight/3+300)
	addText(screen, 20, "Change scene key: 0-3", purple, screenWidth, screenHeight/3+400)
	addText(screen, 20, "*********************", green, screenWidth, screenHeight/3+500)
}

func addText(screen *ebiten.Image, textSize int, t string, color color.Color, width, height float64) {
	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   float64(textSize),
	}

	// t := "YOURE TEXT"
	w, h := text.Measure(
		t,
		face,
		face.Size,
	)

	op := &text.DrawOptions{}
	op.GeoM.Translate(
		width/2-w/2, height/2-h/2,
	)
	op.ColorScale.ScaleWithColor(color)
	text.Draw(
		screen,
		t,
		face,
		op,
	)
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

	// Text, font
	textsource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	checkErr(err)
	mplusFaceSource = textsource

	// TilemapJSON1
	tilemapJSON1, err := tilemaps.NewTilemapJSON("assets/map/level1_bg.json")
	checkErr(err)

	// TilemapJSON2
	tilemapJSON2, err := tilemaps.NewTilemapJSON("assets/map/level2_bg.json")
	checkErr(err)

	// TilemapJSON2 Water
	tilemapJSON3, err := tilemaps.NewTilemapJSON("assets/map/water_bg.json")
	checkErr(err)

	// load tilemapImg image to tilemap 1 - 2
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/map/tileset_floor.png")
	checkErr(err)

	// load tilemapImg image to tilemap water
	tilemapImgWater, _, err := ebitenutil.NewImageFromFile("assets/map/TilesetWater.png")
	checkErr(err)

	// load village image
	village, _, err := ebitenutil.NewImageFromFile("assets/images/village.png")
	checkErr(err)

	// load background image
	bgImg, _, err := ebitenutil.NewImageFromFile("assets/images/grass.png")
	checkErr(err)

	// load Player image
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/playerBlue.png")
	checkErr(err)

	// load Worker image
	workerImg, _, err := ebitenutil.NewImageFromFile("assets/images/player.png")
	checkErr(err)

	// load Work image
	workImg, _, err := ebitenutil.NewImageFromFile("assets/images/workers.png")
	checkErr(err)

	// load coin image
	coinImg, _, err := ebitenutil.NewImageFromFile("assets/images/coin2.png")
	checkErr(err)

	// load coin image
	plantImg, _, err := ebitenutil.NewImageFromFile("assets/images/plants.png")
	checkErr(err)

	// load smoke image
	smokeImg, _, err := ebitenutil.NewImageFromFile("assets/images/smoke.png")
	checkErr(err)

	// Game constructor
	g := &Game{
		Player: &Characters{
			Sprite: &Sprite{
				img: playerImg,
				pos: Point{screenWidth/2 - (imgSize / 2), screenHeight/2 - (imgSize / 2)},
			},
			speed:  PlayerSpeed,
			coin:   0,
			wallet: 1,
		},
	}
	g.Player.rectTop = Point{g.Player.pos.y + imgSize/4, g.Player.pos.y + imgSize/4}
	g.Player.rectBot = Point{imgSize / 2, imgSize / 2}

	// add 10 workers
	for i := 0; i < 10; i++ {
		g.workers = append(g.workers, &Characters{
			Sprite: &Sprite{
				img:     workerImg,
				pos:     Point{40, 20*float64(i) + 60},
				rectPos: image.Rect(0, 0, imgSize/2, imgSize/2),
			},
			speed: 1.5,
			dest:  Point{screenWidth - imgSize - (float64(i * imgSize)), screenHeight/2 - imgSize - (float64(i * imgSize))},
		})
	}
	for i := range g.workers { //set rectTop and rectBot for animation
		g.workers[i].rectTop = Point{0, 0}
		g.workers[i].rectBot = Point{imgSize, imgSize}
	}

	// add 10 coins
	for i := 1; i < 11; i++ {
		g.coins = append(g.coins, &Objects{
			Sprite: &Sprite{
				img:     coinImg,
				pos:     Point{screenWidth/2 + 30 + float64(i)*10, screenHeight/2 + houseTileSize - 30.0},
				rectPos: image.Rect(0, 0, imgSize/2, imgSize/2),
			},
			variety: "coin",
		})
	}
	// add 2 plants: wheat and tomato
	for i := 0; i < 4; i++ {
		g.plants = append(g.plants, &Objects{
			Sprite: &Sprite{
				img:     plantImg,
				pos:     Point{178 + float64(i)*40, 300},
				rectPos: image.Rect(0, 0, imgSize/2, imgSize/2),
				active:  true,
			},
			variety: "wheat",
		})

		g.plants = append(g.plants, &Objects{
			Sprite: &Sprite{
				img:     plantImg,
				pos:     Point{60 + float64(i)*40, 40},
				rectPos: image.Rect(0, 0, imgSize/2, imgSize/2),
				active:  false,
			},
			variety: "tomato",
		})
	}
	//	add house objects
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     Point{250, houseTileSize},
			rectPos: image.Rect(0, 0, houseTileSize, imgSize),
		},
		variety: "house",
	})
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     Point{100, 100},
			rectPos: image.Rect(houseTileSize, 0, houseTileSize*2, imgSize),
		},
		variety: "house",
	})
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     Point{screenWidth/2 + houseTileSize, screenHeight/2 + houseTileSize},
			rectPos: image.Rect(0, imgSize, imgSize, imgSize*2),
		},
		variety: "budda",
	})
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     Point{400, imgSize},
			rectPos: image.Rect(houseTileSize*2+imgSize, 0, houseTileSize*2+imgSize*2, imgSize),
		},
		variety: "small_house",
	})
	g.house = append(g.house, &Objects{
		Sprite: &Sprite{
			img:     village,
			pos:     Point{500, imgSize},
			rectPos: image.Rect(houseTileSize*2+imgSize, imgSize, houseTileSize*2+imgSize*2, imgSize*2),
		},
		variety: "small_house",
	})

	// Add Images and tilemapJSON
	g.coinImg = coinImg
	g.bgImg = bgImg
	g.village = village
	g.tilemapImg = tilemapImg
	g.tilemapImgWater = tilemapImgWater
	g.plantImg = plantImg
	g.workImg = workImg
	g.workerIdleImg = workerImg

	g.smokeSprite = Sprite{
		img:    smokeImg,
		pos:    Point{50, 50},
		active: false,
	}

	g.tilemapJSON1 = tilemapJSON1
	g.tilemapJSON2 = tilemapJSON2
	g.tilemapJSON3 = tilemapJSON3

	g.scene = 1

	////// play background music //////
	_ = audio.NewContext(SampleRate)
	stream, err := vorbis.DecodeWithSampleRate(SampleRate, bytes.NewReader(audioBG))
	checkErr(err)

	// infinite loop Bg music
	audioPlayer, _ := audio.CurrentContext().NewPlayer(
		audio.NewInfiniteLoop(stream,
			int64(len(audioBG)*6*SampleRate)))

	//audioPlayer, _ := audio.CurrentContext().NewPlayer(stream)
	// you pass the audio player to your game struct, and just call
	audioPlayer.SetVolume(0.3)
	audioPlayer.Play() //when you want your music to start, and
	// audioPlayer.Pause()

	// Start game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func playSound(sound []byte) {
	//_ = audio.NewContext(SampleRate)
	stream, err := vorbis.DecodeWithSampleRate(SampleRate, bytes.NewReader(sound))
	checkErr(err)
	audioPlayer, _ := audio.CurrentContext().NewPlayer(stream)
	audioPlayer.Play()
}
