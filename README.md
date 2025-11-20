# init-docker-db

A simple single executable program/script to initialize a disposable docker
container with a database for the development process.

As I can never remember all the environment variables and sql commands (looking
at you, MsSQL) required to initialize a docker container with a db in it.

`init-docker-db` automates this process, creating a database of specified type,
user, password, database name and container name, ensuring consistent parameters
across different database types.

Using `--dry` flag, you can print the required commands to the terminal without
actually executing them.

## Usage

The easiest way is to launch the script in wizard mode:

<https://github.com/user-attachments/assets/e02a4770-91de-43c4-862e-6b3f2b760f77>

This will create a database container with the specified parameters, exposing
its port (depending on the type, 5432 for postgres, 3306 for MySql,
1433 for MsSql, 27017 for Mongo, 6379 for redis) only on localhost (both on
IPv4 and IPv6 interfaces -- depending on the availability). If `--public` flag
is supplied, then port will be exposed on 0.0.0.0 interface available from the
outside world.

Alternatively, you can configure any of the parameters and the port by the CLI
flags:

```
usage: init-docker-db [<containername>] [flags]

Create a disposable database docker container.

Arguments:
  [<containerName>]    name of the database container to be created

Flags:
  -h, --help               Show context-sensitive help.
  -t, --type=STRING        database type
  -u, --user=STRING        database user
  -d, --database=STRING    database name
  -P, --password=STRING    user's password
  -p, --port=PORT,...      port with optional IP address to which database will be mapped to
      --public             expose default port to outside world by mapping to 0.0.0.0 IP address
  -T, --tag=STRING         docker tag to use with the container
  -n, --non-interactive    exit if any required parameters are missing
  -D, --dry                dry run, printing docker command to stdout, without actually running it
  -v, --verbose            run with verbose logging
      --version            show version and exit
  -h, --help               show help message and exit

Examples:
  init-docker-db                       Run in wizard mode
  init-docker-db --dry                 Dry-run in wizard mode
  init-docker-db -t mssql -u app_user  Create a MsSQL database using provided username
```

For MySQL and MsSQL as there is a separate root user with a predefined name
(root and SA correspondingly), we're using the same password for root access
and user access. It's a _disposable_ database after all.

## Installation

The easiest way to install and use this script is to grab an executable
from the latest release and drop it somewhere in your PATH.

Operation requires docker-engine to be available on your machine, with
the `docker` command available in the PATH, unless you're using the dry mode.

Alternatively, if you have `go` installed on your machine you can run:

```bash
go install github.com/religiosa1/init-docker-db@latest
```

## Running/compiling locally

To compile and run project, you need [go](https://go.dev/) version 1.22 or
higher.

To run the application locally from the source:

```bash
go run
```

To build a standalone executable in the `./bin` folder:

```bash
go build
```

To launch unit-tests:

```bash
go test ./...
```

Other tasks and targets are defined in the [taskfile](https://taskfile.dev/).

To cross-compile all executables for all available targets:

```bash
task all
```

## Contribution

If you have suggestions or bugfix, feel free to create a PR!

## License

init-docker-db is MIT licensed.
