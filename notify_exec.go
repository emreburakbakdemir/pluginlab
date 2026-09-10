package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
	"strings"
)

type ExecNotifier struct {
	command		string
	args		[]string
	transport	string // "env" or "json"
	timeout		time.Duration
} 

func (e *ExecNotifier) Name() string { return "exec" }

func (e *ExecNotifier) Configure(cfg map[string]any) error {
	// default values
	e.transport = "env"
	e.timeout = 10 * time.Second

	for key, val := range cfg {
		switch key {

		case "command":
			s, ok := val.(string)
			if !ok { return fmt.Errorf("command must be a string, got %T", val) }
			e.command = s
		
		case "args":
			raw, ok := val.([]any)
			if !ok { return fmt.Errorf("args must be a list, got %T", val) }
			for i, item := range raw {
				s, ok := item.(string)
				if !ok { return fmt.Errorf("args[%d] must be a string, got %T", i, item) }
				e.args = append(e.args, s)
			}
		
		case "transport":
			s, ok := val.(string)
			if !ok { return fmt.Errorf("transport must be a string, got %T", val) }
			if s != "env" && s != "json" { return fmt.Errorf("transport must be %q or %q, got %q",  "env", "json", s) }
			e.transport = s

		case "timeout":
			f, ok := val.(float64)
			if !ok { return fmt.Errorf("timeout must be a number, got %T", val) }
			if f <= 0 { return fmt.Errorf("timeout must positive, got %v", f) }
			e.timeout = time.Duration(f * float64(time.Second))
		
		default:
			return fmt.Errorf("unknown option %q", key)	
		}
	}

	if e.command == "" { return fmt.Errorf("option %q is required", "command") }
	return nil
}

func (e *ExecNotifier) Notify(r BuildResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.command, e.args...)

	// if we leave cmd.Env nil, child inherits entire env
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	/*
		env transport:
			stdout = log noise <- plugin can echo freely
			stderr = log noise
		
		json transport:
			stdout = the response <- must parse, any stray echo breaks it
			stderr = the only log channel
	
		json-transport plugins must log only to stderr.

	*/
	switch e.transport {
	case "env":
		cmd.Env = append(cmd.Env, 
						"BUILD_ID=" + r.ID,
						"BUILD_BRANCH=" + r.Branch,
						"BUILD_STATUS=" + r.Status,)
	
	case "json":
		payload, err := json.Marshal(r)
		if err != nil { return fmt.Errorf("marshal payload: %w", err) }
		cmd.Stdin = bytes.NewReader(payload) // could use cmdStdinPipe(), however must
		// handle the EOF for input=$(cat) for example.
	}

	err := cmd.Run()

	// plugin outputs printed
	if out := strings.TrimSpace(stdout.String()); out != "" {
		for _, line := range strings.Split(out, "\n") {
			fmt.Printf("[%s] %s\n", e.command, line)
		}
	}

	if errOut := strings.TrimSpace(stderr.String()); errOut != "" {
		for _, line := range strings.Split(errOut, "\n") {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", e.command, line)
		}
	}

	// context timeout check
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timed out after %s", e.timeout)
	}

	if err != nil {
		var exitErr *exec.ExitError
		/*
		err.(*exec.ExitError) checks only the outermost error. If anything wrapped it, 
		the assertion silently fails and you misreport a legitimate exit code as 
		"could not start." errors.As walks the wrap chain
		*/
		if errors.As(err, &exitErr) {
			return fmt.Errorf("exit %d: %s", exitErr.ExitCode(), tail(stderr.String()))
		}
		return fmt.Errorf("could not start %q: %w", e.command, err)
	}

	// json transport has 2 success conditions, the process succeeded and the response says so
	if e.transport == "json" {
		var resp struct {
			OK 			bool	`json:"ok"`
			Message		string	`json:"message"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
			return fmt.Errorf("bad response: %w (got %q)", err, tail(stdout.String()))
		}
		if !resp.OK {
			return fmt.Errorf("plugin reported failure: %s", resp.Message)
		}
	}
	return nil
}

func tail(s string) string {
	s = strings.TrimSpace(s)
	if s == "" { return "(no output)" }
	if len(s) > 200 { return "..." + s[len(s) - 200:] }
	return s
}

func init() {
	Register("exec", func() Notifier { return &ExecNotifier{} })
}