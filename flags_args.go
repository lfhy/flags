package flags

var argsFlag = NewFlags()

// 解析
func Parse(args ...string) {
	argsFlag.Parse(args...)
}

// 获取额外的参数列表
func Args() []string {
	return argsFlag.Args()
}

// 获取解析的map[string]string
func Kvargs() map[string]string {
	return argsFlag.Kvargs()
}

// 解析为map[string]string
func ParseToKV(args ...string) map[string]string {
	return argsFlag.ParseToKV(args...)
}

// 定义解析的类型
func Var(fg ...any) {
	argsFlag.Var(fg...)
}

// 打印使用方法
func PrintUsage() {
	argsFlag.PrintUsage()
}

// 添加子命令
func AddSubCommand(sub string, initfn func(...string) error, runfn func() error) {
	var fn subFunc
	fn.initfn = initfn
	fn.runfn = runfn
	argsFlag.AddSubCommand(sub, &fn)
}

// 解析并运行
func ParseToRun(args ...string) error {
	return argsFlag.ParseToRun(args...)
}

// 直接运行
func Run() error {
	return argsFlag.Run()
}
