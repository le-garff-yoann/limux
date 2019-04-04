# Filemux

## Build

```bash
# GO111MODULE=on go mod vendor # Rebuild the vendors.

CGO_ENABLED=0 GOOS=linux go build -o dist/linux/filemux # Build for Linux.
CGO_ENABLED=0 GOOS=windows go build -o dist/windows/filemux.exe # Build for Windows.
```

## Run tests

```bash
go test ./...
```

## Configuration file

Samples are available [here](samples/etc/filemux).

- The `out` key/value pair represent a list of objects where each object represent a ticker who periodically (configured with `interval`) check for files at the glob `${src}/*`. Whenever a tick is fired it verify if there are files availables for archiving (each filename within the archive can be prepended via `archive_inner_dirname`) at `${dst}/${archive_basename}.tar`. Files successfully written in the archive are removed from the `src`. The `exec` command is finally executed.
- The `in` key/value pair represent a list of objects where each object represent a directory notifier who listen for create-like events at `src`. Whenever a new file is created a routine is started to monitor that the writing of this file is over. The tarball is finally extracted at `dst`.

All Go templates are injected with the [sprig helpers](http://masterminds.github.io/sprig).

## Configuration file validation

```bash
filemux validate -c filemux.yml
```

## Run

```bash
# filemux help # Print the global help.
# filemux help run # Print the help for the run subcommand.

filemux run -c filemux.yml
```

## Run it as a `initd` service

Look [here](samples/etc/init.d/filemux) for the *service unit file*.

```bash
chkconfig filemux on
service filemux start
```

An exemple of logrotate file can be found [here](samples/etc/logrotate.d/filemux).

## Run it as a `systemd` service

Look [here](samples/etc/systemd/system/filemux.service) for the *service unit file*.

```bash
systemctl enable filemux
systemctl start filemux
```

## Run it as a Win32 service 

1. [Download the *Non-Sucking Service Manager*](https://nssm.cc/download).
2.
```cmd
nssm install filemux C:\filemux\filemux.exe run -c C:\filemux\filemux.yml
nssm set filemux AppStdout C:\filemux\log\filemux.log

REM "net stop filemux": 10000ms is the time left to the service to gracefully stop before TerminateProcess() is called.
nssm set filemux AppStopMethodConsole 10000

REM The user filemux will run the service
nssm get filemux ObjectName filemux <password>

net start filemux
```
