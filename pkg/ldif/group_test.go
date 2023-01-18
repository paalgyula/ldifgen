package ldif_test

import (
	"testing"

	"github.com/ebauman/ldifgen/pkg/ldif"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorGenerateGroups(t *testing.T) {
	gen, err := ldif.NewGenerator(
		ldif.WithDomain("test.domain.com"),
	)
	assert.NoError(t, err)

	gen.GenerateOrgUnits(5, 2)

	users, err := gen.GenerateUsers(10)
	assert.NoError(t, err)
	assert.Len(t, users, 10)

	gl, err := gen.GenerateGroups(300)
	assert.NoError(t, err)

	assert.Len(t, gl, 300)
}
