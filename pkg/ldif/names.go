package ldif

import (
	"bufio"
	"embed"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

//go:embed data
var content embed.FS

func init() {
	rand.Seed(time.Now().UnixMicro())
}

type NameGenerator struct {
	firstNames  []string
	lastNames   []string
	departments []string
	buzzwords   []string
	groups      []string

	re *regexp.Regexp
}

func NewNameGenerator() (*NameGenerator, error) {
	ng := NameGenerator{}
	ng.re = regexp.MustCompile(`[\w ]+`)

	var err error

	ng.firstNames, err = ng.readWords("first_names")
	if err != nil {
		return nil, fmt.Errorf("cannot read firstNames: %w", err)
	}

	ng.lastNames, err = ng.readWords("last_names")
	if err != nil {
		return nil, fmt.Errorf("cannot read lastNames: %w", err)
	}

	ng.departments, err = ng.readWords("departments")
	if err != nil {
		return nil, fmt.Errorf("cannot read departments: %w", err)
	}

	ng.buzzwords, err = ng.readWords("buzzwords")
	if err != nil {
		return nil, fmt.Errorf("cannot read buzzrowds: %w", err)
	}

	ng.groups, err = ng.readWords("groups")
	if err != nil {
		return nil, fmt.Errorf("cannot read groups: %w", err)
	}

	return &ng, nil
}

func (n *NameGenerator) FirstName() string {
	index := rand.Intn(len(n.firstNames))

	return n.firstNames[index]
}

func (n *NameGenerator) LastName() string {
	index := rand.Intn(len(n.lastNames))

	return n.lastNames[index]
}

func (n *NameGenerator) Department() string {
	index := rand.Intn(len(n.departments))

	return n.departments[index]
}

func (n *NameGenerator) Group() string {
	buzzwordIndex := rand.Intn(len(n.buzzwords))
	groupIndex := rand.Intn(len(n.groups))

	return n.buzzwords[buzzwordIndex] + " " + n.groups[groupIndex]
}

func (n *NameGenerator) readWords(fileName string) ([]string, error) {
	var lines []string

	f, err := content.Open(fmt.Sprintf("data/%s.txt", fileName))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Get rid of empty lines
		if len(line) == 0 {
			continue
		}

		if len(n.re.FindString(line)) < len(line) {
			// this means that the regex matched less than the total string
			continue
		}

		lines = append(lines, line)
	}

	return lines, scanner.Err()
}
