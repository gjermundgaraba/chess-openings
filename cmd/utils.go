package cmd

import (
	"chess-cli/storage"
	"github.com/AlecAivazis/survey/v2"
	"log"
)

func SelectOpening() (string, error) {
	openings, err := storage.GetOpenings()
	if err != nil {
		return "", err
	}

	var selectedIndex int
	if err := survey.AskOne(&survey.Select{
		Message: "Select opening",
		Options: openings,
	}, &selectedIndex, survey.WithValidator(survey.Required)); err != nil {
		log.Fatalf("Error selecting: %v", err)
	}

	selectedOpening := openings[selectedIndex]

	return selectedOpening, nil
}
