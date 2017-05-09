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
	_, err := exec.Command(b.cmd, b.arguments...).Output()
	if err != nil {
		b.report(err)
	}

	return err
}

func (b *bunny) report(err error) {
	ee, ok := err.(*exec.ExitError)
	if !ok {
		b.rprt <- report{
			b.id,
			fmt.Sprintf("unknown error: %v", err),
		}
	}

	b.rprt <- report{
		b.id,
		string(ee.Stderr),
	}
}
