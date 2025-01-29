package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Template struct {
	Name   string
	Port   int
	Engine *gin.Engine
}

func NewTemplate(name string, port int, engine *gin.Engine) Template {
	return Template{
		Name:   name,
		Port:   port,
		Engine: engine,
	}
}

func (t *Template) LongName() string {
	return t.Name + "-service"
}

func (t *Template) DefaultPort() int {
	return t.Port
}

func (t *Template) Run() error {
	return t.Engine.Run(fmt.Sprintf(":%v", t.Port))
}
