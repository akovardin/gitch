package static

import "embed"

//go:embed js/*
//go:embed css/*
//go:embed img/*
var FS embed.FS