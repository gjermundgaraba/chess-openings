package cmd

import (
	"bufio"
	"chess-cli/storage"
	"fmt"
	"github.com/notnil/chess"
	"github.com/spf13/cobra"
	"os"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "chess",
		Short: "chess cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			selectedOpening, err := SelectOpening()
			if err != nil {
				return err
			}
			fmt.Println("Selected opening:", selectedOpening)

			variations, err := storage.GetVariations(selectedOpening)
			if err != nil {
				return err
			}

			// Create a new game
			game := chess.NewGame()

			PrintBoard(game, variations)

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
					storage.StoreVariation(game, selectedOpening)
					isMove = false
				}

				if isMove {
					err := game.MoveStr(move)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}

				PrintBoard(game, variations)
			}

			return nil
		},
	}

	// Add subcommands
	rootCmd.AddCommand(CreateOpeningCmd())
	rootCmd.AddCommand(GetOpeningCmd())

	return rootCmd
}

func PrintBoard(game *chess.Game, variations map[string][]string) {
	// Print the board and clear the terminal so the board is always at the top
	fmt.Print("\033[H\033[2J")
	fmt.Println(game.Position().Board().Draw())

	// Get the next moves
	nextMoves := variations[game.Position().String()]
	if len(nextMoves) != 0 {
		fmt.Println("Next moves:")
		for _, nextMove := range nextMoves {
			fmt.Println(nextMove)
		}
	} else {
		fmt.Println("No next moves")
	}
}
