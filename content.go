package kevinkit

import "io/fs"

// SkillFS contains the skills shipped with this build.
var SkillFS fs.FS

// SkillRoot is the directory inside SkillFS that contains kk-* directories.
var SkillRoot string
