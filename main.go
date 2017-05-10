package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("read: %v", err)
	}

	config := new(config)
	err = json.Unmarshal(configFile, config)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	rprt := make(chan report)

	buns := make(map[string]*bunny)
	for i := range config.Cmds {
		buns[config.Cmds[i].Id] = &bunny{
			id:        config.Cmds[i].Id,
			rprt:      rprt,
			cmd:       config.Cmds[i].Cmd,
			arguments: config.Cmds[i].Arguments,
			slp:       time.Duration(config.Cmds[i].Delta) * time.Second,
		}
	}

	bunnyhole := &hole{
		buns: buns,
		rprt: rprt,
	}

	bunnyhole.orchestrate()
}

type config struct {
	Cmds []cmd `json:"cmds"`
}

type cmd struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"arguments"`
	Id        string   `json:"id"`
	Delta     int      `json:"delta"`
}

type hole struct {
	buns map[string]*bunny
	rprt chan report
}

func (h *hole) remote() {
	// wait for user input
	_, _ = os.Stdin.Read(make([]byte, 1))
}

func (h *hole) orchestrate() {
	for i := range h.buns {
		go h.buns[i].start()
	}

	for rprt := range h.rprt {
		rprt.communicate()
		if rprt.err == nil {
			continue
		}

		h.remote()
		go h.buns[rprt.id].start()
	}
}

type report struct {
	id  string
	msg string
	err error
}

func (rprt *report) communicate() {
	log.Println(rprt.msg)

	if rprt.err != nil {
		log.Println("error", rprt.err)
	}
}
