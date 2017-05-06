package main

import (
	"context"
	"log"
)

func main() {
}

type userinput int

const (
	resume userinput = iota
	stop
)

type state int

const (
	running state = iota
	paused
)

type hole struct {
	input userinput
	buns  []bunny
	rprt  chan report

	cancel context.CancelFunc
	remote chan struct{}
}

func (h *hole) orchestrate() {
	ctx, cancelfunc := context.WithCancel(context.Background())
	h.cancel = cancelfunc

	for _, e := range h.buns {
		go e.setup(ctx)
	}

	for rprt := range h.rprt {
		for _, e := range h.buns {
			e.pause()
		}

		rprt.communicate()
		<-h.remote

		for _, e := range h.buns {
			e.resume()
		}
	}
}

type report struct {
	id  string
	msg string
}

func (rprt *report) communicate() {
	log.Println(rprt.msg)
}
