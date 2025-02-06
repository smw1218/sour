package generator

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/service
var serviceTemplates embed.FS

var mainTmpl *template.Template
var appTmpl *template.Template

func init() {
	var err error
	mainTmpl, err = template.ParseFS(serviceTemplates, "templates/service/main.go.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	appTmpl, err = template.ParseFS(serviceTemplates, "templates/service/service.go.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
}

func NewServiceData(name string, port int, packageName string) (*ServiceData, error) {
	err := ValidateServiceName(name)
	if err != nil {
		return nil, err
	}
	return &ServiceData{
		ServiceName: name,
		Port:        port,
		Package:     packageName,
	}, nil
}

type ServiceData struct {
	ServiceName string
	Port        int
	Package     string
}

func (sd *ServiceData) TitleName() string {
	caser := cases.Title(language.English)
	words := strings.Split(sd.ServiceName, "-")
	upper := ""
	for _, w := range words {
		upper += caser.String(w)
	}
	return upper
}

func (sd *ServiceData) ServiceType() string {
	return sd.TitleName() + "Service"
}

func (sd *ServiceData) ServiceDirectory(others ...string) string {
	return path.Join(slices.Concat([]string{".", "cmd", sd.ServiceName + "-service"}, others)...)
}

func (sd *ServiceData) CreateService() error {
	mainPath := path.Join(sd.ServiceDirectory("main.go"))
	err := sd.doTemplate(mainTmpl, mainPath)
	if err != nil {
		return err
	}

	appPath := path.Join(sd.ServiceDirectory("app", "service.go"))
	err = sd.doTemplate(appTmpl, appPath)
	if err != nil {
		return err
	}
	return nil
}

func (sd *ServiceData) doTemplate(tmpl *template.Template, outpath string) error {
	_, err := os.Stat(outpath)
	if err == nil {
		return fmt.Errorf("%v already exists; not overwriting", outpath)
	}
	err = os.MkdirAll(path.Dir(outpath), 0755)
	if err != nil {
		return fmt.Errorf("failed making directories: %v: %w", path.Dir(outpath), err)
	}
	f, err := os.Create(outpath)
	if err != nil {
		return fmt.Errorf("failed to create %v: %w", outpath, err)
	}
	defer f.Close()
	err = tmpl.Execute(f, sd)
	if err != nil {
		return fmt.Errorf("failed executing template for %v: %w", outpath, err)
	}
	return nil
}

// numbers are for easier testing; but don't encourage them
var validServiceNameRe = regexp.MustCompile(`^[-a-z0-9]+$`)

func ValidateServiceName(serviceName string) error {
	if validServiceNameRe.MatchString(serviceName) {
		return nil
	}
	return fmt.Errorf("service names must only contain lower case letters and \"-\"")
}
