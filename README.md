# Limux

[![pipeline status](https://gitlab.com/le-garff-yoann/limux/badges/master/pipeline.svg)](https://gitlab.com/le-garff-yoann/limux/pipelines)

## Build

```bash
# GO111MODULE=on go mod vendor # Rebuild the vendors.

CGO_ENABLED=0 GOOS=linux go build -o dist/linux/limux # Build for Linux.
CGO_ENABLED=0 GOOS=windows go build -o dist/windows/limux.exe # Build for Windows.
```

## Run tests

```bash
go test ./...
```

## Configuration

Take a look [here](CONFIGURATION.md).

## Run

```bash
# limux help # Print the global help.
# limux help run # Print the help for the run subcommand.

limux run -c config.yml
```

## Frontend

Take a look [here](vue/limux/).
