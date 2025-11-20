// Package mysql implements DBCreator interface for MYSQL
package mysql

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbcreator"
)

type Creator struct{}

const port uint16 = 3306

func (c Creator) GetDefaultOpts() dbcreator.DefaultOpts {
	return dbcreator.DefaultOpts{
		Port:      port,
		User:      "mysql",
		DockerTag: "lts",
		Password:  "",
	}
}

func (c Creator) GetCapabilities() dbcreator.Capabilities {
	return dbcreator.Capabilities{
		DatabaseName: true,
		UserPassword: true,
	}
}

func (c Creator) Create(shell dbcreator.Shell, opts dbcreator.CreateOptions) error {
	// https://hub.docker.com/_/mysql
	args := []string{
		"run", "--name", opts.ContainerName,
		"-e", dbcreator.DockerEnv("MYSQL_USER", opts.User),
		"-e", dbcreator.DockerEnv("MYSQL_ROOT_PASSWORD", opts.Password),
		"-e", dbcreator.DockerEnv("MYSQL_PASSWORD", opts.Password),
		"-e", dbcreator.DockerEnv("MYSQL_DATABASE", opts.Database),
	}
	args = append(args, dbcreator.CreatePortBindingsArgument(port, opts.Ports)...)
	args = append(args, "-d", fmt.Sprintf("mysql:%s", opts.DockerTag))
	return shell.Run("docker", args...)
}

func (c Creator) ValidatePassword(password string) error {
	return nil
}
