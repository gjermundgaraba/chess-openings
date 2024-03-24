package cmd

import (
	"chess-cli/storage"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func CreateOpeningCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-opening",
		Short: "Create a new opening in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			var openingName string
			if err := survey.AskOne(&survey.Input{Message: "Opening name"}, &openingName, survey.WithValidator(survey.Required)); err != nil {
				return err
			}

			if err := storage.CreateOpening(openingName); err != nil {
				return err
			}

			fmt.Println("Opening created successfully:", openingName)
			return nil
		},
	}
}
