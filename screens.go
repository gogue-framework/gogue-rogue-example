package main

import (
	"fmt"
	"github.com/gogue-framework/gogue/ui"
	"os"
)

type TitleScreen struct {}
func (ts *TitleScreen) Enter() {}
func (ts *TitleScreen) Exit() {}
func (ts *TitleScreen) UseEcs() bool { return false }

func (ts *TitleScreen) Render() {
	// Clear the entire window
	ui.ClearWindow(windowWidth, windowHeight, 0)
	centerX := windowWidth / 2
	centerY := windowHeight / 2

	title := "Gogue Roguelike"
	instruction := "Press {Up Arrow} to begin! Or Press {ESC} to exit"

	ui.PrintText(centerX - len(title) / 2, centerY, 0, 0, title, "", "", 0)
	ui.PrintText(centerX - len(instruction) / 2, centerY + 2, 0, 0, instruction, "", "", 0)
}

func (ts *TitleScreen) HandleInput() {
	key := ui.ReadInput(false)

	if key == ui.KeyEscape || key == ui.KeyClose {
		os.Exit(0)
	}

	if key == ui.KeyUp {
		// Change the screen to the play screen
		err := screenManager.SetScreenByName("game")

		if err != nil {
			// If something goes wrong switching screens, log the error and abort
			fmt.Print(err)
			return
		}

		// Immediately call the render method of the new current screen, otherwise, the existing screen will stay
		// rendered until the user provides input
		screenManager.CurrentScreen.Render()
	}
}

type GameScreen struct {}
func (gs *GameScreen) Enter() {}
func (gs *GameScreen) Exit() {}
func (gs *GameScreen) UseEcs() bool { return true }

func (gs *GameScreen) Render() {}

func (gs *GameScreen) HandleInput() {
	key := ui.ReadInput(false)

	if key == ui.KeyEscape || key == ui.KeyClose {
		os.Exit(0)
	}

	if key == ui.KeyDown {
		// Change the screen to the title screen
		err := screenManager.SetScreenByName("title")

		if err != nil {
			// If something goes wrong switching screens, log the error and abort
			fmt.Print(err)
			return
		}

		// Immediately call the render method of the new current screen, otherwise, the existing screen will stay
		// rendered until the user provides input
		screenManager.CurrentScreen.Render()
	}
}
