package cmd

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ebauman/ldifgen/pkg/ldif"
	"github.com/urfave/cli/v2"
)

//go:embed ldif.tmpl
var templateData []byte

const datasetString = "path to an alternative list of %s, used in %s generation. provide list of words, separated by newlines"

func GenerateCommand() *cli.Command {
	generateFlags := []cli.Flag{
		&cli.IntFlag{
			Name:  "users",
			Value: 10,
			Usage: "number of users to generate",
		},
		&cli.IntFlag{
			Name:  "ous",
			Value: 2,
			Usage: "number of organizational units to generate",
		},
		&cli.IntFlag{
			Name:  "ou-depth",
			Value: 1,
			Usage: "depth of generated OUs. specify n>1 to create 'chains' of OUs",
		},
		&cli.IntFlag{
			Name:  "groups",
			Value: 2,
			Usage: "number of groups to generate",
		},
		&cli.StringFlag{
			Name:  "domain",
			Value: "domain.example.org",
			Usage: "domain used to generate DC components, e.g. dc=domain,dc=example,dc=org",
		},
		&cli.StringFlag{
			Name:  "user-classes",
			Value: "top,person,organizationalPerson,inetOrgPerson",
			Usage: "comma-separated list of classes for user objects",
		},
		&cli.StringFlag{
			Name:  "ou-classes",
			Value: "top,organizationalUnit",
			Usage: "comma-separated list of classes for organizational unit objects",
		},
		&cli.StringFlag{
			Name:  "group-classes",
			Value: "top,groupOfNames",
			Usage: "comma-separated list of classes for group objects",
		},
		&cli.StringFlag{
			Name:  "group-membership-attribute",
			Value: "member",
			Usage: "attribute of the group objects specifying membership",
		},
		&cli.StringFlag{
			Name:  "user-change-type",
			Value: "add",
			Usage: "LDIF changetype for users",
		},
		&cli.StringFlag{
			Name:  "group-change-type",
			Value: "add",
			Usage: "LDIF changetype for groups",
		},
		&cli.StringFlag{
			Name:  "ou-change-type",
			Value: "add",
			Usage: "LDIF changetype for OUs",
		},
		&cli.StringFlag{
			Name:  "buzzword-dataset",
			Usage: fmt.Sprintf(datasetString, "buzzwords", "group"),
		},
		&cli.StringFlag{
			Name:  "department-dataset",
			Usage: fmt.Sprintf(datasetString, "department names", "OU"),
		},
		&cli.StringFlag{
			Name:  "first-name-dataset",
			Usage: fmt.Sprintf(datasetString, "first names", "user"),
		},
		&cli.StringFlag{
			Name:  "last-name-dataset",
			Usage: fmt.Sprintf(datasetString, "last names", "user"),
		},
		&cli.StringFlag{
			Name:  "groups-dataset",
			Usage: fmt.Sprintf(datasetString, "group names", "group"),
		},
	}

	return &cli.Command{
		Name:   "generate",
		Usage:  "generate ldif file",
		Action: generateLdif,
		Flags:  generateFlags,
	}
}

func doGenerate(gconf *GenerateConfig) error {
	tmpl, err := template.New("ldif").Parse(string(templateData))
	if err != nil {
		log.Fatalf("error parsing ldif template: %v", err)
	}

	// nameGen, err := generators.NewNameGenerator(gconf.FirstNameDataset, gconf.LastNameDataset, gconf.DepartmentDataset, gconf.BuzzwordDataset, gconf.GroupsDataset)
	gen, err := ldif.NewGenerator(
		ldif.WithDomain(strings.Join(gconf.Domain, ",")),
	)
	if err != nil {
		log.Fatalf("error creating name generator: %v", err)
	}

	gen.GenerateOrgUnits(gconf.OUs, gconf.OUDepth)

	_, _ = gen.GenerateUsers(gconf.Users)
	gg, _ := gen.GenerateGroups(gconf.Groups)

	renderConfig := RenderConfig{
		GenerateConfig: *gconf,
		Users:          gen.Users,
		Domain:         gconf.Domain,
		OUs:            gen.OrgUnits,
		Groups:         gg,
		Time:           time.Now().Format("2006-01-02T15:04:05-0700"),
	}

	err = tmpl.Execute(os.Stdout, renderConfig)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}

	return nil
}

