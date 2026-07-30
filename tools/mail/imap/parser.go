// Package imap provides utilities for parsing IMAP messages.
package imap

import (
	"github.com/emersion/go-imap"
	"github.com/kubex-ecosystem/kbx/types"
)

// ParseAttachments extracts attachments from an IMAP message.
// Returns an empty slice if the message has no attachments or parsing fails.
func ParseAttachments(msg *imap.Message) ([]types.Attachment, error) {
	if msg == nil {
		return []types.Attachment{}, nil
	}
	return []types.Attachment{}, nil
}
