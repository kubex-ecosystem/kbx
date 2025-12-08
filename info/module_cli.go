// Package module provides internal types and functions for the KubexModuleCLIImpl application.
package info

import (
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/kubex-ecosystem/kbx"
	gl "github.com/kubex-ecosystem/logz"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type KubexModuleCLIImpl struct {
	*kbx.GlobalRef
	*kbx.ParamsImpl

	*cobra.Command

	Enable      bool
	PrintBanner bool
	Banners     []string

	alias       string
	shortDesc   string
	longDesc    string
	usage       string
	annotations map[string]string
	examples    []string

	parentCmdName string
}

func RunKbx() KubexModuleCLI {
	var printBannerV = os.Getenv("GROMPT_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &KubexModuleCLIImpl{
		PrintBanner: strings.ToLower(printBannerV) == "true",
	}
}

func (m *KubexModuleCLIImpl) ID() uuid.UUID            { return m.GlobalRef.ID }
func (m *KubexModuleCLIImpl) Module() string           { return m.GlobalRef.Name }
func (m *KubexModuleCLIImpl) Alias() string            { return m.alias }
func (m *KubexModuleCLIImpl) ShortDescription() string { return m.shortDesc }
func (m *KubexModuleCLIImpl) LongDescription() string  { return m.longDesc }
func (m *KubexModuleCLIImpl) Usage() string            { return m.usage }
func (m *KubexModuleCLIImpl) Examples() []string       { return m.examples }
func (m *KubexModuleCLIImpl) Active() bool             { return m.Enable }
func (m *KubexModuleCLIImpl) Execute() error           { return m.command().Execute() }

func (m *KubexModuleCLIImpl) command() *cobra.Command {
	gl.Debugf("Starting %s CLI...", m.Module())

	// Build the root command
	err := m.buildRootCommand()
	if err != nil {
		return nil
	}

	// Add more commands as needed
	m.AddCommand(CliCommand())

	// Set usage definitions for the command and its subcommands
	m.setUsageDefinition(m.Command)
	for _, c := range m.Command.Commands() {
		m.setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return m.Command
}

func (m *KubexModuleCLIImpl) SetParentCmdName(parentCmdName string) {
	m.parentCmdName = parentCmdName
}

func (m *KubexModuleCLIImpl) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}

func (m *KubexModuleCLIImpl) buildRootCommand() error {
	if m.Command == nil {
		m.Command = &cobra.Command{
			Use:   m.Module(),
			Short: m.ShortDescription(),
			Long:  m.LongDescription(),
			Run: func(cmd *cobra.Command, args []string) {
				// Default action when no subcommand is provided
				cmd.Help()
			},
		}
	} else {
		m.Command.Use = m.Module()
		m.Command.Short = m.ShortDescription()
		m.Command.Long = m.LongDescription()
		m.Command.Example = m.concatenateExamples()
		m.Command.Aliases = []string{m.Alias()}
	}
	return nil
}

func (m *KubexModuleCLIImpl) colorYellow(s string) string {
	return color.New(color.FgYellow).SprintFunc()(s)
}
func (m *KubexModuleCLIImpl) colorGreen(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}
func (m *KubexModuleCLIImpl) colorBlue(s string) string {
	return color.New(color.FgBlue).SprintFunc()(s)
}
func (m *KubexModuleCLIImpl) colorRed(s string) string { return color.New(color.FgRed).SprintFunc()(s) }
func (m *KubexModuleCLIImpl) colorHelp(s string) string {
	return color.New(color.FgCyan).SprintFunc()(s)
}

func (m *KubexModuleCLIImpl) hasServiceCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] == "true" {
			return true
		}
	}
	return false
}
func (m *KubexModuleCLIImpl) hasModuleCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] != "true" {
			return true
		}
	}
	return false
}
func (m *KubexModuleCLIImpl) setUsageDefinition(cmd *cobra.Command) {
	cobra.AddTemplateFunc("colorYellow", m.colorYellow)
	cobra.AddTemplateFunc("colorGreen", m.colorGreen)
	cobra.AddTemplateFunc("colorRed", m.colorRed)
	cobra.AddTemplateFunc("colorBlue", m.colorBlue)
	cobra.AddTemplateFunc("colorHelp", m.colorHelp)
	cobra.AddTemplateFunc("hasServiceCommands", m.hasServiceCommands)
	cobra.AddTemplateFunc("hasModuleCommands", m.hasModuleCommands)

	// Altera o template de uso do cobra
	cmd.SetUsageTemplate(m.cliTemplate())
}
func (m *KubexModuleCLIImpl) cliTemplate() string {
	var cliUsageTemplate = `{{- if index .Annotations "banner" }}{{colorBlue (index .Annotations "banner")}}{{end}}{{- if (index .Annotations "description") }}
{{index .Annotations "description"}}
{{- end }}

{{colorYellow "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [args]{{end}}{{if gt (len .Aliases) 0}}

{{colorYellow "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{colorYellow "Example:"}}
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{colorYellow "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{colorGreen (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{colorYellow "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasAvailableInheritedFlags}}

{{colorYellow "Global Options:"}}
  {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasHelpSubCommands}}

{{colorYellow "Additional help topics:"}}
{{range .Commands}}{{if .IsHelpCommand}}
  {{colorGreen (rpad .CommandPath .CommandPathPadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasSubCommands}}

{{colorYellow (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}
`
	return cliUsageTemplate
}
