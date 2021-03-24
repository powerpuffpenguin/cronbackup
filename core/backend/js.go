package backend

import (
	"errors"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// JSBackend js backend
type JSBackend struct {
	runtime       *goja.Runtime
	backup        goja.Callable
	pack          goja.Callable
	removeExpired goja.Callable
	finish        goja.Callable
}

// NewJSBackend new js backend
func NewJSBackend(filename string) (js *JSBackend, e error) {
	runtime := goja.New()
	registry := require.NewRegistry(
		require.WithGlobalFolders(filepath.Join(filepath.Dir(filename), `node_modules`)),
	)
	registry.Enable(runtime)
	// console
	require.RegisterNativeModule(`console`, consoleRequire)
	runtime.Set(`console`, require.Require(runtime, `console`)) // set to global
	// os
	require.RegisterNativeModule(`os`, osRequire)

	// check
	var (
		backup        goja.Callable
		pack          goja.Callable
		removeExpired goja.Callable
		finish        goja.Callable
	)
	v := require.Require(runtime, filename)
	if obj, ok := v.(*goja.Object); ok {
		backup, ok = goja.AssertFunction(obj.Get(`backup`))
		if !ok {
			e = errors.New(`js backend not implement function backup : ` + filename)
			return
		}
		pack, ok = goja.AssertFunction(obj.Get(`pack`))
		if !ok {
			e = errors.New(`js backend not implement function pack : ` + filename)
			return
		}
		removeExpired, ok = goja.AssertFunction(obj.Get(`removeExpired`))
		if !ok {
			e = errors.New(`js backend not implement function removeExpired : ` + filename)
			return
		}
		finish, ok = goja.AssertFunction(obj.Get(`finish`))
		if !ok {
			e = errors.New(`js backend not implement function finish : ` + filename)
			return
		}
	} else {
		e = errors.New(`js backend not implement : ` + filename)
		return
	}

	js = &JSBackend{
		runtime:       runtime,
		backup:        backup,
		pack:          pack,
		removeExpired: removeExpired,
		finish:        finish,
	}
	return
}
func (js *JSBackend) Backup(md Metadata) (e error) {
	_, e = js.backup(goja.Undefined(), js.runtime.ToValue(md))
	return
}
func (js *JSBackend) Pack(md Metadata) (e error) {
	_, e = js.pack(goja.Undefined(), js.runtime.ToValue(md))
	return
}
func (js *JSBackend) RemoveExpired(md Metadata) (e error) {
	_, e = js.removeExpired(goja.Undefined(), js.runtime.ToValue(md))
	return
}
func (js *JSBackend) Finish(md Metadata) (e error) {
	_, e = js.finish(goja.Undefined(), js.runtime.ToValue(md))
	return
}
