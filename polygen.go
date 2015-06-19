package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"go/format"

	"github.com/jessevdk/go-flags"
)

func main() {
	gen := newGenerator()
	parser := flags.NewParser(gen, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

	err = gen.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Generator struct {
	Types map[string]string `short:"t" long:"types" required:"true"`
	Args  struct {
		Source string `positional-arg-name:"source" description:"source file"`
		Dest   string `positional-arg-name:"destination" description:"destination file"`
	} `positional-args:"true" required:"true"`
}

func (g *Generator) Dest() string {
	d := g.Args.Dest
	if strings.HasSuffix(d, ".go") {
		return d
	}

	return d + ".go"
}

func newGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Execute() error {
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
	t := template.New(g.Args.Dest)
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
		"T":  p.T,
		"Id": p.Id,
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
	filename := gopath + "/src/" + g.Args.Source
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

	code, err := g.prettyfy(buf.Bytes())
	if err != nil {
		msg := fmt.Sprintf("invalid code in template '%s'\n%s", t.Name(), err.Error())
		err = errors.New(msg)
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
	filename := g.Dest()
	if err := os.Remove(filename); err != nil {
		return nil, err
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type polymap struct {
	m map[string]string
}

func (p *polymap) T(alias string) (string, error) {
	t, ok := p.m[alias]
	if !ok {
		msg := fmt.Sprintf("Unknown alias '%s'", alias)
		return "", errors.New(msg)
	}

	return t, nil
}

func (p *polymap) Id(alias string) (string, error) {
	s, err := p.T(alias)
	if err != nil {
		return "", err
	}
	s = strings.Replace(s, "[]", "List", -1)
	s = strings.Replace(s, "*", "To", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)

	return s, nil
}