func generateLdif(ctx *cli.Context) error {
	domainList, err := parseDomain(ctx.String("domain"))
	if err != nil {
		log.Fatalf("error parsing domain: %v", err)
	}

	if ctx.Int("ou-depth") < 1 {
		log.Fatalf("invalid ou depth (<1): %v", ctx.Int("ou-depth"))
	}

	userClassList, err := parseClassList(ctx.String("user-classes"))
	if err != nil {
		log.Fatalf("error parsing user classes: %v", err)
	}

	ouClassList, err := parseClassList(ctx.String("ou-classes"))
	if err != nil {
		log.Fatalf("error parsing ou classes: %v", err)
	}

	groupClassList, err := parseClassList(ctx.String("group-classes"))
	if err != nil {
		log.Fatalf("error parsing group classes: %v", err)
	}

	if ctx.String("user-change-type") == "" {
		log.Fatalf("invalid user change type")
	}

	if ctx.String("group-change-type") == "" {
		log.Fatalf("invalid group change type")
	}

	if ctx.String("ou-change-type") == "" {
		log.Fatalf("invalid ou change type")
	}

	if ctx.String("group-membership-attribute") == "" {
		log.Fatalf("invalid group membership attribute")
	}

	if ok := checkPath(ctx.String("buzzword-dataset")); !ok {
		log.Fatalf("invalid buzzword dataset path: %s", ctx.String("buzzword-dataset"))
	}

	if ok := checkPath(ctx.String("department-dataset")); !ok {
		log.Fatalf("invalid department dataset path: %s", ctx.String("department-dataset"))
	}

	if ok := checkPath(ctx.String("first-name-dataset")); !ok {
		log.Fatalf("invalid first name dataset path: %s", ctx.String("first-name-dataset"))
	}

	if ok := checkPath(ctx.String("last-name-dataset")); !ok {
		log.Fatalf("invalid last name dataset path: %s", ctx.String("last-name-dataset"))
	}

	if ok := checkPath(ctx.String("groups-dataset")); !ok {
		log.Fatalf("invalid groups dataset path: %s", ctx.String("groups-dataset"))
	}

	generateConfig := &GenerateConfig{
		Users:                    ctx.Int("users"),
		Groups:                   ctx.Int("groups"),
		OUs:                      ctx.Int("ous"),
		OUDepth:                  ctx.Int("ou-depth"),
		UserChangeType:           ctx.String("user-change-type"),
		GroupChangeType:          ctx.String("group-change-type"),
		OUChangeType:             ctx.String("ou-change-type"),
		GroupMembershipAttribute: ctx.String("group-membership-attribute"),
		Domain:                   *domainList,
		UserClasses:              *userClassList,
		GroupClasses:             *groupClassList,
		OUClasses:                *ouClassList,
		BuzzwordDataset:          ctx.String("buzzword-dataset"),
		DepartmentDataset:        ctx.String("department-dataset"),
		FirstNameDataset:         ctx.String("first-name-dataset"),
		LastNameDataset:          ctx.String("last-name-dataset"),
		GroupsDataset:            ctx.String("groups-dataset"),
	}

	return doGenerate(generateConfig)
}

func checkPath(path string) bool {
	if path == "" {
		return true // it's not invalid just not set
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseDomain(domain string) (*[]string, error) {
	re := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	if !re.Match([]byte(domain)) {
		return nil, fmt.Errorf("invalid domain %s, regex failed", domain)
	}

	domainList := strings.Split(domain, ".")
	if len(domainList) < 2 {
		return nil, fmt.Errorf("invalid domain %s, split resulted in < 2 segments", domain)
	}

	return &domainList, nil
}

func parseClassList(classes string) (*[]string, error) {
	classList := strings.Split(classes, ",")
	if len(classList) == 0 {
		return nil, fmt.Errorf("invalid class list: %v", classes)
	}
	return &classList, nil
}
