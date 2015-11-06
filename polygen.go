package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/dripolles/polygen/generator"
)

func main() {
	gen := newGeneratorCommand()
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

type GeneratorCommand struct {
	Types generator.TypeAssignments `short:"t" long:"types" required:"true"`
	Args  struct {
		Source string `positional-arg-name:"source" description:"source file"`
		Dest   string `positional-arg-name:"destination" description:"destination file"`
	} `positional-args:"true"`
}

func newGeneratorCommand() *GeneratorCommand {
	return &GeneratorCommand{}
}

func (g *GeneratorCommand) Dest() string {
	d := g.Args.Dest
	if d == "" {
		return d
	}

	if strings.HasSuffix(d, ".go") {
		return d
	}

	return d + ".go"
}

func (g *GeneratorCommand) Execute() error {
	if g.Args.Source == "" {
		return errors.New("Source file must be provided")
	}
	gen := generator.NewGenerator(g.Types, g.Args.Source, g.Dest())
	err := gen.Generate()

	return err
}
