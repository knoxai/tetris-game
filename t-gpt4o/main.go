package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"math/rand"
	"time"
)

const (
	screenWidth  = 200
	screenHeight = 400
	blockSize    = 20
	rows         = screenHeight / blockSize
	cols         = screenWidth / blockSize
	initialSpeed = 30 // 控制下落速度的初始值
)

var (
	colors = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // Black
		color.RGBA{0xff, 0x00, 0x00, 0xff}, // Red
		color.RGBA{0x00, 0xff, 0x00, 0xff}, // Green
		color.RGBA{0x00, 0x00, 0xff, 0xff}, // Blue
		color.RGBA{0xff, 0xff, 0x00, 0xff}, // Yellow
		color.RGBA{0xff, 0x00, 0xff, 0xff}, // Magenta
		color.RGBA{0x00, 0xff, 0xff, 0xff}, // Cyan
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // White
	}
)

type Game struct {
	board        [rows][cols]int
	currentPiece *Piece
	tick         int
	speed        int
}

type Piece struct {
	x, y   int
	shape  [][]int
	color  int
}

func NewPiece() *Piece {
	shapes := [][][]int{
		{{1, 1, 1, 1}}, // I
		{{1, 1}, {1, 1}}, // O
		{{0, 1, 0}, {1, 1, 1}}, // T
		{{1, 1, 0}, {0, 1, 1}}, // S
		{{0, 1, 1}, {1, 1, 0}}, // Z
		{{1, 0, 0}, {1, 1, 1}}, // L
		{{0, 0, 1}, {1, 1, 1}}, // J
	}
	shape := shapes[rand.Intn(len(shapes))]
	color := rand.Intn(len(colors)-1) + 1
	return &Piece{x: cols/2 - len(shape[0])/2, y: 0, shape: shape, color: color}
}

func (g *Game) Update() error {
	if g.currentPiece == nil {
		g.currentPiece = NewPiece()
	}

	// 控制方块移动
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.currentPiece.x--
		if g.collides() {
			g.currentPiece.x++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.currentPiece.x++
		if g.collides() {
			g.currentPiece.x--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.currentPiece.y++
		if g.collides() {
			g.currentPiece.y--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}

	// 控制方块下落速度
	g.tick++
	if g.tick >= g.speed {
		g.tick = 0
		g.currentPiece.y++
		if g.collides() {
			g.currentPiece.y--
			g.mergePiece()
			g.clearLines()
			g.currentPiece = NewPiece()
		}
	}

	return nil
}

func (g *Game) rotatePiece() {
	originalShape := g.currentPiece.shape
	rotatedShape := make([][]int, len(originalShape[0]))
	for i := range rotatedShape {
		rotatedShape[i] = make([]int, len(originalShape))
	}
	for y, row := range originalShape {
		for x, cell := range row {
			rotatedShape[x][len(originalShape)-1-y] = cell
		}
	}
	g.currentPiece.shape = rotatedShape
	if g.collides() {
		g.currentPiece.shape = originalShape
	}
}

func (g *Game) collides() bool {
	for y, row := range g.currentPiece.shape {
		for x, cell := range row {
			if cell != 0 {
				boardX := g.currentPiece.x + x
				boardY := g.currentPiece.y + y
				if boardX < 0 || boardX >= cols || boardY >= rows || g.board[boardY][boardX] != 0 {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) mergePiece() {
	for y, row := range g.currentPiece.shape {
		for x, cell := range row {
			if cell != 0 {
				g.board[g.currentPiece.y+y][g.currentPiece.x+x] = g.currentPiece.color
			}
		}
	}
}

func (g *Game) clearLines() {
	for y := 0; y < rows; y++ {
		full := true
		for x := 0; x < cols; x++ {
			if g.board[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			copy(g.board[1:y+1], g.board[0:y])
			for x := 0; x < cols; x++ {
				g.board[0][x] = 0
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colors[0])
	for y, row := range g.board {
		for x, cell := range row {
			if cell != 0 {
				ebitenutil.DrawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, colors[cell])
			}
		}
	}
	for y, row := range g.currentPiece.shape {
		for x, cell := range row {
			if cell != 0 {
				ebitenutil.DrawRect(screen, float64((g.currentPiece.x+x)*blockSize), float64((g.currentPiece.y+y)*blockSize), blockSize, blockSize, colors[g.currentPiece.color])
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tetris")
	game := &Game{speed: initialSpeed}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}