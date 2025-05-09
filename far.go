package far

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func FindAndReplace(root, old, new string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		original, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		replaced := bytes.ReplaceAll(original, []byte(old), []byte(new))

		// Avoid unneccessary writes
		if bytes.Equal(original, replaced) {
			slog.Info("skip", slog.String("path", path))
			return nil
		}

		err = os.WriteFile(path, replaced, d.Type().Perm())
		if err != nil {
			return fmt.Errorf("failed to write out to %s: %w", path, err)
		}

		slog.Info("done", slog.String("path", path))
		return nil
	})
}
