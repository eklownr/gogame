package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
)

var (
	dirUp    = Point{x: 0, y: -1}
	dirDown  = Point{x: 0, y: 1}
	dirRight = Point{x: 1, y: 0}
	dirLeft  = Point{x: -1, y: 0}
)

const (
	gameSpeed    = time.Second / 6
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

	// set direction back to screen if outside
	//if g.snake{Point{x}} <= 0 {
	//}

	return nil
}

// Update the memory of the snake (*snake) not the copy snake
func (g *Game) updateSnake(snake *[]Point, dir Point) {
	head := (*snake)[0]
	newHead := Point{
		x: head.x + dir.x,
		y: head.y + dir.y,
	}
	// update the snake. Head + body-1
	*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
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
		color.RGBA{255, 0, 0, 255},
		true,
	)
}

func (g *Game) Layout(outsidewith, outsideheight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
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

	g.spawnFood()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
