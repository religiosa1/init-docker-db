package main

import (
	"fmt"

	"github.com/religiosa1/init-docker-db/dbCreator"
)

type MySqlCreator struct{}

const port uint16 = 3306

func (pgs MySqlCreator) GetDefaultSettings() dbCreator.DefaultOpts {
	return dbCreator.DefaultOpts{
		Port:      port,
		User:      "mysql",
		DockerTag: "lts",
		Password:  "",
	}
}

func (pgs MySqlCreator) Create(shell dbCreator.Shell, opts dbCreator.CreateOptions) error {
	// https://hub.docker.com/_/mysql
	return shell("docker", "run", "--name", opts.ContainerName,
		"-e", dbCreator.DockerEnv("MYSQL_USER", opts.User),
		"-e", dbCreator.DockerEnv("MYSQL_ROOT_PASSWORD", opts.Password),
		"-e", dbCreator.DockerEnv("MYSQL_PASSWORD", opts.Password),
		"-e", dbCreator.DockerEnv("MYSQL_DATABASE", opts.Database),
		"-p", fmt.Sprintf("%d:%d", port, opts.Port),
		"-d", fmt.Sprintf("mysql:%s", opts.Tag),
	)
}

func (pgs MySqlCreator) IsPasswordValid(password string) error {
	return nil
}
