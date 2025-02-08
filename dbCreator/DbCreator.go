package dbCreator

type CreateOptions struct {
	ContainerName string
	Database      string
	User          string
	Password      string
	Port          uint16
	Tag           string
	Verbose       bool
	DryRun        bool
}

type DefaultOpts struct {
	User      string
	DockerTag string
	Port      uint16
	Password  string
}

type DbCreator interface {
	GetDefaultOpts() DefaultOpts
	Create(shell Shell, opts CreateOptions) error
	ValidatePassword(password string) error
}
