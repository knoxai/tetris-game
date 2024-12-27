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
	screenWidth  = 320
	screenHeight = 600
	gridSize     = 25
	boardWidth   = 10
	boardHeight  = 20
)

type Game struct {
	board         [boardHeight][boardWidth]int
	currentPiece  *Piece
	score         int
	gameOver      bool
	tickCount     int        // 添加这一行
	lastMoveDown  time.Time
	moveDownDelay time.Duration
}

type Piece struct {
	x, y    int
	shape   [][]int
	current int
}

var tetrominos = [][][]int{
	{ // I
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{ // O
		{1, 1},
		{1, 1},
	},
	{ // T
		{0, 1, 0},
		{1, 1, 1},
		{0, 0, 0},
	},
	{ // L
		{1, 0, 0},
		{1, 1, 1},
		{0, 0, 0},
	},
	{ // J
		{0, 0, 1},
		{1, 1, 1},
		{0, 0, 0},
	},
	{ // S
		{0, 1, 1},
		{1, 1, 0},
		{0, 0, 0},
	},
	{ // Z
		{1, 1, 0},
		{0, 1, 1},
		{0, 0, 0},
	},
}

func NewGame() *Game {
	g := &Game{
		moveDownDelay: time.Millisecond * 500,
		lastMoveDown:  time.Now(),
	}
	g.spawnPiece()
	return g
}

func (g *Game) spawnPiece() {
	shapeIdx := rand.Intn(len(tetrominos))
	shape := tetrominos[shapeIdx]
	g.currentPiece = &Piece{
		x:     boardWidth/2 - len(shape[0])/2,
		y:     0,
		shape: shape,
	}
	
	if !g.isValidMove(g.currentPiece.x, g.currentPiece.y, g.currentPiece.shape) {
		g.gameOver = true
	}
}

func (g *Game) isValidMove(x, y int, shape [][]int) bool {
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] == 0 {
				continue
			}
			
			newX := x + j
			newY := y + i
			
			if newX < 0 || newX >= boardWidth || newY >= boardHeight {
				return false
			}
			
			if newY >= 0 && g.board[newY][newX] != 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			*g = *NewGame()
		}
		return nil
	}

	now := time.Now()
	if now.Sub(g.lastMoveDown) >= g.moveDownDelay {
		g.moveDown()
		g.lastMoveDown = now
	}

	// 添加按键检测的防抖
	if g.tickCount%5 == 0 {  // 每5帧才检测一次按键，避免旋转太快
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			if g.isValidMove(g.currentPiece.x-1, g.currentPiece.y, g.currentPiece.shape) {
				g.currentPiece.x--
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			if g.isValidMove(g.currentPiece.x+1, g.currentPiece.y, g.currentPiece.shape) {
				g.currentPiece.x++
			}
		}
		// 添加向上箭头旋转方块的控制
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.rotatePiece()
		}
	}

	// 下箭头可以连续按
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.moveDown()
	}

	// 空格键仍然保留作为备选旋转键
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.rotatePiece()
	}

	g.tickCount++
	return nil
}

func (g *Game) moveDown() {
	if g.isValidMove(g.currentPiece.x, g.currentPiece.y+1, g.currentPiece.shape) {
		g.currentPiece.y++
	} else {
		g.lockPiece()
		g.checkLines()
		g.spawnPiece()
	}
}

func (g *Game) lockPiece() {
	for i := 0; i < len(g.currentPiece.shape); i++ {
		for j := 0; j < len(g.currentPiece.shape[i]); j++ {
			if g.currentPiece.shape[i][j] != 0 {
				y := g.currentPiece.y + i
				if y >= 0 {
					g.board[y][g.currentPiece.x+j] = 1
				}
			}
		}
	}
}

func (g *Game) checkLines() {
	for i := boardHeight - 1; i >= 0; i-- {
		full := true
		for j := 0; j < boardWidth; j++ {
			if g.board[i][j] == 0 {
				full = false
				break
			}
		}
		if full {
			g.score += 100
			for k := i; k > 0; k-- {
				copy(g.board[k][:], g.board[k-1][:])
			}
			for j := 0; j < boardWidth; j++ {
				g.board[0][j] = 0
			}
			i++
		}
	}
}

func (g *Game) rotatePiece() {
	rows := len(g.currentPiece.shape)
	cols := len(g.currentPiece.shape[0])
	rotated := make([][]int, cols)
	for i := range rotated {
		rotated[i] = make([]int, rows)
	}
	
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			rotated[j][rows-1-i] = g.currentPiece.shape[i][j]
		}
	}
	
	if g.isValidMove(g.currentPiece.x, g.currentPiece.y, rotated) {
		g.currentPiece.shape = rotated
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw board
	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {
			if g.board[i][j] != 0 {
				ebitenutil.DrawRect(screen, float64(j*gridSize), float64(i*gridSize), float64(gridSize-1), float64(gridSize-1), color.RGBA{0, 255, 0, 255})
			}
		}
	}

	// Draw current piece
	if g.currentPiece != nil {
		for i := 0; i < len(g.currentPiece.shape); i++ {
			for j := 0; j < len(g.currentPiece.shape[i]); j++ {
				if g.currentPiece.shape[i][j] != 0 {
					x := float64((g.currentPiece.x + j) * gridSize)
					y := float64((g.currentPiece.y + i) * gridSize)
					ebitenutil.DrawRect(screen, x, y, float64(gridSize-1), float64(gridSize-1), color.RGBA{255, 0, 0, 255})
				}
			}
		}
	}

	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over! Press R to restart")
	} else {
		ebitenutil.DebugPrint(screen, "Score: "+string(rune(g.score)))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 600
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tetris")

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}