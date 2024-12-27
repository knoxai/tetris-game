package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 480
	blockSize    = 20
	boardWidth   = 10
	boardHeight  = 20
)

type Game struct {
	board        [][]int
	currentPiece *Piece
	nextPiece    *Piece
	score        int
	gameOver     bool
	lastUpdate   time.Time
}

type Piece struct {
	shape [][]int
	x, y  int
}

var pieces = [][][]int{
	{{1, 1, 1, 1}}, // I
	{{1, 1}, {1, 1}}, // O
	{{1, 1, 1}, {0, 1, 0}}, // T
	{{1, 1, 0}, {0, 1, 1}}, // S
	{{0, 1, 1}, {1, 1, 0}}, // Z
	{{1, 1, 1}, {1, 0, 0}}, // L
	{{1, 1, 1}, {0, 0, 1}}, // J
}

func NewGame() *Game {
	g := &Game{
		board: make([][]int, boardHeight),
	}
	for i := range g.board {
		g.board[i] = make([]int, boardWidth)
	}
	g.spawnPiece()
	g.lastUpdate = time.Now()
	return g
}

func (g *Game) spawnPiece() {
	g.currentPiece = g.nextPiece
	if g.currentPiece == nil {
		g.currentPiece = &Piece{
			shape: pieces[rand.Intn(len(pieces))],
			x:     boardWidth/2 - 1,
			y:     0,
		}
	}
	g.nextPiece = &Piece{
		shape: pieces[rand.Intn(len(pieces))],
		x:     boardWidth/2 - 1,
		y:     0,
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.movePiece(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.movePiece(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.movePiece(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}

	if time.Since(g.lastUpdate) > 500*time.Millisecond {
		if !g.movePiece(0, 1) {
			g.placePiece()
			g.clearLines()
			g.spawnPiece()
			if g.checkCollision(0, 0) {
				g.gameOver = true
			}
		}
		g.lastUpdate = time.Now()
	}

	return nil
}

func (g *Game) movePiece(dx, dy int) bool {
	if !g.checkCollision(dx, dy) {
		g.currentPiece.x += dx
		g.currentPiece.y += dy
		return true
	}
	return false
}

func (g *Game) rotatePiece() {
	oldShape := g.currentPiece.shape
	newShape := make([][]int, len(oldShape[0]))
	for i := range newShape {
		newShape[i] = make([]int, len(oldShape))
	}
	for i := range oldShape {
		for j := range oldShape[i] {
			newShape[j][len(oldShape)-1-i] = oldShape[i][j]
		}
	}
	g.currentPiece.shape = newShape
	if g.checkCollision(0, 0) {
		g.currentPiece.shape = oldShape
	}
}

func (g *Game) checkCollision(dx, dy int) bool {
	for i := range g.currentPiece.shape {
		for j := range g.currentPiece.shape[i] {
			if g.currentPiece.shape[i][j] != 0 {
				x := g.currentPiece.x + j + dx
				y := g.currentPiece.y + i + dy
				if x < 0 || x >= boardWidth || y >= boardHeight || (y >= 0 && g.board[y][x] != 0) {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) placePiece() {
	for i := range g.currentPiece.shape {
		for j := range g.currentPiece.shape[i] {
			if g.currentPiece.shape[i][j] != 0 {
				x := g.currentPiece.x + j
				y := g.currentPiece.y + i
				g.board[y][x] = 1
			}
		}
	}
}

func (g *Game) clearLines() {
	linesCleared := 0
	for y := boardHeight - 1; y >= 0; y-- {
		full := true
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			g.board = append(g.board[:y], g.board[y+1:]...)
			g.board = append([][]int{make([]int, boardWidth)}, g.board...)
			linesCleared++
		}
	}
	g.score += linesCleared * 100
}

func (g *Game) Draw(screen *ebiten.Image) {
	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}

	// Draw the board
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] != 0 {
				ebitenutil.DrawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, white)
			}
		}
	}

	// Draw the current piece
	if g.currentPiece != nil {
		for i := range g.currentPiece.shape {
			for j := range g.currentPiece.shape[i] {
				if g.currentPiece.shape[i][j] != 0 {
					x := g.currentPiece.x + j
					y := g.currentPiece.y + i
					ebitenutil.DrawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, red)
				}
			}
		}
	}

	// Draw the score
	ebitenutil.DebugPrint(screen, "Score: "+string(g.score))

	// Draw game over message
	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}