package generator

import (
	"embed"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/domain
var domainTemplates embed.FS

var domainHandlersTmpl *template.Template

func init() {
	var err error
	domainHandlersTmpl, err = template.ParseFS(domainTemplates, "templates/domain/handlers.go.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
}

func NewDomainData(sd *ServiceData, domainName string) (*DomainData, error) {
	err := ValidateDomainName(domainName)
	if err != nil {
		return nil, err
	}
	return &DomainData{
		ServiceData: *sd,
		DomainName:  domainName,
	}, nil
}

type DomainData struct {
	ServiceData
	DomainName string
}

func (dd *DomainData) DomainPackage() string {
	dd.TitleName()
	return strings.ToLower(dd.DomainName)
}

var lowerCaseRe = regexp.MustCompile(`[a-z]`)

func (dd *DomainData) DomainInitials() string {
	return strings.ToLower(lowerCaseRe.ReplaceAllString(dd.DomainName, ""))
}

func (dd *DomainData) CreateDomain() error {
	handlersPath := path.Join(dd.ServiceDirectory(dd.DomainPackage(), "handlers.go"))
	err := doTemplate(dd, domainHandlersTmpl, handlersPath)
	if err != nil {
		return err
	}
	return nil
}

var validDomainNameRe = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)

func ValidateDomainName(serviceName string) error {
	if validDomainNameRe.MatchString(serviceName) {
		return nil
	}
	return fmt.Errorf("domain names must be pascal-cased starting with a capital letter (valid exported type in go)")
}
