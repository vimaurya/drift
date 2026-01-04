package migration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vimaurya/gomigrate/internal/config"
)

func Create(name string) error {
	version := time.Now().Format("20060102150405")

	data, err := os.ReadFile(".gomigrate.json")

	if err != nil {
		return fmt.Errorf("failed to locate migration directory")
	}

	var cfg config.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to Unmarshal config file %w", err)
	}

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return fmt.Errorf("could not create directory : %w", err)
	}

	files := []struct {
		suffix  string
		content string
	}{
		{suffix: "up", content: "--Write your up migration here"},
		{suffix: "down", content: "--Write your down migration here"},
	}

	for _, f := range files {
		fileName := fmt.Sprintf("%s_%s.%s.sql", version, name, f.suffix)

		fullPath := filepath.Join(cfg.Dir, fileName)

		err := os.WriteFile(fullPath, []byte(f.content), 0644)
		if err != nil {
			return fmt.Errorf("could not create file %s: %w", fileName, err)
		}
		fmt.Printf("created: %s\n", fileName)
	}

	return nil
}
