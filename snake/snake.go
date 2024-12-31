package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/eklownr/pretty"
	_ "github.com/eklownr/pretty"
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
	score           = 0
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
	maxGameSpeed = time.Second / 12
	SPEED        = time.Second / 6
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

// vim-keys to move "hjkl" or arrowkeys
// if g.direction is Up you can´t move Down. Same for all direction
func (g *Game) readKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction != dirUp {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction != dirDown {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction != dirRight {
		g.direction = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction != dirLeft {
		g.direction = dirRight
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) && g.gameOver == true {
		g.restartGame(&g.snake) // move snake to start position and remove body
		g.gameOver = false      // Start the game with Enter-key if GameOver.
	}
}

func (g *Game) Update() error {
	g.readKeys()
	// update speed
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now() // update lastUpdate
	g.updateSnake(&g.snake, g.direction)
	return nil
}

// Update the memory of the snake (*snake) not a copy of the snake
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
	// Grow the snake
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
	} else if newHead == g.food { // Eat Food
		*snake = append([]Point{newHead}, *snake...)
		g.spawnFood()
		score += 10
		if gameSpeed > maxGameSpeed {
			gameSpeed -= time.Second / 66 // get faster eatch food
		}
	} else {
		// Move and update the snakes Head + body-1
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	}
}
func (g *Game) isBadCollision(p Point, snake []Point) bool {
	//// check if snake is out of sceen
	//if p.x < 0 || p.y < 0 || p.x >= screenWidth/gridSize || p.y >= screenHeight/gridSize {return true }

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

	s := fmt.Sprint(score)
	// Add text Score: ... to the screen
	addText(screen, 18, "Score: ", green, 120, 20)
	addText(screen, 18, s, green, 200, 20)

	if g.gameOver {
		addText(screen, 48, "Hit Enter to play", yellow, screenWidth, screenHeight/3)
		//addText(screen, "Game Over!", yellow, screenWidth/2, screenHeight/2)

		vector.DrawFilledRect(
			screen,
			float32(screenWidth/4),  // x position
			float32(screenHeight/3), // y position
			screenWidth/2,           // width size
			screenHeight/3,          // Height size
			red,
			true,
		)

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

func (g *Game) Layout(outsidewith, outsideheight int) (int, int) {
	return screenWidth, screenHeight
}

// Key-Enter restarts the Game
func (g *Game) restartGame(snake *[]Point) {
	newHead := Point{ // Place new head at center of the screen
		x: screenWidth / gridSize / 2,
		y: screenHeight / gridSize / 2,
	}
	*snake = (*snake)[:0]            // remove snake body
	*snake = append(*snake, newHead) // add a new head

	// move snake to the right
	g = &Game{
		direction: Point{x: 1, y: 0},
	}
	gameSpeed = SPEED // set game-speed back to start-speed
	score = 0         // set back score to 0
}

func main() {
	// print memStats
	println("*** Mem before, first in Main() ***")
	printMemStats()

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

	// initial the snake in the center the game
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

	// print memStat
	println("*** Mem just before ebiten.RunGame(&Game) ***")
	printMemStats()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// print memStat
	println("*** Mem after, last of main() ***")
	printMemStats()

}

// replace this line with below section to test memory
//func printMemStats() {}

// This section below is only for Testing... Memory and other stuff
func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	println("Alloc Heap Memory stat in Mb: ", bToMb(m.Alloc))
	println("Total Alloc Memory stat in Mb: ", bToMb(m.TotalAlloc))
	println("Total SYS Heap and stack - memory in Mb: ", bToMb(m.Sys))
	println("Garbage collector times: ", m.NumGC)
	println("*************************")
}
func bToMb(b uint64) uint64 {
	return b / 1000 / 1000
}

// Some colors to print
func println(arg ...interface{}) {
	pretty.Pl(arg...)
}
