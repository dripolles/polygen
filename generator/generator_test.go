package generator

import (
	"fmt"
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
	types := TypeAssignments{
		"a": "int",
	}

	dest, err := s.generate("convertslice.tgo", types)
	c.Assert(err, IsNil)

	c.Assert(os.Remove(dest), IsNil)
}

func (s *GeneratorSuite) generate(
	name string, types TypeAssignments,
) (string, error) {
	source := s.getSource("convertslice.tgo")
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
