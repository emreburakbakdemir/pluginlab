package main

import (
	"fmt"
	"os"
)

type FileNotifier struct {
	path string
}

func (f *FileNotifier) Name() string { return "file" }

func (f *FileNotifier) Configure(cfg map[string]any) error {
	for key, val := range cfg {
		switch key {
		case "path":
			s, ok := val.(string)
			if !ok { return fmt.Errorf("path must be a string, got %T", val) }
			f.path = s
		default:
			return fmt.Errorf("unknown option %q", key)
		}
	}
	if f.path == "" {
		return fmt.Errorf("option %q is required", "path")
	}
	return nil
}

func (f *FileNotifier) Notify(r BuildResult) error {
	file, err := os.OpenFile(f.path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("file notifier: %w", err)
	}
	defer file.Close()

	line := fmt.Sprintf("build %s %s %s\n", r.ID, r.Branch, r.Status)
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("file notifier: %w", err)
	}

	return nil
}

func init() {
	Register("file", func() Notifier{ return &FileNotifier{} } )
}