package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Create(name, dir string) error {
	version := time.Now().Format("20060102150405")

	if dir==""{
		dir = "DB_Migrations"
	}

	if err:= os.MkdirAll(dir, 0755); err!=nil{
		return fmt.Errorf("could not create directory : %w", err)
	}


	files := []struct{
		suffix string
		content string
	}{
		{suffix: "up", content: "--Write your up migration here"},
		{suffix: "down", content: "--Write your down migration here"},
	}

	for _, f := range files {
		fileName := fmt.Sprintf("%s_%s.%s.sql", version, name, f.suffix)

		fullPath := filepath.Join(dir, fileName)

		err := os.WriteFile(fullPath, []byte(f.content), 0644)
		if err!=nil{
			return fmt.Errorf("could not create file %s: %w", fileName, err)
		}
		fmt.Printf("created: %s\n", fileName)
	}

	return nil
}
