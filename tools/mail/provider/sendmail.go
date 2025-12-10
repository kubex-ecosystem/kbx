package provider

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/types"
)

type SendmailProviderImpl struct {
	cfgMap map[reflect.Type]*types.MailConfig
}

type SendmailProvider interface {
	Send(_ *types.MailConfig, msg *types.Email) error
}

func NewProvider[T SendmailProvider](cfgFilePath string) (SendmailProvider, error) {
	// Load the SMTP config from the specified file
	cfgMapper := tools.NewEmptyMapperType[types.MailConfig](cfgFilePath)
	smtpConfig, err := cfgMapper.DeserializeFromFile(filepath.Ext(cfgFilePath)[1:])
	if err != nil {
		var empty T
		return empty, err
	}
	empty := any(&SendmailProviderImpl{cfgMap: map[reflect.Type]*types.MailConfig{
		reflect.TypeOf(types.MailConfig{}): smtpConfig,
	}}).(T)
	return empty, nil
}

func (s SendmailProviderImpl) Send(_ *types.MailConfig, msg *types.Email) error {
	cmd := exec.Command("/usr/sbin/sendmail", "-t", "-i")

	buf := new(bytes.Buffer)
	writeRFC822(buf, msg)

	cmd.Stdin = buf
	return cmd.Run()
}
