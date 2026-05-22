package main

import (
	"context"

	"github.com/kubex-ecosystem/kbx"
	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/tools/flow"
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

	// Demonstração do Kernel Atômico
	kernel := NewKernel()
	gl.Infof("Kernel iniciado no estado: %d", kernel.fsm.Current())

	if kernel.fsm.Trigger(EvStart) {
		gl.Infof("Kernel transicionado para: %d", kernel.fsm.Current())
	}
}

// Estados e Eventos do Kernel (Atômicos)
const (
	StateIdle fsm.State = 1 << iota
	StateRunning
	StateStopped
)

const (
	EvStart fsm.Event = iota
	EvStop
	EvReset
)

// Kernel is a simple example of a kernel with an atomic state machine, a queue, and a retryer.
type Kernel struct {
	fsm   *fsm.FSM
	queue *queue.Queue
	retry *tools.Retryer
	flow  *flow.Flow
}

// NewKernel creates a new Kernel with the given initial state and transitions.
func NewKernel() *Kernel {
	transitions := []fsm.Transition{
		{From: StateIdle, Event: EvStart, To: StateRunning},
		{From: StateRunning, Event: EvStop, To: StateStopped},
		{From: StateStopped, Event: EvReset, To: StateIdle},
	}

	fsm := fsm.NewFSM(StateIdle, transitions)
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
		flow:  flow.NewFlow(),
	}
}

// Execute demonstra a submissão de uma tarefa atômica
func (k *Kernel) Execute(ctx context.Context, task func(context.Context) error) error {
	return k.flow.Submit(ctx, task)
}
