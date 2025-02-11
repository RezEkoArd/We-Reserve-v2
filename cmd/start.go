package cmd

import (
	"wereserve/internal/app"

	"github.com/spf13/cobra"
)

var startCMD = &cobra.Command{
	Use: "start",
	Short: "start",
	Long: `start`,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

func init() {
	rootCMD.AddCommand(startCMD)
}