package get

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kubex-ecosystem/kbx/internal/module/kbx"
	gl "github.com/kubex-ecosystem/logz"
	"github.com/spf13/viper"
)

// ValidateWorkerLimit validates the worker limit value.
// It returns an error if the value is not an integer or is a negative integer.
func ValidateWorkerLimit(value any) error {
	if limit, ok := value.(int); ok {
		if limit < 0 {
			return gl.Errorf("worker limit cannot be negative")
		}
	} else {
		return gl.Errorf("invalid type for worker limit")
	}
	return nil
}

// BootID retrieves the system's unique boot ID from /proc/sys/kernel/random/boot_id.
func BootID() (string, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// BootTimeMac retrieves the system boot time on macOS/BSD using the sysctl command.
func BootTimeMac() (string, error) {
	cmd := exec.Command("sysctl", "-n", "kern.boottime")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// BootTimeWindows retrieves the system boot time on Windows using PowerShell.
func BootTimeWindows() (string, error) {
	cmd := exec.Command("powershell", "-Command", "(Get-WmiObject Win32_OperatingSystem).LastBootUpTime")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// EnvOrDefault retrieves an environment variable or a viper configuration value by key.
// It falls back to defaultValue if neither is set.
func EnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultValue
}

// DefaultConfigPath resolves and retrieves the default configuration path for KUBEX GNyx.
// If the resolved path is empty, invalid, or is ".", it falls back to the user's home directory config or a temporary fallback directory.
func DefaultConfigPath() (string, error) {
	var err error
	vprFile := filepath.Dir(viper.ConfigFileUsed())
	if strings.Contains(vprFile, "gnyx") {
		vprFile = filepath.Dir(vprFile)
	}

	configPath := os.ExpandEnv(EnvOr("KUBEX_GNYX_CONFIG_PATH", ValOrType(vprFile, kbx.DefaultGNyxConfigPath)))
	if strings.TrimSpace(configPath) == "" || configPath == "." {
		configPath, err = os.UserHomeDir()
		if err != nil {
			gl.Log("error", fmt.Sprintf("Failed to get user home directory: %v", err))
			return FallbackTempDir()
		}
		configPath = filepath.Join(configPath, "config.json")
	}

	realPath := configPath
	if filepath.Base(realPath) != "config.json" {
		realPath = filepath.Join(realPath, "config.json")
	}

	if err = os.MkdirAll(filepath.Dir(realPath), 0o755); err != nil {
		gl.Log("error", fmt.Sprintf("Failed to create directory %s: %v", filepath.Dir(realPath), err))
		return FallbackTempDir()
	}

	return configPath, nil
}

// FallbackTempDir creates and returns a temporary fallback directory path, logging warning/fatal events.
func FallbackTempDir() (string, error) {
	base := os.TempDir()
	tmpDir, err := os.MkdirTemp(base, "kubex_gnyx_")
	if err != nil {
		gl.Log("fatal", fmt.Sprintf("Failed to create temp dir for fallback: %v", err))
		return "", err
	}
	gl.Log("warn", fmt.Sprintf("Using temporary directory for config fallback: %s", tmpDir))
	return tmpDir, nil
}
