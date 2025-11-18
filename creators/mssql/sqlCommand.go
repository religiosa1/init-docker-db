package mssql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/religiosa1/init-docker-db/dbcreator"
)

type SQLInContainerRunner struct {
	shell    *dbcreator.Shell
	contID   string
	database string
	password string
	verbose  bool
}

func (r SQLInContainerRunner) RunSilent(sql string) error {
	_, err := r.run(sql)
	return err
}

func (r SQLInContainerRunner) Run(sql string) error {
	out, err := r.run(sql)
	if err != nil && !r.verbose {
		fmt.Println(out)
	}
	return err
}

func (r SQLInContainerRunner) run(sql string) (string, error) {
	if r.verbose {
		if strings.ContainsRune(sql, '\n') {
			fmt.Printf("SQL:\n%s --> END SQL \n", sql)
		} else {
			fmt.Printf("SQL: %s\n", sql)
		}
	}
	// See https://github.com/microsoft/mssql-docker/issues/892
	// Previous versions used mssql-tools, now it's mssql-tools18
	out, err := r.shell.RunWithOutput(
		"docker", "exec", r.contID,
		"/opt/mssql-tools18/bin/sqlcmd", "-C", "-S", "localhost",
		"-U", "SA", "-P", r.password, "-Q", sql,
	)
	if r.verbose {
		fmt.Println(out)
	}
	if err != nil {
		return out, err
	}
	return out, parseSQLCommandError(out)
}

func (r SQLInContainerRunner) RunInDB(sql string) error {
	escapedDBName, err := escapeID(r.database)
	if err != nil {
		return err
	}
	return r.Run(fmt.Sprintf("use %s\n%s", escapedDBName, sql))
}

var mssqlErrRe = regexp.MustCompile(`^Msg (?:\d+), Level (\d+), State (?:\d+), Server (?:[^,]+), Line (?:\d+)`)

func parseSQLCommandError(output string) error {
	matches := mssqlErrRe.FindStringSubmatch(output)
	if matches == nil {
		return nil
	}

	var errorLevel int
	_, err := fmt.Sscanf(matches[1], "%d", &errorLevel)
	if err != nil {
		return err
	}

	// Any severity level less or equal 10 we treat as not an error
	// https://learn.microsoft.com/en-us/sql/relational-databases/errors-events/database-engine-error-severities?view=sql-server-ver16#levels-of-severity
	if errorLevel > 10 {
		return fmt.Errorf("sql error: %s", output)
	}
	return nil
}
