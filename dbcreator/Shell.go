package dbcreator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// DockerEnv formats DockerEnv variable for passing to docker options
func DockerEnv(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// Shell is a child process runner with output verbosity flag
type Shell struct {
	dryRun  bool
	verbose bool
}

// NewShell creates a  new shell instance
func NewShell(dryRun bool, verbose bool) Shell {
	return Shell{
		dryRun:  dryRun,
		verbose: verbose,
	}
}

// RunWithOutput runs a new shell instance capturing it's stdtout as a return value
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

// RunWithTeeOutput runs a child process, streaming its output to Stdout/Stderr while also capturing it as a return value
func (sh Shell) RunWithTeeOutput(name string, args ...string) (string, error) {
	if sh.dryRun || sh.verbose {
		fmt.Println(makeShellCmdString(name, args...))
	}
	if sh.dryRun {
		return "", nil
	}
	cmd := exec.Command(name, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

	err := cmd.Run()

	// Combine stdout and stderr for the return value, similar to CombinedOutput
	combined := outBuf.String() + errBuf.String()
	return combined, err
}

// RunSilent runs a child process, printing its outputs to Stdout/Stderr only in the verbose mode
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
