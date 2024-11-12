package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/adrianpk/gohermes/internal/cmd"
)

const ver = "0.0.1"

//go:embed layout
var layoutFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd(ver)
	rootCmd.AddCommand(cmd.NewInitCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewGenCmd())
	rootCmd.AddCommand(cmd.NewUpgradeCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewNewCmd())
	rootCmd.AddCommand(cmd.NewPublishCmd()) 
	rootCmd.AddCommand(cmd.NewBackupCmd()) 

	if len(os.Args) > 1 && os.Args[1] == "help" {
		rootCmd.Usage()
		return
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
