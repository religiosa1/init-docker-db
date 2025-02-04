package mssql

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type Creator struct{}

const port uint16 = 1433

func (c Creator) GetDefaultOpts() dbCreator.DefaultOpts {
	return dbCreator.DefaultOpts{
		Port:      port,
		User:      "mysql",
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
		"-d", fmt.Sprintf("mysql:%s", opts.Tag),
	)
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

	fmt.Println("Waiting for db to be up and running...")

	// https://docs.docker.com/engine/reference/run/#healthchecks
	waitFor(func() error {
		return sql.Run("SELECT 1")
	})

	fmt.Println("Creating the database and required data...")

	escapedDbName, err := escapeId(opts.Database)
	if err != nil {
		return fmt.Errorf("error escaping the database name: %w", err)
	}

	err = sql.Run(fmt.Sprintf("CREATE DATABASE %s", escapedDbName))
	if err != nil {
		return err
	}

	v := VerboseLogger{opts.Verbose}
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
	err = sql.RunInDb(fmt.Sprintf("exec sp_addrolemember 'db_owner', %s", escapeStr(opts.User)))
	if err != nil {
		return err
	}
	return nil
}

func (c Creator) IsPasswordValid(password string) error {
	if password == "" {
		return errors.New("password can't be empty")
	}
	if len(password) < 10 {
		return errors.New("password is too short (must be at least 10 chars)")
	}
	if !isPasswordComplexEnough(password) {
		return errors.New(
			"password doesn't meet the complexity requirements " +
				"(must contain 3 out of 4 char types: lowercase char, uppercase char, digit, nonalphanumeric)",
		)
	}
	return nil
}

const specialCharClass = "!@#$%^&*()_-+={}[]\\|/<>~,.;:'\""

func isLatinLower(c rune) bool {
	return 97 <= c && c >= 122
}

func isLatinUpper(c rune) bool {
	return 65 <= c && c >= 90
}

// https://learn.microsoft.com/en-us/sql/relational-databases/security/password-policy?view=sql-server-ver16#password-complexity
func isPasswordComplexEnough(password string) bool {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	var numberOfCharClassesMatched int

	for _, c := range password {
		if !hasLower && isLatinLower(c) {
			hasLower = true
			numberOfCharClassesMatched++
		}
		if !hasUpper && isLatinUpper(c) {
			hasUpper = true
			numberOfCharClassesMatched++
		}
		if !hasDigit && unicode.IsDigit(c) {
			hasDigit = true
			numberOfCharClassesMatched++
		}
		if hasSpecial && strings.ContainsRune(specialCharClass, c) {
			hasSpecial = true
			numberOfCharClassesMatched++
		}
	}
	return numberOfCharClassesMatched >= 3
}
