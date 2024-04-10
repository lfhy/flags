package flags

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type FlagSet struct {
	args            []string
	kvargs          map[string]string
	all             map[string][]*Flag
	ok              map[string][]*Flag
	subcmd          map[string]*FlagFunc
	otherArgs       []string
	ishelp          bool
	helpfn          *func()
	isPrintArgsList bool
}

type FlagSubCmd interface {
	// 子命令名称
	CmdName() string
	// 子命令运行
	CmdRun() error
}

type FlagSubInit interface {
	// 子命令初始化
	CmdInit(...string) error
}

func defaultInit(_ ...string) error {
	return nil
}

type FlagSubMark interface {
	// 子命令备注
	CmdMark() string
}

type FlagSubAlias interface {
	// 子命令别名
	CmdAlias() []string
}

func NewFlags() *FlagSet {
	var f FlagSet
	f.kvargs = make(map[string]string)
	f.all = make(map[string][]*Flag)
	f.ok = make(map[string][]*Flag)
	f.subcmd = make(map[string]*FlagFunc)
	return &f
}

func (f *FlagSet) appendToAll(name string, flag *Flag) {
	if f.all[name] == nil {
		f.all[name] = make([]*Flag, 0)
	}
	f.all[name] = append(f.all[name], flag)
}

func (f *FlagSet) appendToOK(name string, flag *Flag) {
	if f.ok[name] == nil {
		f.ok[name] = make([]*Flag, 0)
	}
	f.ok[name] = append(f.ok[name], flag)
}

func (f *FlagSet) SetHelpFunc(fn func()) {
	f.helpfn = &fn
}

func (f *FlagSet) SetHelpPrintArgsList(isPrint bool) {
	f.isPrintArgsList = isPrint
}

func (f *FlagSet) Var(flags ...any) {
	for _, f2 := range flags {
		switch flag := f2.(type) {
		case *Flag:
			f.appendToAll(flag.Name, flag)
		case Flag:
			f.appendToAll(flag.Name, &flag)
		default:
			var subName string
			// 判断是否有定义子命令
			pfn, ok := flag.(*FlagSubCmd)
			if ok {
				subName = (*pfn).CmdName()
				var fn subFunc
				fn.runfn = (*pfn).CmdRun
				initfn, ok := flag.(*FlagSubInit)
				if ok {
					fn.initfn = (*initfn).CmdInit
				} else {
					fn.initfn = defaultInit
				}
				var alias []string
				aliasfn, ok := flag.(*FlagSubAlias)
				if ok {
					alias = (*aliasfn).CmdAlias()
				}
				mark, ok := flag.(*FlagSubMark)
				if ok {
					f.AddSubCommand(subName, alias, &fn, (*mark).CmdMark())
				} else {
					f.AddSubCommand(subName, alias, &fn)
				}

			} else {
				pfn, ok := flag.(FlagSubCmd)
				if ok {
					subName = pfn.CmdName()
					var fn subFunc
					fn.runfn = pfn.CmdRun
					initfn, ok := flag.(FlagSubInit)
					if ok {
						fn.initfn = initfn.CmdInit
					} else {
						fn.initfn = defaultInit
					}
					var alias []string
					aliasfn, ok := flag.(FlagSubAlias)
					if ok {
						alias = aliasfn.CmdAlias()
					}
					mark, ok := flag.(FlagSubMark)
					if ok {
						f.AddSubCommand(subName, alias, &fn, mark.CmdMark())
					} else {
						f.AddSubCommand(subName, alias, &fn)
					}
				}
			}

			// 解析结构体
			t := reflect.ValueOf(flag)
			if t.Kind() == reflect.Ptr {
				if t.IsNil() {
					t = reflect.New(t.Type())
				}
				t = t.Elem()
			}
			if t.Kind() == reflect.Slice {
				t = t.Elem()
			}
			switch t.Kind() {
			case reflect.Struct:
				for i := 0; i < t.NumField(); i++ {
					dtype := t.Type()
					value := dtype.Field(i)
					if !value.IsExported() {
						continue
					}
					v, ok := value.Tag.Lookup("flag")
					if ok {
						// 获取到定义的值
						tmp := Flag{
							Name: v,
						}
						d, ok := value.Tag.Lookup("default")
						if ok {
							tmp.Default = d
						}
						dc, ok := value.Tag.Lookup("dc")
						if ok {
							tmp.Description = dc
						}

						tmp.Value = reflect.ValueOf(flag).Elem().Field(i)
						f.appendToAll(v, &tmp)
						f.subCommandAddFlag(subName, &tmp)
					}
				}
			}
		}
	}
}

