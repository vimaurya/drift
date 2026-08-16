package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Di-Argus/Drift/pkg/config"
)

func Create(name string) error {
	version := time.Now().Format("20060102150405")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config : %w", err)
	}

	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
		return fmt.Errorf("could not create directory : %w", err)
	}

	files := []struct {
		suffix  string
		content string
	}{
		{suffix: "up", content: ""},
		{suffix: "down", content: ""},
	}

	for _, f := range files {
		fileName := fmt.Sprintf("%s_%s.%s.sql", version, name, f.suffix)

		fullPath := filepath.Join(cfg.Dir, fileName)

		err := os.WriteFile(fullPath, []byte(f.content), 0o644)
		if err != nil {
			return fmt.Errorf("could not create file %s: %w", fileName, err)
		}
		fmt.Printf("created: %s\n", fileName)
	}

	return nil
}
