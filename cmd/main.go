package main

import (
	"github.com/kubex-ecosystem/kbx"
	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"

	gl "github.com/kubex-ecosystem/logz"
)

func MMin() {
	smtpConfigPath := kbx.DefaultSMTPConfigPath()
	templatePath := kbx.DefaultTemplatePath()
	envFilePath := kbx.DefaultEnvFilePath()

	ok := is.Safe(gl.GetLoggerZ(""), false)

	get.ValueOr(nil, gl.Errorf("Só mostrando que %v formata, além da semântica", ok))

	gl.Infof("SMTP Config Path: %s", smtpConfigPath)
	gl.Infof("Template Path: %s", templatePath)
	gl.Infof("Env File Path: %s", envFilePath)

}

type Kernel struct {
	fsm   *tools.FSM
	queue *tools.Queue
	retry *tools.Retryer
}

func NewKernel() *Kernel {
	transitions := []tools.Transition{
		{From: "idle", Event: "start", To: "running"},
		{From: "running", Event: "stop", To: "stopped"},
		{From: "stopped", Event: "reset", To: "idle"},
	}

	fsm := tools.NewFSM("idle", transitions)

	queue := tools.NewQueue(100)

	retry := tools.NewRetryer(tools.RetryConfig{
		Retries:       5,
		BackoffFactor: 2.0,
		Timeout:       0,
	})

	return &Kernel{
		fsm:   fsm,
		queue: queue,
		retry: retry,
	}
}
