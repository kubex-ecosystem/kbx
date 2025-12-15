package types

import (
	"fmt"
	"net/url"
	"strings"
)

type DSN struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	Params   map[string]string
}

func ParseDSN(input string) (*DSN, error) {
	if strings.Contains(input, "://") {
		return parseURLDSN(input)
	}
	return parseKVDSN(input)
}

func parseURLDSN(raw string) (*DSN, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	dsn := &DSN{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Port:   u.Port(),
		DBName: strings.TrimPrefix(u.Path, "/"),
		Params: map[string]string{},
	}

	if u.User != nil {
		dsn.User = u.User.Username()
		dsn.Password, _ = u.User.Password()
	}

	for k, v := range u.Query() {
		if len(v) > 0 {
			dsn.Params[k] = v[0]
		}
	}

	return dsn, nil
}

func parseKVDSN(raw string) (*DSN, error) {
	dsn := &DSN{
		Params: map[string]string{},
	}

	parts := strings.Fields(raw)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		val := strings.Trim(kv[1], `"'`)

		switch key {
		case "host":
			dsn.Host = val
		case "port":
			dsn.Port = val
		case "user":
			dsn.User = val
		case "password":
			dsn.Password = val
		case "dbname":
			dsn.DBName = val
		default:
			dsn.Params[key] = val
		}
	}

	return dsn, nil
}

func (d *DSN) ToURL() string {
	u := &url.URL{
		Scheme: d.Scheme,
		Host:   d.Host,
		Path:   "/" + d.DBName,
	}

	if d.Port != "" {
		u.Host += ":" + d.Port
	}

	if d.User != "" {
		if d.Password != "" {
			u.User = url.UserPassword(d.User, d.Password)
		} else {
			u.User = url.User(d.User)
		}
	}

	q := url.Values{}
	for k, v := range d.Params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (d *DSN) ToKV() string {
	var parts []string

	if d.Host != "" {
		parts = append(parts, "host="+d.Host)
	}
	if d.Port != "" {
		parts = append(parts, "port="+d.Port)
	}
	if d.User != "" {
		parts = append(parts, "user="+d.User)
	}
	if d.Password != "" {
		parts = append(parts, "password="+d.Password)
	}
	if d.DBName != "" {
		parts = append(parts, "dbname="+d.DBName)
	}

	for k, v := range d.Params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(parts, " ")
}
