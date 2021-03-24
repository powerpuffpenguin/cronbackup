package backend

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dop251/goja"
)

func osRequire(runtime *goja.Runtime, module *goja.Object) {
	native := osNative{
		runtime: runtime,
	}
	native.Require(module)
}

type osNative struct {
	runtime *goja.Runtime
}

func (native osNative) Require(module *goja.Object) {
	obj := module.Get(`exports`).(*goja.Object)
	obj.Set(`join`, native.join)
	obj.Set(`exec`, native.exec)
	obj.Set(`cwdExec`, native.cwdExec)
	obj.Set(`readFile`, native.readFile)
}
func (native osNative) join(elem ...string) string {
	return filepath.Join(elem...)
}
func (native osNative) exec(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	e := cmd.Run()
	if e != nil {
		panic(native.runtime.ToValue(e))
	}
}
func (native osNative) cwdExec(cwd, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	e := cmd.Run()
	if e != nil {
		panic(native.runtime.ToValue(e))
	}
}
func (native osNative) readFile(filename string) string {
	b, e := ioutil.ReadFile(filename)
	if e == nil {
		return string(b)
	}
	panic(native.runtime.ToValue(e))
}
