package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"codeberg.org/anaseto/goal"
	gos "codeberg.org/anaseto/goal/os"
	"github.com/eiannone/keyboard"
)

//go:embed snake.goal
var gameGoalSource string

func main() {
	ctx := goal.NewContext()
	ctx.Log = os.Stderr
	gos.Import(ctx, "")

	_, err := ctx.Eval(gameGoalSource)
	if err != nil {
		fmt.Printf("Error evaluating Goal game source: %v\n", err)
	}

	os.Exit(runGame(ctx))
}

func runGame(ctx *goal.Context) int {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return 1
	}
	defer keyboard.Close()

	//nolint:mnd // Ticker duration
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	//nolint:mnd // Buffer for key events
	keyEvents, err := keyboard.GetKeys(10)
	if err != nil {
		fmt.Println("Failed to get keys:", err)
		return 1
	}

	_, err = ctx.Eval(`draw""`)
	if err != nil {
		fmt.Printf("Error drawing initial game: %v\n", err)
		return 1
	}

	for {
		select {
		case <-tick.C:
		case event := <-keyEvents:
			if event.Err != nil {
				fmt.Println("Keyboard error:", event.Err)
				return 1
			}
			if event.Rune == 'q' {
				fmt.Println("Bye!")
				return 0
			}
			if event.Rune == 'n' {
				reset(ctx)
			}
			//nolint:exhaustive // Keys other than these should be ignored.
			switch event.Key {
			case keyboard.KeyArrowUp:
				update(ctx, "Up")
			case keyboard.KeyArrowRight:
				update(ctx, "Right")
			case keyboard.KeyArrowDown:
				update(ctx, "Down")
			case keyboard.KeyArrowLeft:
				update(ctx, "Left")
			case keyboard.KeyEsc:
				fmt.Println("Later!")
				return 0
			}
		}
	}
}

func update(ctx *goal.Context, keyPress string) {
	_, err := ctx.Eval(fmt.Sprintf(`game.KeyPress::"%s"; update""; draw""`, keyPress))
	if err != nil {
		fmt.Printf("Error updating key press '%s': %v\n", keyPress, err)
	}
}

func reset(ctx *goal.Context) {
	_, err := ctx.Eval(`reset""; draw""`)
	if err != nil {
		fmt.Printf("Error resetting the game: %v\n", err)
	}
}
