// Package redis implements DBCreator interface for Redis
package redis

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbcreator"
)

type Creator struct{}

const port uint16 = 6379

func (c Creator) GetDefaultOpts() dbcreator.DefaultOpts {
	return dbcreator.DefaultOpts{
		Port:      port,
		User:      "",
		DockerTag: "latest",
		Password:  "",
	}
}

func (c Creator) GetCapabilities() dbcreator.Capabilities {
	return dbcreator.Capabilities{
		DatabaseName: false,
		UserPassword: false,
	}
}

func (c Creator) Create(shell dbcreator.Shell, opts dbcreator.CreateOptions) error {
	// https://hub.docker.com/_/redis/
	return shell.Run("docker", "run", "--name", opts.ContainerName,
		"-p", fmt.Sprintf("%s:%d", opts.Port, port),
		"-d", fmt.Sprintf("redis:%s", opts.DockerTag),
		"redis-server", "--save", "60", "1", "--loglevel", "warning",
	)
}

func (c Creator) ValidatePassword(password string) error {
	return nil
}
