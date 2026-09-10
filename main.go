package main

import (
	"fmt"
	"os"
	"strings"
)

// compile-time registration with runtime configuration

func resolve(cfg *appConfig) ([]Notifier, error) {
	active := make([]Notifier, 0, len(cfg.Notifiers))

	for _, entry := range cfg.Notifiers {
		newNotifier, ok := Lookup(entry.Name)
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: unknown notifier %q, skipping (available: %s)\n", 
						entry.Name, strings.Join(Available(), ", "))
			continue
		}
		
		n := newNotifier()

		if err := n.Configure(entry.Config); err != nil {
			return nil, fmt.Errorf("configuring %q: %w", entry.Name, err)
		}
		active = append(active, n)
	}
	return active, nil
}

func main() {
	path := "config.json"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	cfg, err := loadConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	active, err := resolve(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	result := BuildResult{
		ID:		"build-1042",
		Branch: "main",
		Status: "success",
	}

	failed := 0

	for _, n := range active {
		if err := n.Notify(result); err != nil {
			fmt.Fprintf(os.Stderr, "notifer %q failed: %v\n", n.Name(), err)
			failed++
		}
	}

	fmt.Printf("%d notifiers ran, %d notifiers failed\n", len(active), failed)
	if failed > 0 { os.Exit(1) }

}