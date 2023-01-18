package ldif

import (
	"fmt"
	"strings"
)

type Group struct {
	CommonName         string
	OrganizationalUnit string
	Members            []string

	domain string
}

func (g Group) DistinguishedName() string {
	return "cn=" + g.CommonName + ",ou=" + g.OrganizationalUnit + g.domain
}

func (g *Generator) GenerateGroup(members []*User) (*Group, error) {
	if len(g.OrgUnits) == 0 {
		return nil, ErrEmptyOrgUnits
	}

	newGroup := &Group{}
	newGroup.domain = "dc=" + strings.Join(g.domain, ",dc=")
	newGroup.CommonName = g.nameGenerator.Group()
	newGroup.OrganizationalUnit = g.randomOU()
	newGroup.Members = make([]string, 0)

	for _, m := range members {
		newGroup.Members = append(newGroup.Members, m.DistinguishedName())
	}

	return newGroup, nil
}

func (g *Generator) GenerateGroups(count int) ([]*Group, error) {
	if len(g.Users) == 0 {
		return nil, ErrEmptyUsers
	}

	userChunkSize := len(g.Users) / count

	groupList := make([]*Group, count)
	groupMap := make(map[string]int)

	userChunkPos := 0
	for i := 0; i < count; i++ {
		members := g.Users[userChunkPos:((i + 1) * userChunkSize)]

		g, err := g.GenerateGroup(members)
		if err != nil {
			return nil, err
		}

		if v := groupMap[g.DistinguishedName()]; v > 0 {
			// Append a number to it
			g.CommonName = fmt.Sprintf("%s #%d", g.CommonName, v)
		}

		groupMap[g.DistinguishedName()]++
		groupList[i] = g
		userChunkPos = (i + 1) * userChunkSize
	}

	return groupList, nil
}
