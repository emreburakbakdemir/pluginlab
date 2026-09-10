package main

import "fmt"

type ConsoleNotifier struct {
	prefix string
}

func (c *ConsoleNotifier) Name() string { return "console" }

func (c *ConsoleNotifier) Configure(cfg map[string]any) error {
	c.prefix = "[console]"

	for key, val := range cfg {
		switch key {
		case "prefix":
				s, ok := val.(string)
				if !ok { return fmt.Errorf("prefix must be a string, got %T", val) }
				c.prefix = s
		default:
			return fmt.Errorf("unknown option %q", key)
		}
	}
	return nil
}

func (c *ConsoleNotifier) Notify(r BuildResult) error {
	fmt.Printf("%s build %s on %s: %s\n", c.prefix, r.ID, r.Branch, r.Status)
	return nil
}

// in the runtime go executes init automatically after package-level variable are initialized
// and before main(), can have one per file
func init() {
	Register("console", func() Notifier { return &ConsoleNotifier{} } )
}