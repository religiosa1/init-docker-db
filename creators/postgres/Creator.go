// Package postgres implements DBCreator interface for PostreSQL
package postgres

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbcreator"
)

type Creator struct{}

const port uint16 = 5432

func (c Creator) GetDefaultOpts() dbcreator.DefaultOpts {
	return dbcreator.DefaultOpts{
		Port:      port,
		User:      "postgres",
		DockerTag: "latest",
		Password:  "postgres",
	}
}

func (c Creator) GetCapabilities() dbcreator.Capabilities {
	return dbcreator.Capabilities{
		DatabaseName: true,
		UserPassword: true,
	}
}

func (c Creator) Create(shell dbcreator.Shell, opts dbcreator.CreateOptions) error {
	// https://hub.docker.com/_/postgres
	return shell.Run("docker", "run", "--name", opts.ContainerName,
		"-e", dbcreator.DockerEnv("POSTGRES_PASSWORD", opts.Password),
		"-e", dbcreator.DockerEnv("POSTGRES_USER", opts.User),
		"-e", dbcreator.DockerEnv("POSTGRES_DB", opts.Database),
		"-p", fmt.Sprintf("%s:%d", opts.Port, port),
		"-d", fmt.Sprintf("postgres:%s", opts.DockerTag),
	)
}

func (c Creator) ValidatePassword(password string) error {
	return nil
}
