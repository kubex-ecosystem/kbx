package provider

import (
	"bytes"
	"os/exec"
	"reflect"

	"github.com/kubex-ecosystem/kbx/load"
	"github.com/kubex-ecosystem/kbx/types"
)

type SendmailProviderImpl struct {
	cfgMap map[reflect.Type]*load.MailConnection
}

type SendmailProvider interface {
	Send(_ *types.MailConnection, msg *types.Email) error
}

func NewProvider[T SendmailProvider](cfgFilePath string) (SendmailProvider, error) {
	// Load the SMTP config from the specified file
	mailConfig, err := load.LoadConfigOrDefault[load.MailConfig](cfgFilePath, true)
	if err != nil {
		return nil, err
	}
	empty := SendmailProviderImpl{
		cfgMap: make(map[reflect.Type]*load.MailConnection),
	}
	for _, conn := range mailConfig.Connections {
		if conn == nil {
			continue
		}
		if conn.Protocol == "smtp" || conn.Protocol == "" {
			empty.cfgMap[reflect.TypeOf(*conn)] = conn
			return empty, nil
		}
	}

	return empty, nil
}

func (s SendmailProviderImpl) Send(_ *types.MailConnection, msg *types.Email) error {
	cmd := exec.Command("/usr/sbin/sendmail", "-t", "-i")

	buf := new(bytes.Buffer)
	writeRFC822(buf, msg)

	cmd.Stdin = buf
	return cmd.Run()
}
