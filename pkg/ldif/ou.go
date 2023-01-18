package ldif

import (
	"fmt"
	"math/rand"
	"strings"
)

func (g *Generator) randomOU() string {
	return g.OrgUnits[rand.Intn(len(g.OrgUnits))]
}

func (g *Generator) GenerateOrgUnits(topLevelCount int, maxDepth int) []string {
	ouSlice := make([]string, topLevelCount*maxDepth)
	ouMap := make(map[string]int)

	for i := 0; i < topLevelCount; i++ {

		ouString := ""
		for j := 0; j < maxDepth; j++ {
			if j == 0 {
				ouString = g.nameGenerator.Department()
			} else {
				ouString = g.nameGenerator.Department() + ",ou=" + ouString
			}

			if i := ouMap[ouString]; i > 0 {
				oo := strings.Split(ouString, ",")
				ouString = fmt.Sprintf("%s #%d,%s", oo[0], i, strings.Join(oo[1:], ","))
			}

			ouMap[ouString]++
			ouSlice[(i*maxDepth)+j] = ouString
		}
	}

	g.OrgUnits = ouSlice

	return ouSlice
}
