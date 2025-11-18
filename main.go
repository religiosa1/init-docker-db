package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"text/tabwriter"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/huh"
	"github.com/religiosa1/init-docker-db/creators/mongo"
	"github.com/religiosa1/init-docker-db/creators/mssql"
	"github.com/religiosa1/init-docker-db/creators/mysql"
	"github.com/religiosa1/init-docker-db/creators/postgres"
	"github.com/religiosa1/init-docker-db/creators/redis"
	"github.com/religiosa1/init-docker-db/dbcreator"
	"github.com/religiosa1/init-docker-db/randomname"
)

var ldVersion = "" // Version set by -ldflags during the Taskfile build

type CliArgs struct {
	ContainerName  string `arg:"" optional:"" name:"containerName" help:"name of the database container to be created"`
	Type           string `short:"t" help:"database type"`
	User           string `short:"u" help:"database user"`
	Database       string `short:"d" help:"database name"`
	Password       string `short:"p" help:"user's password"`
	Port           string `short:"P" help:"port with optional IP address to which database will be mapped to"`
	Public         bool   `help:"expose default port to outside world by mapping to 0.0.0.0 IP address"`
	Tag            string `short:"T" help:"docker tag to use with the container"`
	NonInteractive bool   `short:"n" help:"exit if any required parameters are missing"`
	Dry            bool   `short:"D" help:"dry run, printing docker command to stdout, without actually running it"`
	Verbose        bool   `short:"v" help:"run with verbose logging"`
	Version        bool   `help:"show version and exit"`
	Help           bool   `short:"h" help:"show help message and exit"`
}

var CLI CliArgs

type ExitStatus int

const (
	ExitStatusSuccess ExitStatus = iota
	ExitStatusDockerNotFound
	ExitStatusFailedToGetCreator
	ExitStatusFailedToCreateContainer
)

func main() {
	kong.Parse(
		&CLI,
		kong.Description("Create a disposable database docker container."),
		kong.Help(helpPrinter),
	)

	if CLI.Version {
		showVersion()
		return
	}

	// Check if docker is available in PATH (skip for dry-run mode)
	if !CLI.Dry {
		_, err := exec.LookPath("docker")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: 'docker' command not found in PATH.")
			fmt.Fprintln(os.Stderr, "Please install Docker and ensure it's available in your PATH.")
			fmt.Fprintln(os.Stderr, "Run with --dry flag to see commands without requiring Docker.")
			os.Exit(int(ExitStatusDockerNotFound))
		}
	}

	creator, err := getCreator(CLI.Type, CLI.NonInteractive)
	if err != nil {
		fmt.Println(err)
		os.Exit(int(ExitStatusFailedToGetCreator))
	}
	options, err := getOptions(creator, CLI)
	if err != nil {
		fmt.Println(err)
		os.Exit(int(ExitStatusFailedToGetCreator))
	}

	err = creator.Create(dbcreator.NewShell(options.DryRun, options.Verbose), options)
	if err != nil {
		fmt.Println(err)
		os.Exit(int(ExitStatusFailedToCreateContainer))
	}
}

var theme *huh.Theme = huh.ThemeBase16()

func getCreator(dbType string, nonInteractive bool) (dbcreator.DBCreator, error) {
	if dbType != "" {
		return makeCreatorByID(dbType)
	}
	if nonInteractive {
		return nil, fmt.Errorf("must supply database type in non-interactive mode")
	}
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Database type?").
			Options(
				huh.NewOption("postgres", "postgres"),
				huh.NewOption("mssql", "mssql"),
				huh.NewOption("mysql", "mysql"),
				huh.NewOption("mongo", "mongo"),
				huh.NewOption("redis", "redis"),
			).
			Value(&dbType),
	)).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}
	return makeCreatorByID(dbType)
}

func makeCreatorByID(dbType string) (dbcreator.DBCreator, error) {
	switch dbType {
	case "postgres":
		return postgres.Creator{}, nil
	case "mssql":
		return mssql.Creator{}, nil
	case "mysql":
		return mysql.Creator{}, nil
	case "mongo":
		return mongo.Creator{}, nil
	case "redis":
		return redis.Creator{}, nil
	}
	return nil, fmt.Errorf("unknown db type '%s'. Must be one of 'postgres', 'mssql', 'mysql', 'mongo', or 'redis'", dbType)
}

const defaultDBName = "db"

func runWizard(
	capabilities dbcreator.Capabilities,
	validatePassword func(string) error,
	defaultContainerName string,
	defaults dbcreator.DefaultOpts,
	opts *dbcreator.CreateOptions,
) error {
	// We're not setting any values for the fields, opting out for placeholder --
	// in case user wants to modify the default value, they don't need to erase the current value.
	// On a cons side, we need to explicitly check for values afterwards. We're not doing that in
	// the runWizard, as this has to be done for non-interactive mode as well anyway.

	fields := make([]huh.Field, 0)
	if capabilities.DatabaseName && opts.Database == "" {
		fields = append(fields, huh.NewInput().
			Title("Database Name?").
			Placeholder(defaultDBName).
			Value(&opts.Database),
		)
	}
	if capabilities.UserPassword {
		if opts.User == "" {
			fields = append(fields, huh.NewInput().
				Title("Database User?").
				Placeholder(defaults.User).
				Value(&opts.User),
			)
		}
		if opts.Password == "" {
			fields = append(fields, huh.NewInput().
				Title("Database password?").
				EchoMode(huh.EchoModePassword).
				Validate(func(val string) error {
					// if value is empty we're omitting the validation as the default value will be set later
					if val == "" {
						return nil
					}
					return validatePassword(val)
				}).
				Placeholder(defaults.Password).
				Value(&opts.Password),
			)
		}
	}
	if opts.ContainerName == "" {
		fields = append(fields, huh.NewInput().
			Title("Docker Container Name?").
			Validate(ValidateContainerName).
			Placeholder(defaultContainerName).
			Value(&opts.ContainerName),
		)
	}

	return huh.NewForm(huh.NewGroup(fields...).WithTheme(theme)).Run()
}

