# init-docker-db changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## [1.0.0] - 2025.04.02

### Changed

- Full rewrite to golang.
- Postgres default password is now "password" and creation of DB without password
  is forbidden, as it results in failed container.
- `Done` suffix removed from thr program output, only direct `docker run`
  output with container's ID is now printed in non-verbose mode.

## [0.3.0] - 2025.02.11

### Added

- MsSQL handling of potential SQL errors during container's setup.

### Fixed

- max delay clamping in exponential back-off waiting for MsSQL DB to be up and running prior to running SQL commands
- MySQL non-root username fix (previously matched to password)

### Changed

- MsSQL health check switched from `SELECT 1` to `SELECT SERVERPROPERTY('ProductVersion')`
- MsSQL add owner switched from deprecated `sp_addrolemember` to `ALTER ROLE db_owner`
- Minor `--help` and error messages output changes

## [0.2.2] - 2025.01.31

### Fixed

- MsSQL regression introduced by the dry run functionality

## [0.2.1] - 2025.01.30 [YANKED]

### Added

- CLI version flag support
- Dry-run functionality

### Changed

- Bun version bumped to 1.2.1 with lockfile changed to text version

## [0.2.0] - 2024.10.03

### Added

- Changelog itself
- MongoDB support
- Password validation in non-interactive mode for MsSQL
- MsSQL password complexity now allows nonalphnumeric symbols as per docs

### Fixed

- docker tag parameter being ignored for all of the db types
- MsSQL initial db creation with non-default password because of failed SA login
- Ignored container name in non-interactive mode

## [0.1.0] - 2024.06.22

### Added

- MsSQL support
- Docker tag parameter

### Changed

- verbose logging displaying more information

## [0.0.1] - 2024.06.21

First public release, postgres and mysql support.
