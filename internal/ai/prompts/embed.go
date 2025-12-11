package prompts

import _ "embed"

//go:embed system.md
var System string

//go:embed analyze.md
var Analyze string

//go:embed error_analysis.md
var ErrorAnalysis string

//go:embed visual.md
var Visual string
