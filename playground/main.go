package main

import (
	"github.com/red-gold/telar-web/playground/pkg/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdStack := cmd.MakeStack()
	cmdEmail := cmd.MakeEmail()
	cmdTime := cmd.MakeTime()

	var rootCmd = &cobra.Command{
		Use: "tsocial",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmdStack)
	rootCmd.AddCommand(cmdEmail)
	rootCmd.AddCommand(cmdTime)

	rootCmd.Execute()

}
