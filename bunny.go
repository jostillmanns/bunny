package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type recur func() recur

type bunny struct {
	id          string
	cmd         exec.Cmd
	rprt        chan report
	slp         time.Duration
	changestate chan state
	st          state
}

func (b *bunny) pause() {
	b.changestate <- paused
}

func (b *bunny) resume() {
	b.changestate <- running
}

func (b *bunny) checkRemoteControl() {
	select {
	case st := <-b.changestate:
		if st == running {
			return
		}

		<-b.changestate
	}

	if b.st == running {
		return
	}

	// state changed internally
	for st := range b.changestate {
		if st == paused {
			continue
		}

		b.st = running
	}
}

func (b *bunny) run() error {
	_, err := b.cmd.Output()
	if err == nil {
		return nil
	}

	return err
}

func (b *bunny) setup(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	for next := b.initiate(); ; next = next() {
	}
}

func (b *bunny) initiate() recur {
	b.checkRemoteControl()

	err := b.run()
	if err != nil {
		return b.reportError(err)
	}

	time.Sleep(b.slp)
	return b.initiate
}

func (b *bunny) reportError(err error) recur {
	ee, ok := err.(*exec.ExitError)
	if !ok {
		b.rprt <- report{
			b.id,
			fmt.Sprintf("unknown error: %v", err),
		}
		return b.initiate
	}

	b.rprt <- report{
		b.id,
		string(ee.Stderr),
	}

	b.st = paused

	return b.initiate
}
