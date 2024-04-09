package main

import (
	"fmt"

	"github.com/lfhy/flags"
)

type Config struct {
	Port string `flag:"port" default:"1234"`
	IP   string `flag:"bind" default:"127.0.0.1"`
}

var config Config

type Server struct {
}

func (Server) CmdInit(args ...string) error {
	fmt.Println("运行参数:", args)
	return nil
}

func (Server) CmdName() string {
	return "server"
}

func (Server) CmdRun() error {
	fmt.Printf("Server Listen:%v:%v\n", config.IP, config.Port)
	return nil
}

type Client struct {
}

func (Client) CmdName() string {
	return "client"
}

func (Client) CmdRun() error {
	fmt.Printf("Client Connect To:%v:%v\n", config.IP, config.Port)
	return nil
}

func main() {
	var server Server
	var client Client
	flags.Var(&config, &server, &client)
	err := flags.ParseToRun()
	if err != nil {
		panic(err)
	}
}
