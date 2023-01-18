package ldif

import (
	"errors"
	"strings"
)

var (
	ErrEmptyOrgUnits = errors.New("no org units defined/generated")
	ErrEmptyUsers    = errors.New("no users defined/generated")
)

type Generator struct {
	domain []string

	OrgUnits []string
	Users    []*User

	nameGenerator *NameGenerator
}

type GeneratorOption func(*Generator)

// WithDomain
func WithDomain(domain string) GeneratorOption {
	return func(g *Generator) {
		g.domain = strings.Split(domain, ".")
	}
}

func NewGenerator(opts ...GeneratorOption) (*Generator, error) {
	g := &Generator{
		domain: []string{"test", "domain", "com"},
	}

	for _, o := range opts {
		o(g)
	}

	ng, err := NewNameGenerator()
	if err != nil {
		return nil, err
	}

	g.nameGenerator = ng

	return g, nil
}
