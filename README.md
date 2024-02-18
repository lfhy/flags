# flags
一个更优雅的flag包。

# Feature
- 全局参数
```golang
import (
    "fmt"
    "github.com/lfhy/flags"
)
var Port string
var Socks bool
flags.Var(&flags.Flag{Name:"port",Default:"123",&A},flags.Flag{Name:"ss",Default:"false",&B})
flags.Parse()
fmt.Printf("运行如下参数时会解析为如下\n-port=%v -ss=%v",Port,Socks)
fmt.Println("参数列表为:",flags.Args())
```
```sh
go run demo.go -port=233 hello -ss
运行如下参数时会解析为如下
-port=233 -ss=true
参数列表为:[hello]
```
- 支持传入结构体构建
```golang
import (
    "fmt"
    "github.com/lfhy/flags"
)
// 定义结构体
type MyStruct struct{
    Port string `flag:"port" default:"1234"`
    Debug bool `flag:"debug" default:"false"`
}
var demoStruct MyStruct
// 传入flags进行解析
flags.Var(&demoStruct)
flags.Parse()
fmt.Printf("%+v",demoStruct)
```
```sh
go run demo.go -port=233 -debug
# 运行时候可以看到输出信息
{Port:233 Debug:true}
```
- 多种定义方式
```golang
import (
    "fmt"
    "github.com/lfhy/flags"
)
var A,B string
flags.Var(&flags.Flag{Name:"a",Default:"123",&A})
flags.Flag{Name:"b",Default:"321",&B}.Var()
flags.Parse()
fmt.Println("A:",A)
fmt.Println("B:",B)
```
```sh
go run demo.go -a=333
# 运行时可以看到
A:333
B:321
```
- 支持导出kvargs
```golang
import (
    "fmt"
    "github.com/lfhy/flags"
)
var Port string
var Socks bool
flags.Var(&flags.Flag{Name:"port",Default:"123",&A},flags.Flag{Name:"ss",Default:"false",&B})
flags.Parse()
fmt.Println("参数集合为:",flags.Kvargs())
fmt.Println("参数列表为:",flags.Args())
```
```sh
go run demo.go -port=233 hello -ss
参数集合为:map[-port:233 -ss:true]
参数列表为:[hello]
```
- 支持注册子命令并执行
```golang
type Config struct {
	Port string `flag:"port" default:"1234"`
	IP   string `flag:"bind" default:"127.0.0.1"`
}

var config Config

type Server struct {
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
```
```sh
go run demo.go -port=233 server -ss
# 运行时候可以看到输出信息
Server Listen:127.0.0.1:233
```
# TODO
- 兼容旧flag包函数
- 支持从配置文件环境变量导入参数
- 错误处理
- 输出帮助信息