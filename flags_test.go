package flags_test

import (
	"fmt"
	"testing"

	"github.com/lfhy/flags"
)

type Struct struct {
	A string `flag:"sa" default:"321"`
	B int    `flag:"sb" default:"321"`
	C bool   `flag:"sc" default:"true"`
}

func (Struct) CmdName() string {
	return "aaa"
}

func (s *Struct) CmdRun() error {
	fmt.Println("Hello!")
	return nil
}

func TestParse(t *testing.T) {
	args := []string{
		"-a", "123", // 测试无等号
		"aaa",      // 测试非参数
		"-b=321",   // 测试等号
		"ddd",      // 测试参数顺序
		"-c",       // 测试Bool类型
		"-31=2222", // 测试Int类型
		"-3=",      // 测试String置空
		"-sa=233",  // 测试结构体赋值
	}

	type myString string
	var a string
	var b int
	var c bool
	var d string
	var e myString
	var f int
	var f2 string
	var g string

	var h Struct
	flags.Flag{
		Name:    "a",
		Default: "333",
		Value:   &a,
	}.Var()
	flags.Flag{
		Name:    "b",
		Default: "333",
		Value:   &b,
	}.Var()
	flags.Flag{
		Name:    "c",
		Default: "false",
		Value:   &c,
	}.Var()
	flags.Flag{
		Name:    "d",
		Default: "333",
		Value:   &d,
	}.Var()
	flags.Flag{
		Name:    "e",
		Default: "3334",
		Value:   &e,
	}.Var()
	flags.Var(
		&flags.Flag{
			Name:    "31",
			Default: "233",
			Value:   &f,
		},
		flags.Flag{
			Name:    "3",
			Default: "233",
			Value:   &g,
		},
		flags.Flag{
			Name:    "31",
			Default: "233",
			Value:   &f2,
		},
		&h,
	)
	flags.AddSubCommand("bbb", func(s ...string) error {
		return nil
	}, func() error {
		fmt.Println("Hello bbb!")
		return nil
	})
	flags.Parse(args...)
	fmt.Printf("a: %v\n", a)
	fmt.Printf("b: %v\n", b)
	fmt.Printf("c: %v\n", c)
	fmt.Printf("d: %v\n", d)
	fmt.Printf("e: %v\n", e)
	fmt.Printf("f: %v\n", f)
	fmt.Printf("f2: %v\n", f2)
	fmt.Printf("g: %v\n", g)
	fmt.Printf("h: %+v\n", h)
	kvargs := flags.Kvargs()
	fmt.Printf("res: %+v\n", kvargs)
	fargs := flags.Args()
	fmt.Printf("args: %v\n", fargs)
	flags.Run()
}
