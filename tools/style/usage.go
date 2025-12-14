package style

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Set usage definitions for the command and its subcommands

func SetUsageTemplate(cliCommand *cobra.Command) {
	setUsageDefinition(cliCommand)
	for _, c := range cliCommand.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}
}
