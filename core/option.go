package core

import (
	"github.com/powerpuffpenguin/cronbackup/core/backend"
)

var defaultCoreOptions = coreOptions{
	serverHost:     `localhost`,
	serverPort:     3306,
	serverUsername: `root`,

	output:  `output`,
	backend: backend.Empty(),
}

type coreOptions struct {
	serverHost             string
	serverPort             uint16
	serverUsername         string
	serverPassword         string
	output                 string
	backend                backend.Backend
	immediate, description bool
}
type CoreOption interface {
	apply(*coreOptions)
}
type funcCoreOption struct {
	f func(*coreOptions)
}

func (fdo *funcCoreOption) apply(do *coreOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*coreOptions)) *funcCoreOption {
	return &funcCoreOption{
		f: f,
	}
}
func WithServer(host string, port uint16) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.serverHost = host
		o.serverPort = port
	})
}
func WithAuth(username, password string) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.serverUsername = username
		o.serverPassword = password
	})
}
func WithOutput(output string) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.output = output
	})
}
func WithBackend(backend backend.Backend) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.backend = backend
	})
}

func WithImmediate(immediate bool) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.immediate = immediate
	})
}
func WithDescription(description bool) CoreOption {
	return newFuncServerOption(func(o *coreOptions) {
		o.description = description
	})
}
