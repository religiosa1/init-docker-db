package postgres

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type Creator struct{}

const port uint16 = 5432

func (pgs Creator) GetDefaultOpts() dbCreator.DefaultOpts {
	return dbCreator.DefaultOpts{
		Port:      port,
		User:      "postgres",
		DockerTag: "latest",
		Password:  "",
	}
}

func (pgs Creator) Create(shell dbCreator.Shell, opts dbCreator.CreateOptions) error {
	// https://hub.docker.com/_/postgres
	return shell.Run("docker", "run", "--name", opts.ContainerName,
		"-e", dbCreator.DockerEnv("POSTGRES_PASSWORD", opts.Password),
		"-e", dbCreator.DockerEnv("POSTGRES_USER", opts.User),
		"-e", dbCreator.DockerEnv("POSTGRES_DB", opts.Database),
		"-p", fmt.Sprintf("%d:%d", port, opts.Port),
		"-d", fmt.Sprintf("postgres:%s", opts.Tag),
	)
}

func (pgs Creator) ValidatePassword(password string) error {
	return nil
}
