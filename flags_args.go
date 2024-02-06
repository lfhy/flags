package flags

var argsFlag = NewFlags()

func Parse(args ...string) {
	argsFlag.Parse(args...)
}

func Args() []string {
	return argsFlag.Args()
}

func Kvargs() map[string]string {
	return argsFlag.Kvargs()
}

func ParseToKV(args ...string) map[string]string {
	return argsFlag.ParseToKV(args...)
}

func Var(fg ...any) {
	argsFlag.Var(fg...)
}
