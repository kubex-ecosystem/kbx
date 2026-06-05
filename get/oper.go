package get

import (
	"context"
	"os"

	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/logz"
)

// Mapper is an alias for tools.Mapper[T]
type Mapper[T any] = tools.Mapper[T]

// Loader creates a new Mapper[T] to load values from the specified source.
// Example:
//
//	type MyStruct struct {
//	    Name string
//	    Age  int
//	}
//	mapper := get.Loader[MyStruct]("env")
//
// Will create a Mapper[MyStruct] that
// can load MyStruct values from files and environment. Returns: *Mapper[MyStruct]
func Loader[T any](from string) *Mapper[T] { return tools.NewEmptyMapperType[T](from) }

// AvailablePort finds the next available port starting from basePort.
// basePort: the base port to start searching from.
// maxAttempts: the maximum number of attempts to find an available port.
// Returns the available port as a string, or an error if no available port is found.
func AvailablePort(ctx context.Context, basePort, maxAttempts int) (int, error) {
	ar := make([]int, maxAttempts)

	for z := range ar {

		// Avoiding well-known ports (< 1024) and ephemeral ports (> 49151)
		if basePort+z < 1024 || basePort+z > 49151 {
			logz.Noticef("Avoiding well-known port %d...", basePort+z)
			continue
		}

		if !is.PortFree(ctx, logz.Sprintf("%d", basePort+z)) {
			logz.Debugf("Port %d not free, trying the next one...", basePort+z)
			continue
		}

		logz.Debugf("Port %d free, using it.", basePort+z)
		return basePort + z, nil
	}

	return 0, logz.Errorf("no available port in range %d-%d", basePort, basePort+maxAttempts-1)
}

// Hostname retrieves the hostname of the machine.
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
