# init-docker-db changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Changed

- Minor `--help` output changes

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
