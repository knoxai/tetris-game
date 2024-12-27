package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	game := NewGame()
    input := bufio.NewReader(os.Stdin)

	for !game.IsGameOver() {
        game.GameTick()
        fmt.Print("Enter Command(a:left,d:right,s:down,l:rotate left,r:rotate right, x exit):")
		text, _ := input.ReadString('\n')

        switch text {
            case "a\n":
                game.MoveLeft()
            case "d\n":
                game.MoveRight()
            case "s\n":
                game.MoveDown()
            case "l\n":
                game.RotateLeft()
			case "r\n":
				game.RotateRight()
			case "x\n":
				os.Exit(0)
         }
	}
    fmt.Println("Game Over! Your final score :", game.score)
}