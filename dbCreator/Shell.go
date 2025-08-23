package dbCreator

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Format DockerEnv variable for passing to docker options
func DockerEnv(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// Child process runner with output verbosity flag
type Shell struct {
	dryRun  bool
	verbose bool
}

// Create new shell instance
func NewShell(dryRun bool, verbose bool) Shell {
	return Shell{
		dryRun:  dryRun,
		verbose: verbose,
	}
}

// Create new shell instance
func (sh Shell) RunWithOutput(name string, args ...string) (string, error) {
	if sh.dryRun || sh.verbose {
		fmt.Println(makeShellCmdString(name, args...))
	}
	if sh.dryRun {
		return "", nil
	}
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Run a child process, printing its outputs to Stdout/Stderr only in the verbose mode
func (sh Shell) RunSilent(name string, args ...string) error {
	if sh.dryRun || sh.verbose {
		fmt.Println(makeShellCmdString(name, args...))
	}
	if sh.dryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	if !sh.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// Run a child process, printing its output to Stdout/Stderr
func (sh Shell) Run(name string, args ...string) error {
	if sh.dryRun || sh.verbose {
		fmt.Println(makeShellCmdString(name, args...))
	}
	if sh.dryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var pattern *regexp.Regexp

// Quote a shell argument
func Quote(cmd string) string {
	if pattern == nil {
		pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
	}
	if pattern.MatchString(cmd) {
		return "'" + strings.ReplaceAll(cmd, "'", "'\"'\"'") + "'"
	}

	return cmd
}

func makeShellCmdString(name string, args ...string) string {
	var sb strings.Builder
	sb.WriteString(name)

	for _, arg := range args {
		sb.WriteString(" ")
		sb.WriteString(Quote(arg))
	}

	return sb.String()
}
