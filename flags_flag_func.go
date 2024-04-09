package flags

import "fmt"

type FlagFunc struct {
	// 初始化
	initfunc func(...string) error
	// 执行
	cmd func() error
	// 名称
	name string
	// 备注
	dc string
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
