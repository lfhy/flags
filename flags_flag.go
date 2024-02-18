package flags

import (
	"fmt"
	"reflect"
	"strconv"
)

type Flag struct {
	Name        string
	Type        string
	Default     string
	Description string
	Value       any
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

// 打印默认方法
func (flag *Flag) PrintDefault() {
	fmt.Printf("  -%v\n", flag.Name)

	if flag.Description == "" {
		if flag.Default != "" {
			fmt.Printf("        默认值:%v\n", flag.Default)
		} else {
			fmt.Printf("        无默认值\n")
		}
	} else {
		fmt.Printf("        %v ", flag.Description)
		if flag.Default != "" {
			fmt.Printf("(默认值:%v)\n", flag.Default)
		} else {
			fmt.Printf("\n")
		}
	}
}
