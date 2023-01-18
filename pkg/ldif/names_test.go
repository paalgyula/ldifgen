package ldif_test

import (
	"testing"

	generator "github.com/ebauman/ldifgen/pkg/ldif"
	"github.com/stretchr/testify/assert"
)

func TestNameGenerator(t *testing.T) {
	gen, err := generator.NewNameGenerator()
	assert.NoError(t, err)

	for i := 0; i < 50; i++ {
		dep := gen.Department()
		assert.NotEmpty(t, dep)
	}

	for i := 0; i < 50; i++ {
		n := gen.FirstName()
		assert.NotEmpty(t, n)
	}

	for i := 0; i < 50; i++ {
		n := gen.LastName()
		assert.NotEmpty(t, n)
	}

	for i := 0; i < 50; i++ {
		g := gen.Group()
		assert.NotEmpty(t, g)
	}
}
