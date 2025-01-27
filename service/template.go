package service

type Template struct {
	Name string
	Port int
}

func NewTemplate(name string, port int) Template {
	return Template{
		Name: name,
		Port: port,
	}
}

func (t *Template) LongName() string {
	return t.Name + "-service"
}

func (t *Template) DefaultPort() int {
	return t.Port
}
