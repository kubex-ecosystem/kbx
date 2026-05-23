// Package module provides internal types and functions for the GoBE application.
package module

import (
	"github.com/kubex-ecosystem/kbx/internal/module/version"

	"github.com/spf13/cobra"
)

type Kbx struct {
	parentCmdName string
	hideBanner    bool
	certPath      string
	keyPath       string
	configPath    string
	Banners       []string
}

func (m *Kbx) Alias() string {
	return ""
}
func (m *Kbx) ShortDescription() string {
	return "Domus: GKBX Database and Docker manager/service. "
}
func (m *Kbx) LongDescription() string {
	return `Domus: Is a tool to manage GKBX database and Docker services.`
}
func (m *Kbx) Usage() string {
	return "domus [command] [args]"
}
func (m *Kbx) Examples() []string {
	return []string{"domus [command] [args]", "domus database user auth'", "domus db roles list"}
}
func (m *Kbx) Active() bool {
	return true
}
func (m *Kbx) Module() string {
	return "domus"
}
func (m *Kbx) Execute() error {
	dbChanData := make(chan interface{})
	defer close(dbChanData)

	if spyderErr := m.Command().Execute(); spyderErr != nil {
		return spyderErr
	}
	return nil
}
func (m *Kbx) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     m.Module(),
		Example: m.concatenateExamples(),
		Version: version.GetVersion(),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(version.CliCommand())

	return cmd
}

func (m *Kbx) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *Kbx) concatenateExamples() string {
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
