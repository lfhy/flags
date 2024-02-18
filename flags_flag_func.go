package flags

import "fmt"

type FlagFunc struct {
	// 执行
	cmd func() error
	// 名称
	name string
	// 备注
	dc string
}

func (ff *FlagFunc) Run() error {
	return ff.cmd()
}

func (ff *FlagFunc) PrintMark() {
	fmt.Printf("  %v\t\t%v\n", ff.name, ff.dc)
}
