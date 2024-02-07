package flags

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type FlagSet struct {
	args      []string
	kvargs    map[string]string
	all       map[string]*Flag
	ok        map[string]*Flag
	otherArgs []string
	ishelp    bool
}

func NewFlags() *FlagSet {
	var f FlagSet
	f.kvargs = make(map[string]string)
	f.all = make(map[string]*Flag)
	f.ok = make(map[string]*Flag)
	return &f
}

func (f *FlagSet) Var(flags ...any) {
	for _, f2 := range flags {
		switch flag := f2.(type) {
		case *Flag:
			f.all[flag.Name] = flag
		case Flag:
			f.all[flag.Name] = &flag
		default:
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
						f.all[v] = &tmp
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
	for key, flag := range f.all {
		value, ok := argsMap["-"+key]
		if ok {
			err := flag.Parse(value)
			if err == nil {
				f.ok[key] = flag
			}
		}
	}
	// 获取默认值
	for key, flag := range f.all {
		_, ok := f.ok[key]
		if ok {
			continue
		}
		err := flag.Parse(flag.Default)
		if err == nil {
			f.ok[key] = flag
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
	fmt.Println("使用方法:")
	if len(f.all) == 0 {
		fmt.Printf("  暂无绑定参数\n")
		return
	}
	// 排序列表
	var flags []*Flag
	for _, f2 := range f.all {
		flags = append(flags, f2)
	}
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	// 遍历列表打印使用方法
	for _, flag := range flags {
		// 打印使用方法
		flag.PrintDefault()
	}
}
