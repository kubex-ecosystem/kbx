package imap

import "encoding/xml"

func ParseXMLAttachment[T any](data []byte) (*T, error) {
    var v T
    if err := xml.Unmarshal(data, &v); err != nil {
        return nil, err
    }
    return &v, nil
}