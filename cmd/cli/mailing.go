package cli

import (
	"github.com/kubex-ecosystem/kbx"
	"github.com/spf13/cobra"
)

func MailCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mail",
		Short: "GNyx Mail - Send test emails using configured SMTP settings",
		Long: `GNyx Mail provides a CLI to send test emails using the configured SMTP settings.

Features:
  • Send test emails using configured SMTP settings
`,
	}

	rootCmd.AddCommand(sendCommand())

	return rootCmd
}

func sendCommand() *cobra.Command {
	var debug = false
	var mailParams = kbx.NewMailSrvParams("")

	// Make HTTP GET request to health endpoint

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a test email using the configured SMTP settings",
		Run: func(cmd *cobra.Command, args []string) {
			_ = debug
			// sender := kbx.NewMailSender(mailParams)
			// ...
		},
	}

	sendCmd.Flags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")
	sendCmd.Flags().StringVar(&mailParams.ConfigPath, "smtp-config", "", "Path to SMTP configuration file")
	sendCmd.Flags().StringVar(&mailParams.From, "from-email", "", "From email address")
	sendCmd.Flags().StringVar(&mailParams.Name, "from-name", "", "From name")

	return sendCmd
}
