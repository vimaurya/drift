package migrations

import (
	"os"
	"sort"
	"strings"
)

func GetAvailableMigrations(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err!=nil{
		return nil, err
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql"){
			upFiles = append(upFiles, f.Name())
		}
	}

	sort.Strings(upFiles)
	return upFiles, nil
}
