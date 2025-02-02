package dbCreator

import "fmt"

type Shell func(argv ...string) error

func DockerEnv(key string, value string) string {
	// TODO
	return fmt.Sprintf("%s=%s", key, value)
}
