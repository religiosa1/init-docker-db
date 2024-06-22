# init-docker-db

Simple script to initialize a disposable docker container with a database.

During the development process I want to create database containers for my apps.
And I can never remember the environment variables and sql commands to initialize them.

This script automates this process, creating a database of specified type, user, password
database name and container name.

Written using [bun](https://bun.sh/)

## Usage

The easiest way is to launch the script in wizard mode:

```bash
./init-docker-db
database type? [postgres,mysql] (postgres):
> your database type here
database name? (db):
> your database name here
database user? (postgres):
> your user here
database password? (123456):
> your password here
docker container name? (apathetic-devotion):
> container name
```

It will create a database container with set parameters, exposing its port (depending
on the type, 5432 for postgres, 3306 for MySql).

Alternatively, you can configure any of the parameters and the port by the CLI
flags:

```
Positionals:
  name  name of the database container to be created                    [string]

Options:
      --version          Show version number                           [boolean]
  -t, --type             database type   [string] [choices: "postgres", "mysql"]
  -u, --user             database user                                  [string]
  -d, --database         database name                                  [string]
  -p, --password         user's password                                [string]
  -P, --port             TCP port to which database will be mapped to   [number]
  -T, --tag              docker tag to use with the container           [string]
  -n, --non-interactive  exit, if some of the required params are missing
                                                                       [boolean]
  -v, --verbose          Run with verbose logging                      [boolean]
  -h, --help             Show help                                     [boolean]
```

For MySQL as there is a separate root user with a predefined name
(root and SA correspondingly), we're using the same password for root access
and user access. It's a _disposable_ database after all.

## Installation

The easiest way to install and use this script is to grab an executable
from the latest release and drop it somewhere in your PATH.

## Compilation

To install dependencies:

```bash
bun install
```

To run:

```bash
bun run index.ts
```

To build a standalone executable:

```bash
bun run compile
```

To cross-compile all executables for all available targets:

```base
bun run compile-all
```

For additional compilation targets, please refer to the [bun docs](https://bun.sh/docs/bundler/executables#cross-compile-to-other-platforms)

## Contribution

If you have suggsetions or bugfix, feel free to create a PR!

## License

init-docker-db is MIT licensed.
