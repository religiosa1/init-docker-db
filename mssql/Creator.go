package mssql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type Creator struct{}

const port uint16 = 1433

func (c Creator) GetDefaultOpts() dbCreator.DefaultOpts {
	return dbCreator.DefaultOpts{
		Port:      port,
		User:      "mssql",
		DockerTag: "2022-latest",
		Password:  "Password12",
	}
}

func (c Creator) Create(shell dbCreator.Shell, opts dbCreator.CreateOptions) error {
	// https://mcr.microsoft.com/product/mssql/server/about
	shellOutput, err := shell.RunWithOutput("docker", "run", "-e", "ACCEPT_EULA=Y",
		"--name", opts.ContainerName,
		"--hostname", opts.ContainerName,
		"-e", dbCreator.DockerEnv("MSSQL_SA_PASSWORD", opts.Password),
		"-p", fmt.Sprintf("%d:%d", port, opts.Port),
		"-d", fmt.Sprintf("mcr.microsoft.com/mssql/server:%s", opts.DockerTag),
	)
	fmt.Println(shellOutput)
	if err != nil {
		return err
	}
	contId := strings.TrimSpace(shellOutput)

	sql := SqlInContainerRunner{
		shell:    &shell,
		contId:   contId,
		database: opts.Database,
		password: opts.Password,
		verbose:  opts.Verbose,
	}

	v := VerboseLogger{opts.Verbose}
	fmt.Println("Waiting for db to be up and running...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// predelaying waiting for 1 seconds, as there's no way MsSQL can launch that
	// fast, and connectivity timeouts take quite some time to resolve.
	waitOpts := WaitForOpts{PreDelay: 1000}

	err = waitFor(ctx, func() error {
		start := time.Now()
		// can be SELECT SERVERPROPERTY('ProductVersion'), but its ouput is too long
		err := sql.Run("SELECT 1")
		end := time.Since(start)
		v.Log("sql health check duration", end)
		return err
	}, waitOpts)
	if err != nil {
		return fmt.Errorf("failed to wait for the database to be operational: %w", err)
	}

	fmt.Println("Creating the database and required data...")

	escapedDbName, err := escapeId(opts.Database)
	if err != nil {
		return fmt.Errorf("error escaping the database name: %w", err)
	}

	err = sql.Run(fmt.Sprintf("CREATE DATABASE %s", escapedDbName))
	if err != nil {
		return err
	}

	v.Log("Creating login")

	escapedUser, err := escapeUser(opts.User)
	if err != nil {
		return fmt.Errorf("error escaping the username: %w", err)
	}

	err = sql.Run(fmt.Sprintf("CREATE LOGIN %s WITH PASSWORD = %s", escapedUser, escapeStr(opts.Password)))
	if err != nil {
		return err
	}

	v.Log("Creating user")
	err = sql.RunInDb(fmt.Sprintf(`create user %s for login %s`, escapedUser, escapedUser))
	if err != nil {
		return err
	}

	// To check available roles: Select	[name] From sysusers Where issqlrole = 1
	v.Log("Adding required permissions")
	err = sql.RunInDb(fmt.Sprintf("ALTER ROLE db_owner ADD MEMBER %s", escapedUser))
	if err != nil {
		return err
	}
	return nil
}

var ErrPasswordEmpty error = errors.New("password can't be empty")
var ErrPasswordTooShort error = errors.New("password is too short (must be at least 10 chars)")
var ErrPasswordTooSimple error = errors.New(
	"password doesn't meet the complexity requirements " +
		"(must contain 3 out of 4 char types: lowercase char, uppercase char, digit, nonalphanumeric)",
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
