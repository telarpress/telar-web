package cmd

import (
	"fmt"

	"github.com/red-gold/telar-core/utils"
	"github.com/spf13/cobra"
)

func MakeEmail() *cobra.Command {
	var command = &cobra.Command{
		Use:          "email",
		Short:        "Send email",
		Example:      `  tsocial email --e email.example.com --p password`,
		SilenceUsage: false,
	}

	command.Flags().String("e", "", "Reference email")
	command.Flags().String("p", "", "Email password")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		refEmail, _ := command.Flags().GetString("e")
		password, _ := command.Flags().GetString("p")
		sendEmail(refEmail, password)
		return nil
	}
	return command
}

func sendEmail(refEmail string, password string) {
	email := utils.NewEmail(refEmail, password, "smtp.gmail.com")
	req := utils.NewEmailRequest([]string{"amir.gholzam@live.com"}, "test subject", "")
	templateData := struct {
		Name    string
		Code    string
		AppName string
	}{
		Name:    "Amir",
		Code:    "1293740",
		AppName: "Playground",
	}
	status, err := email.SendEmail(req, "tmpl.html", templateData)
	if err != nil {
		fmt.Printf("Error in email: %s", err)
	}

	fmt.Printf("Email status is %v", status)
}
