<p align="center">
  <img alt="piglatin logo" src="assets/logo.png" height="150"/>
  <h3 align="center">Detask</h3>
  <p align="center">Task management tool</p>
</p>

---

`detask` is a simple server-side application that provides API for managing tasks.<br />
The app is available online on Heroku: [go-detask.herokuapp.com](https://go-detask.herokuapp.com/api/v1/health) <br />

## Badges
[![Circleci](https://circleci.com/gh/dnozdrin/detask.svg?style=shield)](https://circleci.com/gh/dnozdrin/detask)
[![Codecov](https://codecov.io/gh/dnozdrin/detask/branch/master/graph/badge.svg)](https://codecov.io/gh/dnozdrin/detask)
[![Coreportcard](https://goreportcard.com/badge/github.com/dnozdrin/detask)](https://goreportcard.com/report/github.com/dnozdrin/detask)
[![License](https://img.shields.io/github/license/dnozdrin/detask)](/LICENSE)
[![Release](https://img.shields.io/github/release/dnozdrin/detask.svg)](https://github.com/dnozdrin/detask/releases/latest)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

To build, install and run the app you will need:

- [Go 1.14](https://golang.org/dl)
- [PostgreSQL 12](https://www.postgresql.org/download/)

As alternative, you may use docker containers.
To run tests or serve REST API docs locally, you will need [Docker Compose](https://docs.docker.com/compose) and [Make](https://en.wikipedia.org/wiki/Make_(software)).

### Installing

- Clone the project.
- Compile the app by running the next command in the project directory:

```shell script
make build
```

The compiled binary will be available as `./bin/detask`

- To run the app you will need to set the next environment variables:

| Variable | Description | Example |
|:--------|-------------|---------|
| DB_HOST | database host | `localhost` |
| DB_PORT | database port | `5432` |
| DB_NAME | database name | `postgres` |
| DB_USER | database user | `postgres` |
| DB_PASS |  database password | `superMegaPass123#!` |
| DB_MIGRATION_PATH | path to the sql migrations | `file:///app/internal/db/migrations` |
| PORT | port where the app server will work | `80` |
| APP_ALLOWED_ORIGINS | allowed origins for 'Access-Control-Allow-Origin' header, separated with comma | `http://localhost:8081,http://localhost:80` |
| APP_CONTEXT | application context | `development` |
| APP_LOG_PATH | path where app log will be stored | `stderr` |

Supported application contexts:

| Context | Description |
|:--------|-------------|
| production | for production usage |
| testing | for running automatic tests |
| development | should be used during development |

## REST API
REST API documentation is available on [dnozdrin.github.io/detask](https://dnozdrin.github.io/detask)

To start a local instance of the Swagger UI, run the next command in the project directory:

```shell script
make swagger-start
```

This will make the documentation available on `http://localhost:8081/`

To stop the Swagger UI container, run:

```shell script
make swagger-stop
```

## Running the tests

The project contains with unit and end to end tests.

### Unit tests

Unit tests are placed in the same directory as the target code. To start unit tests run in the project directory:

```shell script
make unit-test
```

### End to end tests

End to end tests are placed in the `./test` directory. To start these tests run in the project directory:

```shell script
make integration-test
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/dnozdrin/detask/tags). 

## Authors

* **Dmytro Nozdrin** - *Initial work* - [dnozdrin](https://github.com/dnozdrin)
