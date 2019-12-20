package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func MakeTime() *cobra.Command {
	var command = &cobra.Command{
		Use:          "time",
		Short:        "Show time",
		Example:      `  tsocial time --now --utc`,
		SilenceUsage: false,
	}

	command.Flags().Bool("now", false, "Time now")
	command.Flags().Bool("utc", false, "Time UTC")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		now, _ := command.Flags().GetBool("now")
		utc, _ := command.Flags().GetBool("utc")
		showTime(now, utc)
		return nil
	}
	return command
}

func showTime(now bool, utc bool) {
	if now {
		fmt.Printf("Time now : %s \n", time.Now())
	}

	if now {
		fmt.Printf("Time utc now : %v \n", (time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))))
	}
}
