package provider

import (
	"fmt"

	"github.com/kubex-ecosystem/kbx/types"
)

type GmailProvider struct{}
type OutlookProvider struct{}
type MicrosoftProvider struct{}

func (p *GmailProvider) Send(cfg *types.MailConnection, msg *types.Email) error {
	return fmt.Errorf("gmail provider: not implemented")
}

func (p *OutlookProvider) Send(cfg *types.MailConnection, msg *types.Email) error {
	return fmt.Errorf("outlook provider: not implemented")
}

func (p *MicrosoftProvider) Send(cfg *types.MailConnection, msg *types.Email) error {
	return fmt.Errorf("microsoft provider: not implemented")
}
