// Package imap provides functions to connect to an IMAP server.
package imap

import (
	"github.com/emersion/go-imap/client"
)

func Connect(host string, user string, pass string) (*client.Client, error) {
	c, err := client.DialTLS(host, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Login(user, pass); err != nil {
		return nil, err
	}
	return c, nil
}
