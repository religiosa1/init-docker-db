package dbCreator

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

// Creator capabilities list
type Capabilities struct {
	DatabaseName bool
	UserPassword bool
}

// Creator default options
type DefaultOpts struct {
	User      string
	DockerTag string
	Port      uint16
	Password  string
}

type DbCreator interface {
	GetDefaultOpts() DefaultOpts
	GetCapabilities() Capabilities
	Create(shell Shell, opts CreateOptions) error
	ValidatePassword(password string) error
}
