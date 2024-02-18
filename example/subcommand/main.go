package main

import (
	"fmt"

	"github.com/lfhy/flags"
)

type Server struct {
	Port string `flag:"port" default:"1234"`
	IP   string `flag:"bind" default:"127.0.0.1"`
}

func (Server) CmdName() string {
	return "server"
}

func (s *Server) CmdRun() error {
	fmt.Printf("Server Listen:%v:%v\n", s.IP, s.Port)
	return nil
}

func main() {
	var server Server
	flags.Var(&server)
	err := flags.ParseToRun()
	if err != nil {
		panic(err)
	}
}
