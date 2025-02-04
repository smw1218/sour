package service

import (
	"github.com/gin-gonic/gin"
)

// ServiceInterface will be used in two specific ways to control the
// lifecycle of a microservice. In normal operation, main will call
// Setup, RegisterRoutes then service.Run.
// The local proxy that can route to any service will only call
// RegisterRoutes and DefaultPort to find out where each endpoint should
// be routed. This means that RegisterRoutes must not dereference any nil
// pointers that my be needed from calling Setup. Typically, this not a
// problem because the gin.HandlerFuncs will exist on a nil pointer and
// since the local proxy will never actually call the method, then nothing
// blows up.
type ServiceInterface interface {
	Setup() error
	RegisterRoutes(gin.IRouter)
	Shutdown() error
	Name() string
	LongName() string
	DefaultPort() int
}

type Template struct {
	name string
	Port int
}

func NewTemplate(name string, port int) Template {
	return Template{
		name: name,
		Port: port,
	}
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) LongName() string {
	return t.name + "-service"
}

func (t *Template) DefaultPort() int {
	return t.Port
}

func (t *Template) Shutdown() error {
	return nil
}