func getOptions(creator dbcreator.DBCreator, args CliArgs) (dbcreator.CreateOptions, error) {
	capabilities := creator.GetCapabilities()
	defaultOpts := creator.GetDefaultOpts()
	opts := dbcreator.CreateOptions{
		Database:      args.Database,
		User:          args.User,
		Password:      args.Password,
		ContainerName: args.ContainerName,
		Port:          args.Port,
		DockerTag:     args.Tag,
		Verbose:       args.Verbose,
		DryRun:        args.Dry,
	}
	// Setting non-interactive-only defaults
	if opts.Port == "" {
		if args.Public {
			opts.Port = fmt.Sprintf("%d", defaultOpts.Port)
		} else {
			opts.Port = fmt.Sprintf("127.0.0.1:%d", defaultOpts.Port)
		}
	}
	if opts.DockerTag == "" {
		opts.DockerTag = defaultOpts.DockerTag
	}

	// validating existing password first if it's there for early exit
	if opts.Password != "" {
		err := creator.ValidatePassword(opts.Password)
		if err != nil {
			return opts, fmt.Errorf("provided password does not meet the requirements: %w", err)
		}
	}

	randomContainerName := randomname.Generate()

	if !args.NonInteractive {
		err := runWizard(capabilities, creator.ValidatePassword, randomContainerName, defaultOpts, &opts)
		if err != nil {
			return opts, fmt.Errorf("error running the wizard: %w", err)
		}
	}

	// Setting default values
	if opts.User == "" && defaultOpts.User != "" {
		opts.User = defaultOpts.User
	}
	if opts.Password == "" && defaultOpts.Password != "" {
		opts.Password = defaultOpts.Password
	}
	if opts.ContainerName == "" {
		opts.ContainerName = randomContainerName
	}

	if capabilities.DatabaseName && opts.Database == "" {
		opts.Database = defaultDBName
	}
	if !capabilities.DatabaseName && opts.Database != "" {
		fmt.Fprintln(os.Stderr, "This DB type doesn't support database name, so provided argument is ignored")
	}

	if capabilities.UserPassword {
		if args.NonInteractive {
			if opts.User == "" {
				return opts, fmt.Errorf("db username is required in non-interactive mode, but not provided")
			}
			if defaultOpts.Password == "" {
				return opts, fmt.Errorf("password is required in non-interactive mode, but not provided")
			}
		}
	} else {
		if opts.User != "" {
			fmt.Fprintln(os.Stderr, "This DB type doesn't support user/password for its auth, so provided username argument is ignored")
		}
		if opts.Password != "" {
			fmt.Fprintln(os.Stderr, "This DB type doesn't support user/password for its auth, so provided password argument is ignored")
		}
	}

	return opts, nil
}

var (
	containerFirstChar   *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9]`)
	containerNamePattern *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

// ValidateContainerName validates container name in wizard mode.
// We're omitting empty string, as it will be autogenerated. In non-interactive
// mode we're not validating name at all, allowing docker to error out instead
func ValidateContainerName(name string) error {
	if name == "" {
		return nil
	}
	if len(name) > 128 {
		return errors.New("maximum length is 128 characters")
	}
	if !containerNamePattern.MatchString(name) {
		return errors.New("can only contain alphanumeric characters from 0 to 9, A to Z, a to z, and the _ and - characters")
	}
	if !containerFirstChar.MatchString(name) {
		return errors.New("must start with an alphanumeric character")
	}
	return nil
}

func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	fmt.Println("\nExamples:")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintf(w, "  %s\tRun in wizard mode\n", ctx.Model.Name)
	_, _ = fmt.Fprintf(w, "  %s --dry\tDry-run in wizard mode\n", ctx.Model.Name)
	_, _ = fmt.Fprintf(w, "  %s -t mssql -u app_user\tCreate a MsSQL database using provided username\n", ctx.Model.Name)

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func showVersion() {
	if ldVersion != "" {
		fmt.Printf("%s\n", ldVersion)
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Build information not available")
		return
	}
	buildInfoVersion := info.Main.Version
	if buildInfoVersion == "(devel)" {
		var commit string
		var dirty bool
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
			if setting.Key == "vcs.modified" {
				dirty = setting.Value == "true"
			}
		}
		buildInfoVersion += " " + commit
		if dirty {
			buildInfoVersion += " dirty"
		}
	}

	fmt.Printf("%s\n", buildInfoVersion)
}
