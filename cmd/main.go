package main

import (
	"github.com/kubex-ecosystem/kbx"
	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/tools/flow/fsm"
	"github.com/kubex-ecosystem/kbx/tools/flow/queue"

	gl "github.com/kubex-ecosystem/logz"
)

// MMin is the main function of the application.
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

// Kernel is a simple example of a kernel with a finite state machine, a queue, and a retryer.
type Kernel struct {
	fsm   *fsm.FSM
	queue *queue.Queue
	retry *tools.Retryer
}

// NewKernel creates a new Kernel with the given initial state and transitions.
func NewKernel() *Kernel {
	transitions := []fsm.Transition{
		{From: "idle", Event: "start", To: "running"},
		{From: "running", Event: "stop", To: "stopped"},
		{From: "stopped", Event: "reset", To: "idle"},
	}

	fsm := fsm.NewFSM("idle", transitions)

	queue := queue.NewQueue(100)

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
