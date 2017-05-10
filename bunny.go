package main

import (
	"fmt"
	"os/exec"
	"time"
)

type bunny struct {
	id        string
	rprt      chan report
	cmd       string
	arguments []string
	slp       time.Duration
}

func (b *bunny) start() {
	timer := time.NewTimer(b.slp)

	for {
		select {
		case <-timer.C:
			err := b.runbun()
			if err != nil {
				return
			}
			timer.Reset(b.slp)
		}
	}
}

func (b *bunny) runbun() error {
	msg, err := exec.Command(b.cmd, b.arguments...).Output()
	b.report(string(msg), err)

	return err
}

func (b *bunny) report(msg string, err error) {
	if err == nil {
		b.rprt <- report{
			id: b.id,
			msg: msg,
		}
		return
	}

	ee, ok := err.(*exec.ExitError)
	if !ok {
		b.rprt <- report{
			id: b.id,
			err: fmt.Errorf("unknown error: %v", err),
			msg: msg,
		}
		return
	}

	b.rprt <- report{
		id: b.id,
		err: fmt.Errorf(string(ee.Stderr)),
		msg: msg,
	}
}
