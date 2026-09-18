//go:build !e2ea && !e2eb

package kevinkit

import "embed"

//go:embed all:skills/kk-*
var productionSkills embed.FS

func init() {
	SkillFS = productionSkills
	SkillRoot = "skills"
}
