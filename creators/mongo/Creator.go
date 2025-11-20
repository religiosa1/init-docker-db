// Package mongo implements DBCreator interface for MongoDB
package mongo

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbcreator"
)

type Creator struct{}

const port uint16 = 27017

func (c Creator) GetDefaultOpts() dbcreator.DefaultOpts {
	return dbcreator.DefaultOpts{
		Port:      port,
		User:      "mongo",
		DockerTag: "latest",
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
	// https://hub.docker.com/_/mongo
	args := []string{
		"run", "--name", opts.ContainerName,
		"-e", dbcreator.DockerEnv("MONGO_INITDB_ROOT_PASSWORD", opts.Password),
		"-e", dbcreator.DockerEnv("MONGO_INITDB_ROOT_USERNAME", opts.User),
		"-e", dbcreator.DockerEnv("MONGO_INITDB_DATABASE", opts.Database),
	}
	args = append(args, dbcreator.CreatePortBindingsArgument(port, opts.Ports)...)
	args = append(args, "-d", fmt.Sprintf("mongo:%s", opts.DockerTag))

	return shell.Run("docker", args...)
}

func (c Creator) ValidatePassword(password string) error {
	return nil
}
