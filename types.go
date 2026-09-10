package main

// BuildResult is the slice of host state a notifier gets to see
// every field here is a promise you have to keep
type BuildResult struct {
	ID     string
	Branch string
	Status string
}


type Factory func() Notifier

// Notifier is the extension point. Anything implementing this can be plugged 
// into the host without the host knowing its name
type Notifier interface {
	Name() 							string
	Configure(cfg map[string]any) 	error
	Notify(r BuildResult) 			error
}

type notifierEntry struct {
	Name 	string 			`json:"name"`
	Config 	map[string]any 	`json:"config"`
}

type appConfig struct {
	Notifiers []notifierEntry `json:"notifiers"`
}