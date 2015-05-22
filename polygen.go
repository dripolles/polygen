package main

import (
	"bytes"
	"fmt"
	"io"
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

	gen.Execute()
}

type Generator struct {
	Types map[string]string `short:"t" long:"types"`
	Args  struct {
		Source string `positional-arg-name:"source"`
		Dest   string `positional-arg-name:"destination"`
	} `positional-args:"true"`
}

func (g *Generator) Dest() string {
	return g.Args.Dest + ".go"
}

func newGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Execute() error {
	t := g.getTemplate()
	f := g.getDestFile()
	defer f.Close()
	g.executeTemplate(f, t)

	return nil
}

func (g *Generator) getTemplate() *template.Template {
	t := template.New(g.Args.Dest)
	g.addTemplateFunctions(t)
	tcode := g.getTemplateCode()
	t = template.Must(t.Parse(tcode))

	return t
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

func (g *Generator) getTemplateCode() string {
	gopath := os.Getenv("GOPATH")
	filename := gopath + "/src/" + g.Args.Source
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	str := string(code)
	str = fmt.Sprintf("package %s\n%s", g.getPackage(), str)

	return str
}

func (g *Generator) getPackage() string {
	return os.Getenv("GOPACKAGE")
}

func (g *Generator) executeTemplate(wr io.Writer, t *template.Template) error {
	var buf bytes.Buffer

	err := t.Execute(&buf, nil)
	if err != nil {
		panic(err)
		return err
	}

	return prettyfy(buf.Bytes(), wr)
}

func prettyfy(input []byte, wr io.Writer) error {
	output, err := format.Source(input)
	if err != nil {
		panic(err)
		return err
	}

	_, err = wr.Write(output)
	return err
}

func (g *Generator) getDestFile() *os.File {
	filename := g.Dest()
	os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return file
}

type polymap struct {
	m map[string]string
}

func (p *polymap) T(alias string) string {
	t, ok := p.m[alias]
	if !ok {
		msg := fmt.Sprintf("Unknown alias '%s'", alias)
		panic(msg)
	}

	return t
}

func (p *polymap) Id(alias string) string {
	s := p.T(alias)
	s = strings.Replace(s, "[]", "List", -1)
	s = strings.Replace(s, "*", "To", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)

	return s
}
