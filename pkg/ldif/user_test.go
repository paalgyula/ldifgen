package ldif

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupGenerator(t *testing.T) *Generator {
	t.Helper()

	gen, err := NewGenerator(WithDomain("test.domain.com"))
	assert.NoError(t, err)

	return gen
}

func TestUserGenerator_Generate(t *testing.T) {
	ug := setupGenerator(t)

	t.Run("TestNoOrgUnits", func(t *testing.T) {
		_, err := ug.GenerateUser()
		assert.Error(t, err)
		assert.EqualError(t, ErrEmptyOrgUnits, err.Error())
	})

	ug.GenerateOrgUnits(2, 10)

	user, err := ug.GenerateUser()
	assert.NoError(t, err)

	assert.NotEmpty(t, user.Description)
	assert.NotEmpty(t, user.CommonName)
	assert.NotEmpty(t, user.GivenName)
	assert.NotEmpty(t, user.Surname)

	assert.Contains(t, user.DistinguishedName(), "cn=")
	assert.Contains(t, user.DistinguishedName(), "ou=")
	assert.Contains(t, user.DistinguishedName(), "dc=")
}

func TestUserGenerator_GenerateN(t *testing.T) {
	ug := setupGenerator(t)
	ug.GenerateOrgUnits(5, 2)

	gg, err := ug.GenerateUsers(10)
	assert.NoError(t, err)

	assert.Len(t, gg, 10)
}
