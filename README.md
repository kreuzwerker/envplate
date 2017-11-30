# Envplate

[![Build Status](https://travis-ci.org/kreuzwerker/envplate.svg)](https://travis-ci.org/kreuzwerker/envplate)

Trivial templating for configuration files using environment keys. References to such keys are declared in arbitrary config files either as:

1. `${key}` or
* `${key:-default value}`

Envplate (`ep`) parses arbitrary configuration files (using glob patterns) and replaces all references with values from the environment or with default values (if given). These values replace the keys *inline* (= the files will be changed).

Failure to resolve the supplied glob pattern(s) to at least one file results in an error.

Optionally, `ep` can:

* backup (`-b` flag): create backups of the files it changes, appending a `.bak` extension to backup copies
* dry-run (`-d` flag): output to stdout instead of replacing values inline
* strict (`-s` flag): refuse to fallback to default values
* verbose (`-v` flag): be verbose about it's operations

`ep` can also `exec()` another command by passing

* `--` after all arguments to `ep`
* the path to the binary and it's arguments

Example: `/usr/local/bin/ep -v *.conf -- /usr/sbin/nginx -c /etc/nginx.conf`

This can be used to use `ep` to parse configs and execute the container process using Dockers `CMD`

## Escaping

In case the file you want to modify already uses the pattern envplate is searching for ( e.g. for reading environment variables ) you can escape the sequence by adding a leading backslash `\`. It's also possible to escape a leading backslash by adding an additional backslash. Basically a sequence with an even number of leading backslashes will be parsed, is the number of leading backslashes odd the sequence will be escaped.

See https://github.com/kreuzwerker/envplate#full-example

## Why?

For apps running Docker which rely (fully or partially) on configuration files instead of being purely configured through environment variables.

You can directly download envplate binaries into your Dockerfile using Github releases like this:

```
RUN curl -sLo /usr/local/bin/ep https://github.com/kreuzwerker/envplate/releases/download/v0.0.8/ep-linux && chmod +x /usr/local/bin/ep

...

CMD [ "/usr/local/bin/ep", "-v", "/etc/nginx/nginx.conf", "--", "/usr/sbin/nginx", "-c", "/etc/nginx/nginx.conf" ]
```
Have a look at https://github.com/avthart/docker-nginx-env/blob/master/Dockerfile to see a working example Dockerfile.

## Full example

```
$ cat /etc/foo.conf
Database=${FOO_DATABASE}
DatabaseSlave=${BAR_DATABASE:-db2.example.com}
Mode=fancy
Escaped1=\${FOO_DATABASE}
NotEscaped1=\\${FOO_DATABASE}
Escaped2=\\\${BAR_DATABASE:-db2.example.com}
NotEscaped2=\\\\${BAR_DATABASE:-db2.example.com}

$ export FOO_DATABASE=db.example.com

$ ep /etc/f*.conf

$ cat /etc/foo.conf
Database=db.example.com
DatabaseSlave=db2.example.com
Mode=fancy
Escaped1=${FOO_DATABASE}
NotEscaped1=\db.example.com
Escaped2=\${BAR_DATABASE:-db2.example.com}
NotEscaped2=\\db2.example.com
```
