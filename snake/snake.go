package main

import (
	"bytes"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"golang.org/x/exp/rand"
)

var (
	dirUp           = Point{x: 0, y: -1}
	dirDown         = Point{x: 0, y: 1}
	dirRight        = Point{x: 1, y: 0}
	dirLeft         = Point{x: -1, y: 0}
	gameSpeed       = time.Second / 6
	mplusFaceSource *text.GoTextFaceSource
	red             = color.RGBA{255, 0, 0, 255}
	yellow          = color.RGBA{220, 200, 0, 255}
	green           = color.RGBA{0, 220, 0, 255}
	blue            = color.RGBA{0, 0, 0, 255}
	purple          = color.RGBA{200, 0, 200, 255}
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
}

func (g *Game) spawnFood() {
	g.food = Point{
		rand.Intn(screenWidth / gridSize),
		rand.Intn(screenHeight / gridSize),
	}
}

func (g *Game) readKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.direction = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.direction = dirRight
	}
}

func (g *Game) Update() error {
	g.readKeys()
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now() // update lastUpdate
	g.updateSnake(&g.snake, g.direction)
	return nil
}

// Update the memory of the snake (*snake) not a copy of snake
func (g *Game) updateSnake(snake *[]Point, dir Point) {
	head := (*snake)[0]
	newHead := Point{
		x: head.x + dir.x,
		y: head.y + dir.y,
	}
	//  collision detection
	if g.isBadCollision(newHead, *snake) {
		g.gameOver = true
		return // Stop the Game
	}

	// set snake back to screen if out of screen
	if newHead.x < 0 {
		newHead.x = screenWidth / gridSize
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	} else if newHead.y < 0 {
		newHead.y = screenHeight / gridSize
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	} else if newHead.y >= screenHeight/gridSize {
		newHead.y = 0
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	} else if newHead.x >= screenWidth/gridSize {
		newHead.x = 0
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	} else if newHead == g.food {
		*snake = append([]Point{newHead}, *snake...)
		g.spawnFood()
		gameSpeed -= time.Second / 66 // get faster eatch food
	} else {
		// Move and update the snakes Head + body-1
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	}
}
func (g *Game) isBadCollision(p Point, snake []Point) bool {
	//// check if snake is out of sceen
	//if p.x < 0 || p.y < 0 || p.x >= screenWidth/gridSize || p.y >= screenHeight/gridSize {
	//	return true
	//}

	// check snake body collide with self body
	for _, sp := range snake {
		if sp == p {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake {
		vector.DrawFilledRect(
			screen,
			float32(p.x*gridSize),
			float32(p.y*gridSize),
			gridSize,
			gridSize,
			color.White,
			true,
		)
	}
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*gridSize),
		float32(g.food.y*gridSize),
		gridSize,
		gridSize,
		red,
		true,
	)

	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		t := "Game Over!"
		w, h := text.Measure(
			t,
			face,
			face.Size,
		)

		op := &text.DrawOptions{}
		op.GeoM.Translate(
			screenWidth/2-w/2, screenHeight/2-h/2,
		)
		op.ColorScale.ScaleWithColor(yellow)
		text.Draw(
			screen,
			t,
			face,
			op,
		)
	}
}

func (g *Game) Layout(outsidewith, outsideheight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// game over
	s, err := text.NewGoTextFaceSource(
		bytes.NewReader(
			fonts.MPlus1pRegular_ttf,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	// setup window size nd title
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Snake")

	// initial the snake in to the game
	g := &Game{
		snake: []Point{
			{
				x: screenWidth / gridSize / 2,
				y: screenHeight / gridSize / 2,
			}},
		direction: Point{x: 1, y: 0},
	}
	// init food to the game
	g.spawnFood()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
