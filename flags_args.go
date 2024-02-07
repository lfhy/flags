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
