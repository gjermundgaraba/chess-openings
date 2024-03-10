package main

import (
	"bufio"
	"chess-cli/storage"
	"fmt"
	"github.com/notnil/chess"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	// Set up the root command
	rootCmd := &cobra.Command{
		Use:   "chess",
		Short: "chess cli",
		Run: func(cmd *cobra.Command, args []string) {
			// Create a new game
			game := chess.NewGame()

			// Print the board
			fmt.Print("\033[H\033[2J")
			fmt.Println(game.Position().Board().Draw())

			// Create a scanner to read the moves
			scanner := bufio.NewScanner(os.Stdin)

			// Loop through the moves
			for scanner.Scan() {
				// Get the move
				move := scanner.Text()

				isMove := true
				if move == "\x1b[D\n" {
					// Go back one move in the game
					isMove = false
				}
				if move == "save" {
					storage.StoreGame(game)
					isMove = false
				}

				if isMove {
					err := game.MoveStr(move)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}

				// Print the board and clear the terminal so the board is always at the top
				fmt.Print("\033[H\033[2J")
				fmt.Println(game.Position().Board().Draw())
			}
		},
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
