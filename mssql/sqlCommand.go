package mssql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type SqlInContainerRunner struct {
	shell    *dbCreator.Shell
	contId   string
	database string
	password string
	verbose  bool
}

func (r SqlInContainerRunner) Run(sql string) error {

	if r.verbose {
		if strings.ContainsRune(sql, '\n') {
			fmt.Printf("SQL:\n%s --> END SQL \n", sql)
		} else {
			fmt.Printf("SQL: %s\n", sql)
		}
	}
	out, err := r.shell.RunWithOutput(
		"docker", "exec", r.contId,
		"/opt/mssql-tools/bin/sqlcmd", "-S", "localhost",
		"-U", "SA", "-P", r.password, "-Q", sql,
	)
	if r.verbose {
		fmt.Println(out)
	}
	if err != nil {
		return err
	}
	err = hasSqlCommandError(out)
	if err != nil && !r.verbose {
		fmt.Println(out)
	}
	return err
}

func (r SqlInContainerRunner) RunInDb(sql string) error {
	escapedDbName, err := escapeId(r.database)
	if err != nil {
		return err
	}
	return r.Run(fmt.Sprintf("use %s\n%s", escapedDbName, sql))
}

var mssqlErrRe *regexp.Regexp

func hasSqlCommandError(output string) error {
	if mssqlErrRe == nil {
		mssqlErrRe = regexp.MustCompile(`^Msg (?:\d+), Level (\d+), State (?:\d+), Server (?:[^,]+), Line (?:\d+)`)
	}

	matches := mssqlErrRe.FindStringSubmatch(output)
	if matches == nil {
		return nil
	}

	var errorLevel int
	fmt.Sscanf(matches[1], "%d", &errorLevel)

	// Any severity level less or equal 10 we treat as not an error
	// https://learn.microsoft.com/en-us/sql/relational-databases/errors-events/database-engine-error-severities?view=sql-server-ver16#levels-of-severity
	if errorLevel > 10 {
		return errors.New("sql error")
	}
	return nil
}
