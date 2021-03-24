package backend

import "errors"

var defaultBackup emptyBackend
var ErrNotImplemented = errors.New(`backend not implement`)

func Empty() Backend {
	return defaultBackup
}

type Metadata struct {
	ID uint64

	Host     string
	Port     uint16
	Username string
	Password string

	Output string
}

type Backend interface {
	Backup(md Metadata) (e error)
	Package(md Metadata) (e error)
	RemoveExpired(md Metadata) (e error)
	Finish(md Metadata) (e error)
}
type emptyBackend uint

func (emptyBackend) Backup(md Metadata) (e error) {
	return ErrNotImplemented
}
func (emptyBackend) Package(md Metadata) (e error) {
	return ErrNotImplemented
}
func (emptyBackend) RemoveExpired(md Metadata) (e error) {
	return ErrNotImplemented
}
func (emptyBackend) Finish(md Metadata) (e error) {
	return ErrNotImplemented
}
