package imap

import (
	"context"
	"fmt"
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	imapparser "github.com/kubex-ecosystem/kbx/tools/mail/imap"
)

// Config define parâmetros mínimos para acesso IMAP.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
	Mailbox  string
	Limit    int
}

// Message representa um email básico retornado do IMAP.
type Message struct {
	UID         uint32
	From        string
	Subject     string
	Body        string
	Attachments []imapparser.ParsedAttachment
}

// FetchUnread obtém mensagens não lidas da mailbox (default INBOX).
func FetchUnread(ctx context.Context, cfg Config) ([]*Message, error) {
	mailbox := cfg.Mailbox
	if mailbox == "" {
		mailbox = "INBOX"
	}
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var c *client.Client
	var err error
	if cfg.UseTLS {
		c, err = client.DialTLS(address, nil)
	} else {
		c, err = client.Dial(address)
	}
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	if err := c.Login(cfg.Username, cfg.Password); err != nil {
		return nil, err
	}

	if _, err = c.Select(mailbox, false); err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	uids, err := c.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return []*Message{}, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, section.FetchItem()}
	messages := make(chan *imap.Message, len(uids))

	go func() {
		_ = c.Fetch(seqset, items, messages)
	}()

	out := []*Message{}
	limit := cfg.Limit
	if limit <= 0 || limit > len(uids) {
		limit = len(uids)
	}

	for msg := range messages {
		if len(out) >= limit {
			break
		}
		body := ""
		if msg != nil {
			if r := msg.GetBody(section); r != nil {
				b, _ := io.ReadAll(r)
				body = string(b)
			}
		}
		attachments, _ := imapparser.ParseAttachments(msg)
		out = append(out, &Message{
			UID:         msg.Uid,
			From:        envelopeAddr(msg),
			Subject:     msg.Envelope.Subject,
			Body:        body,
			Attachments: attachments,
		})
	}

	return out, nil
}

func envelopeAddr(msg *imap.Message) string {
	if msg == nil || msg.Envelope == nil || len(msg.Envelope.From) == 0 {
		return ""
	}
	addr := msg.Envelope.From[0]
	name := addr.PersonalName
	email := addr.Address()
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
