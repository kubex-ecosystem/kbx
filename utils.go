package kbx

import (
	"io"
	"os"

	gl "github.com/kubex-ecosystem/logz"
)


func asyncCopyFile(src, dst string) error {
	//go func() {
	_, err := copyFile(src, dst)
	if err != nil {
		gl.Printf("Erro ao fazer backup do arquivo: %v\n", err)
	}
	//}()
	return nil
}
func copyFile(src, dst string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)

	return io.Copy(destination, source)
}
