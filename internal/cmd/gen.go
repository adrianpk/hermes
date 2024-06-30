package cmd

import (
	"fmt"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate HTML from Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.GenerateHTML()
		if err != nil {
			fmt.Println("Error generating HTML:", err)
			return
		}

		fmt.Println("HTML generated.")
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
