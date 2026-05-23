package main

import (
	"fmt"

	"github.com/kubex-ecosystem/kbx"
)

// MMin is the main function of the application.
func MMin() {
	smtpConfigPath := kbx.DefaultSMTPConfigPath()
	templatePath := kbx.DefaultTemplatePath()
	envFilePath := kbx.DefaultEnvFilePath()

	fmt.Printf("SMTP Config Path: %s\n", smtpConfigPath)
	fmt.Printf("Template Path: %s\n", templatePath)
	fmt.Printf("Env File Path: %s\n", envFilePath)
}
