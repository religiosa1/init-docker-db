// Package dbcreator provides an interface for a specific type of DBCreator
// e.g. Postgres, MYSQL, etc.
package dbcreator

import "fmt"

type CreateOptions struct {
	ContainerName string
	Database      string
	User          string
	Password      string
	// host port with optional IP address;
	// see https://docs.docker.com/reference/cli/docker/container/run/#publish
	Ports     []string
	DockerTag string
	Verbose   bool
	DryRun    bool
}

// Capabilities are the list of DBCreator capabilities
type Capabilities struct {
	DatabaseName bool
	UserPassword bool
}

// DefaultOpts are default options for the DBCreator
type DefaultOpts struct {
	User      string
	DockerTag string
	Port      uint16
	Password  string
}

type DBCreator interface {
	GetDefaultOpts() DefaultOpts
	GetCapabilities() Capabilities
	Create(shell Shell, opts CreateOptions) error
	ValidatePassword(password string) error
}

func CreatePortBindingsArgument(containerPort uint16, bindings []string) []string {
	args := make([]string, len(bindings)*2)
	for i := range bindings {
		args[i*2] = "-p"
		args[i*2+1] = fmt.Sprintf("%s:%d", bindings[i], containerPort)
	}
	return args
}
