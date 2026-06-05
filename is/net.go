package is

import (
	"context"
	"net"
	"strings"
	"time"
)

// IsIPv6 checks if an IP address is an IPv6 address.
// ip: the IP address to check.
// Returns true if the IP address is an IPv6 address, otherwise false.
func IsIPv6(ip string) bool {
	var sanitizedIP string
	if strings.Contains(ip, ":") {
		sanitizedIP = strings.Replace(ip, ":", "", -1)

		// If the sanitized IP is still the same as the original, then it's an IPv6
		if sanitizedIP == ip {
			return true
		}
	}
	return false
}

// PortOpen checks if a port is open.
// port: the port to be checked.
// Returns true if the port is open, otherwise false. Returns an error, if any.
func PortFree(ctx context.Context, port string) bool {
	lc := net.ListenConfig{KeepAlive: 10 * time.Second}
	ln, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		// Se deu erro (já em uso, permissão negada, etc), não está livre.
		return false
	}
	_ = ln.Close()
	return true
}

// PortUsed checks if a port is in use.
// port: the port to be checked.
// Returns true if the port is in use, otherwise false.
func PortUsed(ctx context.Context, port string) bool { return !PortFree(ctx, port) }

// LocalIP attempts to find a local IP address that is not a loopback address.
// It tries to connect to a public DNS server (Google's DNS on port 53) to discover the
// primary network interface and its IP address. If it fails to find a non-loopback address,
// it returns an empty string.
func LocalIP() string {
	// A more robust way to get the local IP, avoiding loopback addresses
	// by trying to connect to a known external address.
	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		return ""
	}
	defer func() { _ = conn.Close() }()

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	return localAddr.IP.String()
}
