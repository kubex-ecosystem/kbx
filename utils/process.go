// Package utils provides utility functions for the kbx package.
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/kubex-ecosystem/kbx/get"

	gl "github.com/kubex-ecosystem/logz"
)

// GenerateProcessFileName generates a unique process file name based on the process name,
// PID, and the system boot ID.
func GenerateProcessFileName(processName string, pid int) string {
	bootID, err := get.BootID()
	if err != nil {
		gl.Log("error", fmt.Sprintf("Failed to get boot ID: %v", err))
		return ""
	}
	return fmt.Sprintf("%s_%d_%s.pid", processName, pid, bootID)
}

// CreateProcessFile creates a temporary process file, writing the process name, PID,
// and current UNIX timestamp to it.
func CreateProcessFile(processName string, pid int) (*os.File, error) {
	fileName := GenerateProcessFileName(processName, pid)
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	// Escrever os detalhes do processo no arquivo
	_, err = file.WriteString(fmt.Sprintf("Process Name: %s\nPID: %d\nTimestamp: %d\n", processName, pid, time.Now().Unix()))
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

// RemoveProcessFile closes and deletes the given process file, logging the result.
func RemoveProcessFile(file *os.File) {
	if file == nil {
		return
	}

	fileName := file.Name()
	file.Close()

	// Apagar o arquivo temporário
	if err := os.Remove(fileName); err != nil {
		gl.Log("error", fmt.Sprintf("Failed to remove process file %s: %v", fileName, err))
	} else {
		gl.Log("debug", fmt.Sprintf("Successfully removed process file: %s", fileName))
	}
}
