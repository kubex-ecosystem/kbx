package utils

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	logz "github.com/kubex-ecosystem/logz"
)

// CheckPortOpen verifica se uma porta está aberta.
// port: a porta a ser verificada.
// Retorna true se a porta estiver aberta, caso contrário, false. Retorna um erro, se houver.
func CheckPortOpen(port string) (bool, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false, err // Porta já está em uso ou bloqueada
	}
	_ = ln.Close()
	return true, nil // Porta disponível
}

// ListOpenPorts lista todas as portas abertas no sistema.
// Executa o comando `netstat` para obter as portas abertas.
// Retorna uma lista de strings com as portas abertas e um erro, se houver.
func ListOpenPorts() ([]string, error) {
	cmd := exec.Command("sh", "-c", "netstat -tuln | grep LISTEN | awk '{print $4}' | sed 's/.*://' | sort -n | uniq")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	ports := strings.Split(strings.TrimSpace(string(output)), "\n")
	return ports, nil
}

// ClosePort fecha uma porta específica.
// port: a porta a ser fechada.
// Retorna um erro, se houver.
func ClosePort(port string) error {
	cmd := exec.Command("fuser", "-k", port+"/tcp")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao fechar a porta %s: %v", port, err)
	}
	return nil
}

// OpenPort abre uma porta específica.
// port: a porta a ser aberta.
// Retorna um erro, se houver.
func OpenPort(port string) error {
	cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", port, "-j", "ACCEPT")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao abrir a porta %s: %v", port, err)
	}
	return nil
}

// IsIPv6 verifica se um endereço IP é um endereço IPv6.
// ip: o endereço IP a ser verificado.
// Retorna true se o endereço IP for um endereço IPv6, caso contrário, false.
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

// GetAvailablePort returns the first available port in the range.
// startPort: the first port to check.
// endPort: the last port to check.
// Returns the first available port and an error, if any.
func GetAvailablePort(startPort int, endPort int) (string, error) {
	for i := startPort; i <= endPort; i++ {
		port := strconv.Itoa(i)
		if isAvailable, _ := CheckPortOpen(port); isAvailable {
			return port, nil
		}
	}
	return "", fmt.Errorf("no available ports found in range %d-%d", startPort, endPort)
}

// PortInUse verifica se uma porta específica está em uso.
// port: a porta a ser verificada.
// Retorna true se a porta estiver em uso, caso contrário, false.
func PortInUse(port string) bool {
	_, err := net.Dial("tcp", ":"+port)
	return err == nil
}

// GetNextPort returns the next available port starting from the given port.
func GetNextPort(port string) string {
	startPort, err := strconv.Atoi(port)
	if err != nil {
		logz.Error("Failed to convert port to integer: " + err.Error())
		return ""
	}
	logz.Info("Searching for available port starting from " + strconv.Itoa(startPort))
	for i := startPort; i <= 65535; i++ {
		port := strconv.Itoa(i + 1)
		if !PortInUse(port) {
			logz.Info("Available port found: " + port)
			return port
		}
	}
	logz.Error("No available ports found in range " + strconv.Itoa(startPort) + "-" + strconv.Itoa(65535))
	return ""
}
