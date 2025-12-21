// Package module provides internal types and functions for the GoBE application.
package module

import (
	"github.com/kubex-ecosystem/kbx/internal/module/version"

	kbxInfo "github.com/kubex-ecosystem/kbx/tools/info"
	kbxStyle "github.com/kubex-ecosystem/kbx/tools/style"
	logz "github.com/kubex-ecosystem/logz"

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
	return "CanalizeDS: GKBX Database and Docker manager/service. "
}
func (m *Kbx) LongDescription() string {
	return `CanalizeDS: Is a tool to manage GKBX database and Docker services. It provides many DB flavors like MySQL, PostgreSQL, MongoDB, Redis, etc. It also provides Docker services like Docker Swarm, Docker Compose, etc. It is a command line tool that can be used to manage GKBX database and Docker services.`
}
func (m *Kbx) Usage() string {
	return "canalizeds [command] [args]"
}
func (m *Kbx) Examples() []string {
	return []string{"canalizeds [command] [args]", "canalizeds database user auth'", "canalizeds db roles list"}
}
func (m *Kbx) Active() bool {
	return true
}
func (m *Kbx) Module() string {
	return "canalizeds"
}
func (m *Kbx) Execute() error {
	dbChanData := make(chan interface{})
	defer close(dbChanData)

	if spyderErr := m.Command().Execute(); spyderErr != nil {
		logz.Log("error", spyderErr.Error())
		return spyderErr
	} else {
		return nil
	}
}
func (m *Kbx) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: m.Module(),
		//Aliases:     []string{m.Alias(), "w", "wb", "webServer", "http"},
		Example: m.concatenateExamples(),
		Annotations: kbxInfo.CLIBannerStyle(
			m.Banners,
			[]string{
				m.LongDescription(),
				m.ShortDescription(),
			},
			m.hideBanner,
		),
		Version: version.GetVersion(),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(version.CliCommand())

	kbxStyle.SetUsageTemplate(cmd)

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
