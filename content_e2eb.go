//go:build e2eb

package kevinkit

import "embed"

//go:embed all:internal/e2e/testdata/payload-b/kk-*
var fixtureSkillsB embed.FS

func init() {
	SkillFS = fixtureSkillsB
	SkillRoot = "internal/e2e/testdata/payload-b"
}
