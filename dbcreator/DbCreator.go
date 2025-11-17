// Package dbcreator provides an inteface for a specific type of DBCreator
// e.g. Posrgres, MYSQL, etc.
package dbcreator

type CreateOptions struct {
	ContainerName string
	Database      string
	User          string
	Password      string
	Port          uint16
	DockerTag     string
	Verbose       bool
	DryRun        bool
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
