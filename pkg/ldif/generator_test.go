package ldif

import "testing"

func BenchmarkGenerator_GenerateN(b *testing.B) {
	gen, err := NewGenerator()
	if err != nil {
		b.Errorf("cannot create generator: %s", err.Error())
		b.FailNow()
	}

	gen.GenerateOrgUnits(3, 1)
	gen.GenerateUsers(100)

	for i := 0; i < b.N; i++ {
		_, err := gen.GenerateGroups(100)
		if err != nil {
			b.Errorf("group generation error: %s", err.Error())
		}
	}
}
