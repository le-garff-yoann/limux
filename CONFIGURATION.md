# Limux configuration

## Configuration file

```yaml
out:
  - interval: 1s
    src: /tmp/data/A/*/out
    dst: /tmp/spool/A/out
    archive_basename: A_{{ .Now.UnixNano }}_0
    archive_inner_dirname: "{{ index .SrcFragments 4 }}/in"
    exec:
      - bash
      - -c
      - mv {{ .ArchiveFullPath }} /tmp/spool/A/in

in:
  - src: /tmp/spool/A/in
    dst: /tmp/data/A
```

- The `out` key/value pair represent a list of objects where each object represent a ticker who periodically (configured with `interval`) check for files at the glob `${src}/*`. Whenever a tick is fired it verify if there are files availables for archiving (each filename within the archive can be prepended via `archive_inner_dirname`) at `${dst}/${archive_basename}.tar`. Files successfully written in the archive are removed from the `src`. The `exec` command is finally executed.
- The `in` key/value pair represent a list of objects where each object represent a directory notifier who listen for create-like events at `src`. Whenever a new file is created a routine is started to monitor that the writing of this file is over. The tarball is finally extracted at `dst`.

All Go templates are injected with the [sprig helpers](http://masterminds.github.io/sprig).

## Configuration file validation

```bash
limux validate -c config.yml
```
