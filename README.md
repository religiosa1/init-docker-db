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

```bash
./init-docker-db
database type? [postgres,mysql,mssql,mongo] (postgres):
> postgres
database name? (db):
> my_awesome_db
database user? (postgres):
> my_user
database password? (123456):
> qwerty
docker container name? (random-shortname):
> mydb

# Outputs ID of the created container on success:
2d42f1fbfcd63c64a56a9034af26f0bffe1157c125f4921a5f6e08c4e22a311c
```

This will create a database container with the specified parameters, exposing
its port (depending on the type, 5432 for postgres, 3306 for MySql,
1433 for MsSql, 27017 for Mongo).

Alternatively, you can configure any of the parameters and the port by the CLI
flags:

```
Positionals:
  containerName  name of the database container to be created           [string]

Options:
  -t, --type             database type
                       [string] [choices: "postgres", "mysql", "mssql", "mongo"]
  -u, --user             database user                                  [string]
  -d, --database         database name                                  [string]
  -p, --password         user's password                                [string]
  -P, --port             TCP port to which database will be mapped to   [number]
  -T, --tag              docker tag to use with the container           [string]
  -n, --non-interactive  exit if any required parameters are missing   [boolean]
  -D, --dry              dry run, printing docker command to stdout, without
                         actually running it                           [boolean]
  -v, --verbose          Run with verbose logging                      [boolean]
  -h, --help             Show help                                     [boolean]
      --version          Show version number                           [boolean]

Examples:
  init-docker-db                       Run in wizard mode
  init-docker-db --dry                 Dry-run in wizard mode
  init-docker-db -t mssql -u app_user  Create a MsSQL database using provided
                                       username
```

For MySQL and MsSQL as there is a separate root user with a predefined name
(root and SA correspondingly), we're using the same password for root access
and user access. It's a _disposable_ database after all.

## Installation

The easiest way to install and use this script is to grab an executable
from the latest release and drop it somewhere in your PATH.

Operation requires docker-engine to be available on your machine, with
the `docker` command available in the PATH, unless you're using the dry mode.

## Running/compiling localy

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

If you have suggsetions or bugfix, feel free to create a PR!

## License

init-docker-db is MIT licensed.
