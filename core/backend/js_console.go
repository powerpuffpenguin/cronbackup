package backend

import (
	"fmt"

	"github.com/dop251/goja"
)

const (
	// ModuleID .
	ModuleID = `console`
)

func nativeLog(call goja.FunctionCall) goja.Value {
	count := len(call.Arguments)
	if count == 0 {
		fmt.Println()
		return nil
	}
	for i := 0; i < count; i++ {
		if i != 0 {
			fmt.Print(` `)
		}
		arg := call.Arguments[i]
		if obj, ok := arg.(*goja.Object); ok {
			if obj.ClassName() == "Error" {
				fmt.Print(arg.String())
			} else {
				b, e := obj.MarshalJSON()
				if e != nil {
					fmt.Print(obj.ClassName())
				} else {
					fmt.Printf(`%s`, b)
				}
			}
		} else {
			fmt.Print(arg.String())
		}
	}
	fmt.Println()
	return nil
}

func consoleRequire(vm *goja.Runtime, module *goja.Object) {
	obj := module.Get(`exports`).(*goja.Object)
	obj.Set(`trace`, nativeLog)
	obj.Set(`debug`, nativeLog)
	obj.Set(`log`, nativeLog)
	obj.Set(`info`, nativeLog)
	obj.Set(`warn`, nativeLog)
	obj.Set(`error`, nativeLog)
}
