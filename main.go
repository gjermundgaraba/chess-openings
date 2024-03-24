package main

import (
	"chess-cli/cmd"
	"fmt"
	"os"
)

func main() {
	rootCmd := cmd.RootCmd()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
