package main

import (
	"fmt"

	"github.com/lfhy/flags"
)

type Config struct {
	Port string `flag:"port" default:"1234" dc:"启动的端口"`
	IP   string `flag:"bind" default:"127.0.0.1" dc:"绑定的IP地址"`
}

func (c Config) ToAddr() string {
	return fmt.Sprintf("%v:%v", c.IP, c.Port)
}

var config Config

func main() {
	flags.Var(&config)
	flags.Parse()
	flags.PrintUsage()
	fmt.Printf("解析后地址: %v\n", config.ToAddr())
}
