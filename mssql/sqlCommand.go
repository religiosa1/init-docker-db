package mssql

import (
	"fmt"
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

// TODO: parse the resposnse os SQL server and check for potential errors

func (r SqlInContainerRunner) Run(sql string) error {
	if r.verbose {
		if strings.ContainsRune(sql, '\n') {
			fmt.Printf("SQL:\n%s --> END SQL \n", sql)
		} else {
			fmt.Printf("SQL: %s\n", sql)
		}
	}
	err := r.shell.RunSilent("docker", "exec", "-it", r.contId,
		"/opt/mssql-tools/bin/sqlcmd", "-S", "localhost",
		"-U", "SA", "-P", r.password, "-Q", sql, "||", "exit 1")
	return err
}

func (r SqlInContainerRunner) RunInDb(sql string) error {
	escapedDbName, err := escapeId(r.database)
	if err != nil {
		return err
	}
	return r.Run(fmt.Sprintf("use %s\n%s", escapedDbName, sql))
}
