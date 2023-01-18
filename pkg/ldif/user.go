package ldif

import (
	"fmt"
	"strings"
)

type User struct {
	GivenName          string
	CommonName         string
	Surname            string
	Manager            string
	Secretary          string
	Title              string
	Description        string
	OrganizationalUnit string

	domain string
}

func (u User) DistinguishedName() string {
	return fmt.Sprintf("cn=%s,ou=%s,%s",
		u.CommonName,
		u.OrganizationalUnit,
		u.domain,
	)
}

func (u User) UID() string {
	return fmt.Sprintf("%s.%s", u.GivenName, u.Surname)
}

func (g *Generator) GenerateUser() (*User, error) {
	if len(g.OrgUnits) == 0 {
		return nil, ErrEmptyOrgUnits
	}

	u := &User{}
	u.domain = "dc=" + strings.Join(g.domain, ",dc=")
	u.GivenName = g.nameGenerator.FirstName()
	u.Surname = g.nameGenerator.LastName()
	u.CommonName = u.GivenName + " " + u.Surname
	u.Description = "This is the description for " + u.CommonName + " " + u.Surname
	u.OrganizationalUnit = g.randomOU()

	return u, nil
}

func (g *Generator) GenerateUsers(count int) ([]*User, error) {
	users := make([]*User, count)
	var err error

	for i := 0; i < count; i++ {
		users[i], err = g.GenerateUser()
		if err != nil {
			return nil, err
		}
	}

	g.Users = users

	return users, nil
}
