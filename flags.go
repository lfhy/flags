package flags

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type FlagSet struct {
	args      []string
	kvargs    map[string]string
	all       map[string]*Flag
	ok        map[string]*Flag
	otherArgs []string
}

func NewFlags() *FlagSet {
	var f FlagSet
	f.kvargs = make(map[string]string)
	f.all = make(map[string]*Flag)
	f.ok = make(map[string]*Flag)
	return &f
}

type Flag struct {
	Name    string
	Type    string
	Default string
	Value   any
}

func (flag Flag) Var(f ...*FlagSet) {
	if len(f) == 0 {
		f = append(f, argsFlag)
	}
	for _, fs := range f {
		fs.Var(&flag)
	}
}

func (flag *Flag) errorInput(err any) error {
	return fmt.Errorf("输入错误: %v", err)
}

func (flag *Flag) errorUnknowType(err any) error {
	return fmt.Errorf("不支持的类型: %v", err)
}

func (flag *Flag) Parse(value string) error {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if rv, ok := flag.Value.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(flag.Value)
	}
	// 取出真实类型

	for {
		reflectKind = reflectValue.Kind()
		switch reflectKind {
		case reflect.Ptr:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// 为空就创一个默认值出来
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			} else {
				reflectValue = reflectValue.Elem()
			}
		case reflect.Int:
			data, err := strconv.Atoi(value)
			if err != nil {
				return flag.errorInput(err)
			} else {
				reflectValue.SetInt(int64(data))
				return nil
			}
		case reflect.String:
			reflectValue.SetString(value)
			return nil
		case reflect.Bool:
			data, err := strconv.ParseBool(value)
			if err != nil {
				return flag.errorInput(err)
			} else {
				reflectValue.SetBool(data)
				return nil
			}
		case reflect.Float64:
			data, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return flag.errorInput(err)
			} else {
				reflectValue.SetFloat(data)
				return nil
			}
		case reflect.Uint:
			data, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return flag.errorInput(err)
			} else {
				reflectValue.SetUint(data)
				return nil
			}
		default:
			// 不支持的类型
			return flag.errorUnknowType(reflectKind)
		}
	}
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
	f.otherArgs = []string{}
	f.kvargs = make(map[string]string)
	for index, arg := range f.args {
		if strings.HasPrefix(arg, "-") {
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
