package cmd

import (
	"chess-cli/storage"
	"github.com/spf13/cobra"
)

func GetOpeningCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-opening",
		Short: "Get opening",
		RunE: func(cmd *cobra.Command, args []string) error {
			selectedOpening, err := SelectOpening()
			if err != nil {
				return err
			}

			positionsFromOpening, err := storage.GetVariations(selectedOpening)
			if err != nil {
				return err
			}

			for fen, moves := range positionsFromOpening {
				for _, move := range moves {
					cmd.Printf("FEN: %s, Move: %s\n", fen, move)
				}
			}

			return nil
		},
	}
}
