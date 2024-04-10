package flags

import "fmt"

type FlagFunc struct {
	// 初始化
	initfunc func(...string) error
	// 执行
	cmd func() error
	// 名称
	name string
	// 别名
	alias []string
	// 备注
	dc string
	// 参数列表
	flags []*Flag
}

func (ff *FlagFunc) Run(args ...string) error {
	err := ff.initfunc(args...)
	if err != nil {
		return err
	}
	return ff.cmd()
}

func (ff *FlagFunc) PrintMark() {
	fmt.Printf("  %v\t\t%v\n", ff.name, ff.dc)
}

// 大于参数帮助信息
func (ff *FlagFunc) PrintHelp() {
	fmt.Printf("子命令:%v", ff.name+" 使用方法\n\n")
	if len(ff.flags) == 0 {
		fmt.Printf("  暂无绑定参数\n")
		return
	}
	fmt.Println("参数列表:")
	for _, f := range ff.flags {
		f.PrintDefault()
	}
}
