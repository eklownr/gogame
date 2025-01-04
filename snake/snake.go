package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/eklownr/pretty"
	_ "github.com/eklownr/pretty"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"golang.org/x/exp/rand"
)

var (
	dirUp           = Point{x: 0, y: -1}
	dirDown         = Point{x: 0, y: 1}
	dirRight        = Point{x: 1, y: 0}
	dirLeft         = Point{x: -1, y: 0}
	gameSpeed       = SPEED
	mplusFaceSource *text.GoTextFaceSource
	red             = color.RGBA{255, 0, 0, 255}
	yellow          = color.RGBA{220, 200, 0, 255}
	green           = color.RGBA{0, 220, 0, 255}
	blue            = color.RGBA{0, 20, 120, 255}
	purple          = color.RGBA{200, 0, 200, 255}
	orange          = color.RGBA{180, 160, 0, 255}
	white           = color.RGBA{255, 255, 255, 255}
	black           = color.RGBA{0, 0, 0, 255}
)

const (
	screenWidth  = 1920 / 2
	screenHeight = 1080 / 2
	gridSize     = 20
	maxGameSpeed = time.Second / 12
	SPEED        = time.Second / 6
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	snakeColor color.Color
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
	gamePause  bool
	fullWindow bool
	score      int
	screen     *ebiten.Image
}

func (g *Game) Layout(outsidewith, outsideheight int) (int, int) {
	return screenWidth, screenHeight
}

// set random food position
func (g *Game) spawnFood() {
	g.food = Point{
		rand.Intn(screenWidth / gridSize),
		rand.Intn(screenHeight / gridSize),
	}
}

// vim-keys to move "hjkl" or Arrowkeys
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // Pause the game
		g.pauseGame()
	} else if ebiten.IsKeyPressed(ebiten.KeyQ) { // Quit the Game!
		g.quitGame()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) { // Full screen
		g.fullScreen()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) { // SpawnFood
		g.spawnFood()
		g.drawFood(g.screen, purple)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyY) { // buy yellow snake
		if g.score >= 20 {
			g.snakeColor = yellow
			g.score -= 20
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyP) { // buy yellow snake
		if g.score >= 30 {
			g.snakeColor = purple
			g.score -= 30
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) { // buy yellow snake
		if g.score >= 40 {
			g.snakeColor = red
			g.score -= 40
		}
	}
}

func (g *Game) quitGame() {
	println("Warning! Quit the game!")
	ebiten.SetRunnableOnUnfocused(false)
	os.Exit(1)
}

func (g *Game) Update() error {
	g.readKeys()
	// check if the snake can move, ckeck every 60 FPS
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	if g.gameOver {
		g.score = 0          // set score back to 0
		g.snakeColor = white // set snake back color to white
	}
	g.lastUpdate = time.Now() // update lastUpdate
	g.updateSnake(&g.snake, g.direction)
	return nil
}

// update snak, check collision and pause
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
	if g.gamePause {
		return // pause the game
	}

	//// check if snake is out of sceen, set snake back to screen
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
		*snake = append([]Point{newHead}, *snake...) // Grow the snake
		g.spawnFood()
		g.spawnFood() // set rand position for food
		g.score += 10
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
	g.screen = screen
	// Draw background color
	оп := &ebiten.DrawImageOptions{}
	img := ebiten.NewImage(screenWidth, screenHeight)
	img.Fill(color.NRGBA{0, 20, 80, 255}) // Blue color
	screen.DrawImage(img, оп)

	if !g.gameOver && !g.gamePause {
		// draw the snake eatch Time = gameSpeed
		for _, p := range g.snake {
			vector.DrawFilledRect(
				screen,
				float32(p.x*gridSize),
				float32(p.y*gridSize),
				gridSize,
				gridSize,
				g.snakeColor,
				true,
			)
		}
		g.drawFood(screen, green)

		// Add text Score: ... to the screen
		s := fmt.Sprint(g.score)
		addText(screen, 18, "Score: ", green, 120, 20)
		addText(screen, 18, s, green, 200, 20)
	}
	// Game Over
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
		// Add text Score: ... to the screen
		s := fmt.Sprint(g.score)
		addText(screen, 18, "Score: ", green, 120, 20)
		addText(screen, 18, s, green, 200, 20)
		// Pause the game
	} else if g.gamePause {
		vector.DrawFilledRect(
			screen,
			float32(20),     // x position
			float32(20),     // y position
			screenWidth-40,  // width size
			screenHeight-40, // Height size
			blue,
			true,
		)
		addText(screen, 56, "Pause", black, screenWidth+5, screenHeight/3+4)
		addText(screen, 56, "Pause", yellow, screenWidth, screenHeight/3)
		addText(screen, 34, "Pause the Game - Esc", yellow, screenWidth, screenHeight/3+100)
		addText(screen, 34, "Quit the game - q", yellow, screenWidth, screenHeight/3+200)
		addText(screen, 34, "Full screen - f", yellow, screenWidth, screenHeight/3+300)
		addText(screen, 34, "*************** Shop ***************", purple, screenWidth, screenHeight/3+450)
		addText(screen, 34, "Buy yellow snake, cost: 20 - y", green, screenWidth, screenHeight/3+550)
		addText(screen, 34, "Buy purple snake, cost: 30 - p", green, screenWidth, screenHeight/3+650)
		addText(screen, 34, "Buy red snake, cost:    40 - r", green, screenWidth, screenHeight/3+750)
		// Add text Score: ... to the screen
		s := fmt.Sprint(g.score)
		addText(screen, 28, "Score: ", green, 200, 90)
		addText(screen, 28, s, green, 400, 90)
	}
}

// Draw the food
func (g *Game) drawFood(screen *ebiten.Image, color color.Color) {
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*gridSize),
		float32(g.food.y*gridSize),
		gridSize,
		gridSize,
		color,
		true,
	)
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

// Enter-key to restart the Game
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
}

// Escape-key to Pause the game
func (g *Game) pauseGame() {
	if !g.gamePause {
		g.gamePause = true
	} else {
		g.gamePause = false
	}
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
	ebiten.SetWindowResizingMode(ebiten.WindowResizingMode())

	// initial the snake in the center the game
	g := &Game{
		snake: []Point{
			{
				x: screenWidth / gridSize / 2,
				y: screenHeight / gridSize / 2,
			}},
		direction: Point{x: 1, y: 0},
	}
	g.snakeColor = white
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

// //////////////// TEST //////////////////////////////
// replace this line with below section to test memory
func printMemStats() {}

// // This section below is only for Testing... Memory and other stuff
// func printMemStats() {
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)
// 	println("Alloc Heap Memory stat in Mb: ", bToMb(m.Alloc))
// 	println("Total Alloc Memory stat in Mb: ", bToMb(m.TotalAlloc))
// 	println("Total SYS Heap and stack - memory in Mb: ", bToMb(m.Sys))
// 	println("Garbage collector times: ", m.NumGC)
// 	println("*************************")
// }
// func bToMb(b uint64) uint64 {
// 	return b / 1000 / 1000
// }

// Print colors in terminal
func println(arg ...interface{}) {
	pretty.Pl(arg...)
}
