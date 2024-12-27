package main

import "math/rand"

// 定义俄罗斯方块的类型
type TetrominoType int

const (
	I TetrominoType = iota
	J
	L
	O
	S
	T
	Z
)


type Rotation int

const (
	R0	Rotation = iota
	R90
	R180
	R270
)

// 定义俄罗斯方块的结构
type Tetromino struct {
	Type     TetrominoType
	Rotation Rotation
	Shape    [][]int // 俄罗斯方块的形状
    X, Y int
}

// 定义每种俄罗斯方块的初始形状
var tetrominoShapes = map[TetrominoType][][][]int{
    I: {
        {{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}},
        {{0, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 0, 0}},
		{{0, 0, 0, 0}, {0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}},
		{{0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}},
    },
	J: {
        {{1, 0, 0}, {1, 1, 1}, {0, 0, 0}},
		{{0, 1, 1}, {0, 1, 0}, {0, 1, 0}},
		{{0, 0, 0}, {1, 1, 1}, {0, 0, 1}},
		{{0, 1, 0}, {0, 1, 0}, {1, 1, 0}},
    },
    L: {
        {{0, 0, 1}, {1, 1, 1}, {0, 0, 0}},
		{{0, 1, 0}, {0, 1, 0}, {0, 1, 1}},
		{{0, 0, 0}, {1, 1, 1}, {1, 0, 0}},
		{{1, 1, 0}, {0, 1, 0}, {0, 1, 0}},
    },
    O: {
        {{0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
    },
    S: {
        {{0, 1, 1}, {1, 1, 0}, {0, 0, 0}},
		{{0, 1, 0}, {0, 1, 1}, {0, 0, 1}},
		{{0, 0, 0}, {0, 1, 1}, {1, 1, 0}},
		{{1, 0, 0}, {1, 1, 0}, {0, 1, 0}},
    },
    T: {
        {{0, 1, 0}, {1, 1, 1}, {0, 0, 0}},
		{{0, 1, 0}, {0, 1, 1}, {0, 1, 0}},
		{{0, 0, 0}, {1, 1, 1}, {0, 1, 0}},
		{{0, 1, 0}, {1, 1, 0}, {0, 1, 0}},
    },
    Z: {
        {{1, 1, 0}, {0, 1, 1}, {0, 0, 0}},
		{{0, 0, 1}, {0, 1, 1}, {0, 1, 0}},
		{{0, 0, 0}, {1, 1, 0}, {0, 1, 1}},
		{{0, 1, 0}, {1, 1, 0}, {1, 0, 0}},
    },
}

func NewTetromino() *Tetromino {
    tetrominoType := TetrominoType(rand.Intn(7)) // 随机生成一个 0-6 的数来决定方块类型
    return &Tetromino{
            Type: tetrominoType,
			Rotation: R0,
            Shape: tetrominoShapes[tetrominoType][R0], 
			X: 3, // 初始位置
        	Y: 0,
        }
}

func (t *Tetromino) RotateRight() {
	newRotation  := (t.Rotation + 1) % 4;
	t.Rotation = newRotation
	t.Shape = tetrominoShapes[t.Type][newRotation]
}

func (t *Tetromino) RotateLeft() {
	newRotation  := (t.Rotation + 3) % 4;
	t.Rotation = newRotation
	t.Shape = tetrominoShapes[t.Type][newRotation]
}