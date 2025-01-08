package sdblueprint

import (
	"github.com/gaorx/stardust6/sdfile"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

func getGoModDir() string {
	goModFn := sdfile.FirstExistsAbs(
		"go.mod",
		"../go.mod",
		"../../go.mod",
		"../../../go.mod",
		"../../../../go.mod",
		"../../../../../go.mod",
	)
	if goModFn == "" {
		wd, _ := os.Getwd()
		if wd != "" {
			return wd
		}
	}
	return filepath.Dir(goModFn)
}

func getGoModName() string {
	dir := getGoModDir()
	if dir == "" {
		return ""
	}
	goModFn := filepath.Join(dir, "go.mod")
	goModData, err := sdfile.ReadBytes(goModFn)
	if err != nil {
		return ""
	}
	goModF, err := modfile.ParseLax(goModFn, goModData, nil)
	if err != nil {
		return ""
	}
	if goModF == nil || goModF.Module == nil {
		return ""
	}
	return goModF.Module.Mod.Path
}
