package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GeneratorSuite struct {
	c *C
}

var _ = Suite(&GeneratorSuite{})

func (s *GeneratorSuite) SetUpSuite(c *C) {
	err := os.Setenv("GOPACKAGE", "testpackage")
	c.Assert(err, IsNil)
}

func (s *GeneratorSuite) SetUpTest(c *C) {
	s.c = c
}

func (s *GeneratorSuite) TestGenerateBasic(c *C) {
	types := TypeAssignments{"a": "int"}
	dest, err := s.generate("convertslice.tgo", types)
	c.Assert(err, IsNil)

	info := s.checkTypes(dest)
	s.checkTypeDef(info, "convertintslice", "func(xs []int) interface{}")

	c.Assert(os.Remove(dest), IsNil)
}

func (s *GeneratorSuite) TestSyntaxError(c *C) {
	dest, err := s.generate("syntaxerror.tgo", TypeAssignments{"a": "int"})
	c.Assert(err, Not(IsNil))

	generated, err := ioutil.ReadFile(dest)
	c.Assert(err, IsNil)
	c.Assert(os.Remove(dest), IsNil)

	c.Assert(string(generated), DeepEquals, syntaxerrorFixture)
}

func (s *GeneratorSuite) generate(
	name string, types TypeAssignments,
) (string, error) {
	source := s.getSource(name)
	dest, err := s.getDestination()
	s.c.Assert(err, IsNil)

	gen := NewGenerator(types, source, dest)
	err = gen.Generate()

	return dest, err
}

func (s *GeneratorSuite) getSource(name string) string {
	return fmt.Sprintf(
		"github.com/dripolles/polygen/generator/fixtures/%s",
		name,
	)
}

func (s *GeneratorSuite) getDestination() (string, error) {
	destFile, err := ioutil.TempFile("", "polygen")
	if err != nil {
		return "", err
	}
	defer destFile.Close()
	dest := destFile.Name()

	return dest, nil
}

func (s *GeneratorSuite) checkTypes(filename string) *types.Info {
	fs := token.NewFileSet()
	checker, info := s.getChecker(fs)
	f := s.getAstFile(filename, fs)
	checker.Files([]*ast.File{f})

	return info
}

func (s *GeneratorSuite) getChecker(
	fs *token.FileSet,
) (*types.Checker, *types.Info) {
	config := &types.Config{
		Importer: importer.Default(),
	}
	pkg := types.NewPackage("", "testpackage")
	info := &types.Info{
		Defs: map[*ast.Ident]types.Object{},
	}
	checker := types.NewChecker(config, fs, pkg, info)

	return checker, info
}

func (s *GeneratorSuite) getAstFile(
	filename string, fs *token.FileSet,
) *ast.File {
	f, err := parser.ParseFile(fs, filename, nil, 0)
	s.c.Assert(err, IsNil)

	return f
}

func (s *GeneratorSuite) checkTypeDef(info *types.Info, name, typestr string) {
	for ident, obj := range info.Defs {
		if ident.Name == name {
			s.c.Assert(obj.Type(), Equals, typestr)
			return
		}
	}

	s.c.Error("Name %s not found", name)
}

var syntaxerrorFixture = `package testpackage
func syntaxerror(x int) {
	this is not valid code
}
`
