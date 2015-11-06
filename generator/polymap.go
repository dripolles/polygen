package generator

import (
	"errors"
	"fmt"
	"strings"
)

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
