// Package provider contains common utilities for mail providers
package provider

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"net/smtp"
	"strings"

	"github.com/kubex-ecosystem/kbx/types"
)

// formatAddr builds host:port
func formatAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// sendSMTPMessage normalizes RFC-style message assembly
func sendSMTPMessage(c *smtp.Client, msg *types.Email) error {
	if err := c.Mail(msg.From); err != nil {
		return err
	}
	for _, r := range msg.To {
		if err := c.Rcpt(r); err != nil {
			return err
		}
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	buf := bytes.NewBuffer(nil)
	writeRFC822(buf, msg)

	_, err = wc.Write(buf.Bytes())
	return err
}

// writeRFC822 – builds a *very clean* RFC message
func writeRFC822(w io.Writer, msg *types.Email) {
	boundary := "KBXMAIL-" + "BOUNDARY"

	fmt.Fprintf(w, "From: %s\r\n", msg.From)
	fmt.Fprintf(w, "To: %s\r\n", joinAddressList(msg.To))
	fmt.Fprintf(w, "Subject: %s\r\n", msg.Subject)
	fmt.Fprintf(w, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(w, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary)

	// TEXT
	if msg.Text != "" {
		fmt.Fprintf(w, "--%s\r\n", boundary)
		fmt.Fprintf(w, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
		fmt.Fprintf(w, "Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp := quotedprintable.NewWriter(w)
		qp.Write([]byte(msg.Text))
		qp.Close()
		fmt.Fprintf(w, "\r\n")
	}

	// HTML
	if msg.HTML != "" {
		fmt.Fprintf(w, "--%s\r\n", boundary)
		fmt.Fprintf(w, "Content-Type: text/html; charset=\"utf-8\"\r\n")
		fmt.Fprintf(w, "Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp := quotedprintable.NewWriter(w)
		qp.Write([]byte(msg.HTML))
		qp.Close()
		fmt.Fprintf(w, "\r\n")
	}

	// attachments
	for _, att := range msg.Attachments {
		fmt.Fprintf(w, "--%s\r\n", boundary)
		fmt.Fprintf(w, "Content-Type: %s\r\n", att.Mime)
		fmt.Fprintf(w, "Content-Disposition: attachment; filename=\"%s\"\r\n", att.Filename)
		fmt.Fprintf(w, "Content-Transfer-Encoding: base64\r\n\r\n")

		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(att.Data)))
		base64.StdEncoding.Encode(encoded, att.Data)
		w.Write(encoded)
		fmt.Fprintf(w, "\r\n")
	}

	fmt.Fprintf(w, "--%s--", boundary)
}

func joinAddressList(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0]
	}
	return strings.Join(list, ", ")
}
