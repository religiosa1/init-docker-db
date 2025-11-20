// Package mssql implements DBCreator interface for MS SQL Server
package mssql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/religiosa1/init-docker-db/dbcreator"
	"github.com/religiosa1/init-docker-db/wait"
)

type Creator struct{}

const port uint16 = 1433

func (c Creator) GetDefaultOpts() dbcreator.DefaultOpts {
	return dbcreator.DefaultOpts{
		Port:      port,
		User:      "mssql",
		DockerTag: "2022-latest",
		Password:  "Password12",
	}
}

func (c Creator) GetCapabilities() dbcreator.Capabilities {
	return dbcreator.Capabilities{
		DatabaseName: true,
		UserPassword: true,
	}
}

func (c Creator) Create(shell dbcreator.Shell, opts dbcreator.CreateOptions) error {
	// https://mcr.microsoft.com/product/mssql/server/about
	args := []string{
		"run", "-e", "ACCEPT_EULA=Y",
		"--name", opts.ContainerName,
		"--hostname", opts.ContainerName,
		"-e", dbcreator.DockerEnv("MSSQL_SA_PASSWORD", opts.Password),
	}
	args = append(args, dbcreator.CreatePortBindingsArgument(port, opts.Ports)...)
	args = append(args, "-d", fmt.Sprintf("mcr.microsoft.com/mssql/server:%s", opts.DockerTag))
	shellOutput, err := shell.RunWithTeeOutput("docker", args...)
	if err != nil {
		return err
	}
	contID := strings.TrimSpace(shellOutput)

	sql := SQLInContainerRunner{
		shell:    &shell,
		contID:   contID,
		database: opts.Database,
		password: opts.Password,
		verbose:  opts.Verbose,
	}

	v := NewProgressLogger(opts.Verbose)
	defer v.Done()
	v.LogState("Waiting for db to be up and running...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// predelaying waiting for 1 seconds, as there's no way MsSQL can launch that
	// fast, and connectivity timeouts take quite some time to resolve.
	waitOpts := wait.Opts{PreDelay: 1000}

	err = wait.For(ctx, func() error {
		start := time.Now()
		err := sql.RunSilent("SELECT SERVERPROPERTY('ProductVersion')")
		end := time.Since(start)
		v.LogVerbose("sql health check duration", end)
		return err
	}, waitOpts)
	if err != nil {
		return fmt.Errorf("failed to wait for the database to be operational: %w", err)
	}

	v.LogState("Creating the database and required data...")

	escapedDBName, err := escapeID(opts.Database)
	if err != nil {
		return fmt.Errorf("error escaping the database name: %w", err)
	}

	err = sql.Run(fmt.Sprintf("CREATE DATABASE %s", escapedDBName))
	if err != nil {
		return err
	}

	v.LogState("Creating login")

	escapedUser, err := escapeUser(opts.User)
	if err != nil {
		return fmt.Errorf("error escaping the username: %w", err)
	}

	err = sql.Run(fmt.Sprintf("CREATE LOGIN %s WITH PASSWORD = %s", escapedUser, escapeStr(opts.Password)))
	if err != nil {
		return err
	}

	v.LogState("Creating user")
	err = sql.RunInDB(fmt.Sprintf(`create user %s for login %s`, escapedUser, escapedUser))
	if err != nil {
		return err
	}

	// To check available roles: Select	[name] From sysusers Where issqlrole = 1
	v.LogState("Adding required permissions")
	err = sql.RunInDB(fmt.Sprintf("ALTER ROLE db_owner ADD MEMBER %s", escapedUser))
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrPasswordEmpty     error = errors.New("password can't be empty")
	ErrPasswordTooShort  error = errors.New("password is too short (must be at least 10 chars)")
	ErrPasswordTooSimple error = errors.New(
		"password doesn't meet the complexity requirements " +
			"(must contain 3 out of 4 char types: lowercase char, uppercase char, digit, non-alphanumeric)",
	)
)

func (c Creator) ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}
	if len(password) < 10 {
		return ErrPasswordTooShort
	}
	if !isPasswordComplexEnough(password) {
		return ErrPasswordTooSimple
	}
	return nil
}
