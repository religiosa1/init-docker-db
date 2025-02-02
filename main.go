package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"text/tabwriter"

	"github.com/alecthomas/kong"
)

var CLI struct {
	ContainerName  string `arg:"" optional:"" name:"containerName" help:"name of the database container to be created"`
	Type           string `short:"t" default:"postgres" enum:"postgres,mysql,mssql,mongo" help:"database type"`
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

func main() {
	ctx := kong.Parse(
		&CLI,
		kong.Description("Create a disposable database docker container."),
		kong.Help(helpPrinter),
	)

	if CLI.Version {
		showVersion()
		return
	}

	panic("TODO" + ctx.Command())
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
