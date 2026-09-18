//go:build e2ea

package kevinkit

import "embed"

//go:embed all:internal/e2e/testdata/payload-a/kk-*
var fixtureSkillsA embed.FS

func init() {
	SkillFS = fixtureSkillsA
	SkillRoot = "internal/e2e/testdata/payload-a"
}
