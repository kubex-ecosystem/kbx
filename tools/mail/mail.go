// Package mail provides email sending functionality with multiple provider support and fallback mechanisms.
package mail

import (
	"errors"
	"strings"
	"time"

	"github.com/kubex-ecosystem/kbx/tools/mail/provider"
	"github.com/kubex-ecosystem/kbx/types"
)

var provMap = map[string]types.MailProvider{
	"gmail":     &provider.GmailProvider{},
	"outlook":   &provider.OutlookProvider{},
	"microsoft": &provider.MicrosoftProvider{},
	// "sendmail":  provider.SendmailProvider{},
}

// fallback order: Kubex-style chaos-first resiliency
var fbkOrder = []string{
	"gmail",
	"outlook",
	"microsoft",
	"sendmail",
}

func Send(cfg *types.SMTPConfig, msg *types.Email) error {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	primary := strings.ToLower(cfg.Provider)
	if p, ok := provMap[primary]; ok {
		if err := p.Send(cfg, msg); err == nil {
			return nil
		}
	}

	// fallback
	for _, name := range fbkOrder {
		if name == primary {
			continue
		}
		if p, ok := provMap[name]; ok {
			if err := p.Send(cfg, msg); err == nil {
				return nil
			}
		}
	}

	return errors.New("KBX-Mail: all providers failed")
}
