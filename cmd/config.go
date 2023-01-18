package cmd

import (
	"strings"

	"github.com/ebauman/ldifgen/pkg/ldif"
)

type GenerateConfig struct {
	Users                    int
	Groups                   int
	OUs                      int
	OUDepth                  int
	Domain                   []string
	UserClasses              []string
	GroupClasses             []string
	OUClasses                []string
	UserChangeType           string
	GroupChangeType          string
	OUChangeType             string
	BuzzwordDataset          string
	DepartmentDataset        string
	FirstNameDataset         string
	LastNameDataset          string
	GroupsDataset            string
	GroupMembershipAttribute string
}

type RenderConfig struct {
	GenerateConfig

	Users  []*ldif.User
	Domain []string
	OUs    []string
	Groups []*ldif.Group
	Time   string
}

func (c RenderConfig) DC() string {
	return "dc=" + strings.Join(c.Domain, ",dc=")
}

func (c RenderConfig) TrimOU(ou string) string {
	return strings.Split(ou, ",")[0]
}
