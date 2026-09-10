package main

import (
	"sort"
	"fmt"
)

var registry = map[string]Factory{}

func Register(name string, newNotifier Factory) {
	if _, exists := registry[name]; exists {
		//this panic is a startup failure, not runtime
		panic(fmt.Sprintf("notifier already registered: %q", name))
	}
	registry[name] = newNotifier
}

func Lookup(name string) (Factory, bool) {
	newNotifier, ok := registry[name]
	return newNotifier, ok
}

func Available() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
} 


