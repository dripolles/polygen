package generator

import (
	"errors"
	"fmt"
	"strings"
)

type polymap struct {
	m map[string]string
}

func (p *polymap) Type(alias string) (string, error) {
	t, ok := p.m[alias]
	if !ok {
		msg := fmt.Sprintf("Unknown alias '%s'", alias)
		return "", errors.New(msg)
	}

	return t, nil
}

func (p *polymap) Name(alias string) (string, error) {
	s, err := p.Type(alias)
	if err != nil {
		return "", err
	}
	s = strings.Replace(s, "[]", "List", -1)
	s = strings.Replace(s, "*", "PtrTo", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)

	return s, nil
}
