package imap

import (
    "github.com/emersion/go-imap"
    "github.com/emersion/go-message/mail"
    "io"
)

type ParsedAttachment struct {
    Filename string 
    Data     []byte
}

func ParseAttachments(msg *imap.Message) ([]ParsedAttachment, error) {
    r := msg.GetBody(&imap.BodySectionName{})
    if r == nil {
        return nil, nil
    }

    mr, err := mail.CreateReader(r)
    if err != nil {
        return nil, err
    }

    var result []ParsedAttachment

    for {
        p, err := mr.NextPart()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        switch h := p.Header.(type) {
        case *mail.AttachmentHeader:
            filename, _ := h.Filename()
            data, _ := io.ReadAll(p.Body)
            result = append(result, ParsedAttachment{
                Filename: filename,
                Data:     data,
            })
        }
    }

    return result, nil
}