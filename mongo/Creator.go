package mongo

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type Creator struct{}

const port uint16 = 27017

func (pgs Creator) GetDefaultOpts() dbCreator.DefaultOpts {
	return dbCreator.DefaultOpts{
		Port:      port,
		User:      "mongo",
		DockerTag: "latest",
		Password:  "",
	}
}

func (pgs Creator) Create(shell dbCreator.Shell, opts dbCreator.CreateOptions) error {
	// https://hub.docker.com/_/mongo
	return shell("docker", "run", "--name", opts.ContainerName,
		"-e", dbCreator.DockerEnv("MONGO_INITDB_ROOT_PASSWORD", opts.Password),
		"-e", dbCreator.DockerEnv("MONGO_INITDB_ROOT_USERNAME", opts.User),
		"-e", dbCreator.DockerEnv("MONGO_INITDB_DATABASE", opts.Database),
		"-p", fmt.Sprintf("%d:%d", port, opts.Port),
		"-d", fmt.Sprintf("mongo:%s", opts.Tag),
	)
}

func (pgs Creator) IsPasswordValid(password string) error {
	return nil
}
