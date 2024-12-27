package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	boardWidth  = 10
    boardHeight = 20
    fallSpeed = 500
)

type Game struct {
	Board    [][]int
	CurrentTetromino *Tetromino
	gameOver bool
	score int
	lastFallTime time.Time
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	board := make([][]int, boardHeight)
	for i := range board {
		board[i] = make([]int, boardWidth)
	}
    
	return &Game{
        Board:    board,
		CurrentTetromino: NewTetromino(),
		gameOver: false,
		score: 0,
		lastFallTime: time.Now(),
    }
}


func (game *Game) MoveLeft() {
    game.CurrentTetromino.X -= 1
	if game.IsCollision() {
        game.CurrentTetromino.X += 1
    }
}

func (game *Game) MoveRight() {
	game.CurrentTetromino.X += 1
	if game.IsCollision() {
        game.CurrentTetromino.X -= 1
	}
}

func (game *Game) MoveDown() {
	game.CurrentTetromino.Y += 1
	if game.IsCollision() {
        game.CurrentTetromino.Y -= 1 // 撞击之后退回到原来的位置
		game.FreezeTetromino()
		game.ClearLines()
		game.CurrentTetromino = NewTetromino()

		if game.IsCollision() {
			game.gameOver = true
		}
	}
}



func (game *Game) RotateRight() {
	oldShape := game.CurrentTetromino.Shape
	game.CurrentTetromino.RotateRight()
	if game.IsCollision() {
		game.CurrentTetromino.Shape = oldShape
		game.CurrentTetromino.RotateLeft()
	}
}

func (game *Game) RotateLeft() {
	oldShape := game.CurrentTetromino.Shape
	game.CurrentTetromino.RotateLeft()
	if game.IsCollision() {
		game.CurrentTetromino.Shape = oldShape
		game.CurrentTetromino.RotateRight()
	}
}

func (game *Game) IsCollision() bool {
	for y, row := range game.CurrentTetromino.Shape {
		for x, cell := range row {
			if cell == 1 {
				boardX := game.CurrentTetromino.X + x
				boardY := game.CurrentTetromino.Y + y
				
                //check out of boundry
				if boardX < 0 || boardX >= boardWidth || boardY >= boardHeight{
					return true
				}

				if boardY >= 0 && game.Board[boardY][boardX] == 1 {
					return true
				}

			}
		}
	}
	return false
}


func (game *Game) FreezeTetromino() {
	for y, row := range game.CurrentTetromino.Shape {
		for x, cell := range row {
			if cell == 1 {
                gameY := game.CurrentTetromino.Y + y
				gameX := game.CurrentTetromino.X + x
				if gameY >= 0 {
                	game.Board[gameY][gameX] = 1
				}
			}
		}
	}
}

func (game *Game) ClearLines() {
	for y := 0; y < boardHeight; y++ {
		fullLine := true // 一行默认是满的
        for x := 0; x < boardWidth; x++ {
           if game.Board[y][x] == 0 {
			   fullLine=false
			   break
		   }
        }

        if fullLine {
			//删除该行，并将上面的所有行往下移动一行
			game.score = game.score + 100
			for i:=y; i> 0; i-- {
				game.Board[i] = game.Board[i-1]
			}
			game.Board[0] = make([]int, boardWidth)
			y--
		}
	}
}

func (game *Game) Drawboard() {
	ClearScreen()
    fmt.Println("Score: ", game.score)
	//Copy board for render
	tempBoard := make([][]int, boardHeight);
	for i := range tempBoard {
		tempBoard[i] = make([]int, boardWidth)
		copy(tempBoard[i], game.Board[i])
	}

	for y, row := range game.CurrentTetromino.Shape {
		for x, cell := range row {
			if cell == 1 {
				boardX := game.CurrentTetromino.X + x
				boardY := game.CurrentTetromino.Y + y
				if boardY >= 0 && boardY < boardHeight {
					tempBoard[boardY][boardX]=2
				}
			}
		}
	}


	for _,row := range tempBoard{
		for _, cell := range row {
            if cell == 0 {
                fmt.Print(". ")    
            } else if cell == 1 {
                fmt.Print("# ")
            } else {
                fmt.Print("* ")
            }
		}
		fmt.Println()
	}
	fmt.Println()
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout // 将输出重定向到控制台
	cmd.Run()             // 执行命令
}

func (game *Game) GameTick() {
	if time.Since(game.lastFallTime).Milliseconds() >= fallSpeed {
		game.MoveDown()
		game.lastFallTime = time.Now()
	}
	game.Drawboard()
}

func (game *Game) IsGameOver() bool {
	return game.gameOver
}
