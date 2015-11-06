package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

type Generator struct {
	Types       TypeAssignments
	Source      string
	Destination string
}

type TypeAssignments map[string]string

func NewGenerator(
	types TypeAssignments, source, destination string,
) *Generator {
	if destination == "" {
		destination = getDestinationFromTypes(source, types)
	}
	return &Generator{
		Types:       types,
		Source:      source,
		Destination: destination,
	}
}

func (g *Generator) Generate() error {
	t, err := g.getTemplate()
	if err != nil {
		return err
	}

	f, err := g.getDestFile()
	if err != nil {
		return err
	}
	defer f.Close()

	code, codeErr := g.executeTemplate(t)

	if _, err := f.Write(code); err != nil {
		return err
	}

	return codeErr
}

func (g *Generator) getTemplate() (*template.Template, error) {
	t := template.New(g.Destination)
	g.addTemplateFunctions(t)
	tcode, err := g.getTemplateCode()
	if err != nil {
		return nil, err
	}
	t, err = t.Parse(tcode)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (g *Generator) addTemplateFunctions(t *template.Template) {
	p := g.createPolymap()
	fm := template.FuncMap{
		"T":    p.Type,
		"Name": p.Name,
	}

	t.Funcs(fm)
}

func (g *Generator) createPolymap() *polymap {
	p := &polymap{
		m: map[string]string{},
	}
	for alias, t := range g.Types {
		p.m[alias] = t
	}

	return p
}

func (g *Generator) getTemplateCode() (string, error) {
	gopath := os.Getenv("GOPATH")
	filename := gopath + "/src/" + g.Source
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	str := string(code)
	pkg, err := g.getPackage()
	if err != nil {
		return "", err
	}
	str = fmt.Sprintf("package %s\n%s", pkg, str)

	return str, nil
}

func (g *Generator) getPackage() (string, error) {
	pkg := os.Getenv("GOPACKAGE")
	if pkg == "" {
		err := errors.New("No $GOPACKAGE environment variable set")
		return "", err
	}

	return pkg, nil
}

func (g *Generator) executeTemplate(t *template.Template) ([]byte, error) {
	var buf bytes.Buffer

	err := t.Execute(&buf, nil)
	if err != nil {
		return []byte{}, err
	}

	raw := buf.Bytes()
	code, err := g.prettyfy(raw)
	if err != nil {
		msg := fmt.Sprintf("invalid code in template '%s'\n%s", t.Name(), err.Error())
		err = errors.New(msg)
		code = raw
	}

	return code, err
}

func (g *Generator) prettyfy(code []byte) ([]byte, error) {
	pretty, err := format.Source(code)
	if err != nil {
		pretty = code
	}

	return pretty, err
}

func (g *Generator) getDestFile() (*os.File, error) {
	if err := g.ensureDestFileDoesNotExist(); err != nil {
		return nil, err
	}

	file, err := os.Create(g.Destination)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (g *Generator) ensureDestFileDoesNotExist() error {
	if _, err := os.Stat(g.Destination); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(g.Destination)
}