func (f *FlagSet) Parse(args ...string) {
	if len(args) == 0 && len(f.args) == 0 {
		args = os.Args
	}
	// 获取args的KV集合
	if len(f.args) == 0 {
		f.args = args
	}
	argsMap := f.ParseToKV()
	// 获取有效的值
	for key, flags := range f.all {
		value, ok := argsMap["-"+key]
		if ok {
			for _, flag := range flags {
				err := flag.Parse(value)
				if err == nil {
					f.appendToOK(key, flag)
				}
			}

		}
	}
	// 获取默认值
	for key, flags := range f.all {
		_, ok := f.ok[key]
		if ok {
			continue
		}
		for _, flag := range flags {
			err := flag.Parse(flag.Default)
			if err == nil {
				f.appendToOK(key, flag)
			}
		}
	}
	if f.ishelp {
		f.PrintUsage()
		os.Exit(0)
	}
}

func (f *FlagSet) Kvargs() map[string]string {
	return f.kvargs
}

func (f *FlagSet) ParseToKV(args ...string) map[string]string {
	if len(args) == 0 && len(f.args) == 0 {
		args = os.Args
	}
	if len(f.args) == 0 {
		f.args = args
	}
	// 清空旧的参数信息
	f.ishelp = false
	f.otherArgs = []string{}
	f.kvargs = make(map[string]string)
	for index, arg := range f.args {
		if strings.HasPrefix(arg, "-") {
			if arg == "-h" || arg == "--help" {
				f.ishelp = true
			}
			// 判断是否有等号
			argInfo := strings.Split(arg, "=")
			if len(argInfo) >= 2 {
				f.kvargs[argInfo[0]] = argInfo[1]
				// 否则取下一个值
			} else {
				next := index + 1
				if len(f.args) > next {
					info := f.args[next]
					if strings.HasPrefix(info, "-") {
						f.kvargs[argInfo[0]] = "true"
					} else {
						f.kvargs[argInfo[0]] = info
					}
				} else {
					f.kvargs[argInfo[0]] = "true"
				}
			}
		} else {
			// 获取之前那个
			prev := index - 1
			if prev >= 0 {
				info := f.args[prev]
				if strings.HasPrefix(info, "-") && len(strings.Split(info, "=")) == 1 {
					continue
				}
			}
			f.otherArgs = append(f.otherArgs, arg)
		}
	}
	return f.kvargs
}

func (f *FlagSet) Args() []string {
	return f.otherArgs
}

// 打印使用方法
func (f *FlagSet) PrintUsage() {
	if f.helpfn != nil {
		(*f.helpfn)()
		return
	}
	fmt.Printf("使用方法\n\n")
	if len(f.all) == 0 && len(f.subcmd) == 0 {
		fmt.Printf("  暂无绑定参数\n")
		return
	}

	// 打印子命令信息
	var subCmd []*FlagFunc
	for _, ff := range f.subcmd {
		subCmd = append(subCmd, ff)
	}
	sort.Slice(subCmd, func(i, j int) bool {
		return subCmd[i].name < subCmd[j].name
	})

	if len(subCmd) > 0 {
		fmt.Println("子命令列表:")
	}
	for _, ff := range subCmd {
		ff.PrintMark()
	}
	if len(subCmd) > 0 {
		fmt.Println("")
	}

	if !f.isPrintArgsList {
		return
	}
	// 排序列表
	var flags []*Flag
	for _, f2 := range f.all {
		flags = append(flags, f2...)
	}
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	if len(flags) > 0 {
		fmt.Println("参数列表:")
	}

	// 遍历列表打印使用方法
	for _, flag := range flags {
		// 打印使用方法
		flag.PrintDefault()
	}
}

type subFunc struct {
	initfn func(...string) error
	runfn  func() error
}

// 添加子命令
func (f *FlagSet) AddSubCommand(sub string, alias []string, fn *subFunc, dc ...string) {
	subCmd := FlagFunc{initfunc: fn.initfn, cmd: fn.runfn, name: sub}
	if len(dc) > 0 {
		subCmd.dc = strings.Join(dc, ",")
	} else {
		subCmd.dc = fmt.Sprintf("运行 %v 子命令", sub)
	}
	f.subcmd[sub] = &subCmd
	for _, sa := range alias {
		f.subcmd[sa] = &subCmd
	}
}

// 子命令添加定义参数
func (f *FlagSet) subCommandAddFlag(sub string, fg *Flag) {
	if sub == "" || f.subcmd[sub] == nil {
		return
	}
	f.subcmd[sub].flags = append(f.subcmd[sub].flags, fg)
}

// 解析并运行
func (f *FlagSet) ParseToRun(args ...string) error {
	f.Parse(args...)
	return f.Run()
}

// 直接运行
func (f *FlagSet) Run() error {
	if len(f.otherArgs) > 0 {
		fn := f.subcmd[f.otherArgs[0]]
		if fn != nil {
			return fn.Run(f.otherArgs...)
		}
	}
	if len(f.otherArgs) > 1 {
		fn := f.subcmd[f.otherArgs[1]]
		if fn != nil {
			return fn.Run(f.otherArgs[:1]...)
		}
	}
	if len(f.otherArgs) > 2 && f.otherArgs[1] == "help" {
		helpCmd := f.otherArgs[2]
		fn := f.subcmd[helpCmd]
		if fn != nil {
			fn.PrintHelp()
			return nil
		}
	}
	return fmt.Errorf("not sub command")
}
