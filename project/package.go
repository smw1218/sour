package project

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

func ReadPackage() (string, error) {
	mf, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed reading go.mod: %w", err)
	}
	gomod, err := modfile.Parse("go.mod", mf, nil)
	if err != nil {
		return "", err
	}
	return gomod.Module.Mod.Path, nil
}
