package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"text/tabwriter"

	"github.com/alecthomas/kong"
	"github.com/religiosa1/init-docker-db/RandomName"
	"github.com/religiosa1/init-docker-db/dbCreator"
	"github.com/religiosa1/init-docker-db/mongo"
	"github.com/religiosa1/init-docker-db/mssql"
	"github.com/religiosa1/init-docker-db/mysql"
	"github.com/religiosa1/init-docker-db/postgres"
)

type CliArgs struct {
	ContainerName  string `arg:"" optional:"" name:"containerName" help:"name of the database container to be created"`
	Type           string `short:"t" help:"database type"`
	User           string `short:"u" help:"database user"`
	Database       string `short:"d" help:"database name"`
	Password       string `short:"p" help:"user's password"`
	Port           uint16 `short:"P" help:"TCP port to which database will be mapped to"`
	Tag            string `short:"T" help:"docker tag to use with the container"`
	NonInteractive bool   `short:"n" help:"exit if any required parameters are missing"`
	Dry            bool   `short:"D" help:"dry run, printing docker command to stdout, without actually running it"`
	Verbose        bool   `short:"v" help:"run with verbose logging"`
	Version        bool   `help:"show version and exit"`
	Help           bool   `short:"h" help:"show help message and exit"`
}

var CLI CliArgs

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

	rl := NewReadline()
	creator, err := getCreator(rl, CLI.Type, CLI.NonInteractive)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	options, err := getOptions(rl, creator, CLI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if rl.hadOutput {
		fmt.Println("") // if we printed anything on console, using empty string as a delimiter
	}

	err = creator.Create(dbCreator.NewShell(options.DryRun, options.Verbose), options)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if !CLI.Dry {
		fmt.Println("Done")
	}
}

func getCreator(rl Readline, dbType string, nonInteractive bool) (dbCreator.DbCreator, error) {
	if dbType != "" {
		return makeCreatorById(dbType)
	}
	if nonInteractive {
		return nil, fmt.Errorf("must supply database type in non-interactive mode")
	}
	const defaultDbType string = "postgres"
	for {
		creatorId, err := rl.Question("database type? [postgres, mysql, mssql, mongo]", defaultDbType)
		if err != nil {
			return nil, err
		}
		creator, err := makeCreatorById(creatorId)
		if err == nil {
			return creator, nil
		}
		fmt.Println(err)
	}
}

func makeCreatorById(dbType string) (dbCreator.DbCreator, error) {
	switch dbType {
	case "postgres":
		return postgres.Creator{}, nil
	case "mssql":
		return mssql.Creator{}, nil
	case "mysql":
		return mysql.Creator{}, nil
	case "mongo":
		return mongo.Creator{}, nil
	}
	return nil, fmt.Errorf("unknown db type '%s'. Must be one of 'postgres', 'mysql', 'mongo'", dbType)
}

func getOptions(rl Readline, creator dbCreator.DbCreator, args CliArgs) (dbCreator.CreateOptions, error) {
	defaultOpts := creator.GetDefaultOpts()
	opts := dbCreator.CreateOptions{
		Database:      args.Database,
		User:          args.User,
		Password:      args.Password,
		ContainerName: args.ContainerName,
		Port:          args.Port,
		Tag:           args.Tag,
		Verbose:       args.Verbose,
		DryRun:        args.Dry,
	}
	// Setting non-interactive defaults
	if opts.Port == 0 {
		opts.Port = defaultOpts.Port
	}
	if opts.Tag == "" {
		opts.Tag = defaultOpts.DockerTag
	}

	// validating existing password first if it's there
	if opts.Password != "" {
		if defaultOpts.Password != "" {
			opts.Password = defaultOpts.Password
		}
		err := creator.IsPasswordValid(opts.Password)
		if err != nil {
			return opts, fmt.Errorf("provided password does not meet the requirements: %w", err)
		}
	}

	// Filling out the rest of missing data in interactive mode if allowed
	if opts.Database == "" {
		if args.NonInteractive {
			return opts, fmt.Errorf("database name is requied in non-interactive mode, but not provided")
		}
		val, err := rl.Question("database name?", "db")
		if err != nil {
			return opts, nil
		}
		opts.Database = val
	}
	if opts.User == "" {
		if args.NonInteractive {
			return opts, fmt.Errorf("db username is requied in non-interactive mode, but not provided")
		}
		val, err := rl.Question("database user?", defaultOpts.User)
		if err != nil {
			return opts, nil
		}
		opts.User = val
	}
	if opts.Password == "" {
		if args.NonInteractive {
			return opts, fmt.Errorf("password is requied in non-interactive mode, but not provided")
		}
		for {
			val, err := rl.Question("database password?", defaultOpts.Password)
			if err != nil {
				return opts, nil
			}
			err = creator.IsPasswordValid(val)
			if err != nil {
				fmt.Println(err)
				continue
			}

			break
		}
	}
	if opts.ContainerName == "" {
		val, err := rl.Question("docker container name?", RandomName.Generate())
		if err != nil {
			return opts, nil
		}
		opts.ContainerName = val
	}

	return opts, nil
}

func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	fmt.Println("\nExamples:")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "  %s\tRun in wizard mode\n", ctx.Model.Name)
	fmt.Fprintf(w, "  %s --dry\tDry-run in wizard mode\n", ctx.Model.Name)
	fmt.Fprintf(w, "  %s -t mssql -u app_user\tCreate a MsSQL database using provided username\n", ctx.Model.Name)

	w.Flush()

	return nil
}

func showVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Build information not available")
		return
	}

	version := info.Main.Version
	if version == "(devel)" {
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
		version += " " + commit
		if dirty {
			version += " dirty"
		}
	}

	fmt.Printf("%s\n", version)
}
