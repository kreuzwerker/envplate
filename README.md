# Envplate

[![Release](https://img.shields.io/github/v/release/kreuzwerker/envplate)](https://github.com/kreuzwerker/envplate/releases)
[![Build Status](https://github.com/kreuzwerker/envplate/workflows/build/badge.svg)](https://github.com/kreuzwerker/envplate/actions)
[![Documentation](https://godoc.org/github.com/kreuzwerker/envplate?status.svg)](http://godoc.org/github.com/kreuzwerker/envplate) 
[![Go Report Card](https://goreportcard.com/badge/github.com/kreuzwerker/envplate)](https://goreportcard.com/report/github.com/kreuzwerker/envplate) 

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
* charset (`-c` flag): specify a custom charset to replace variables. Example: `/usr/local/bin/ep -c iso8859-1 *.conf`

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
RUN wget -q https://github.com/kreuzwerker/envplate/releases/download/v1.0.2/envplate_1.0.2_$(uname -s)_$(uname -m).tar.gz -O - | tar xz && mv envplate /usr/local/bin/ep && chmod +x /usr/local/bin/ep

...

CMD [ "/usr/local/bin/ep", "-v", "/etc/nginx/nginx.conf", "--", "/usr/sbin/nginx", "-c", "/etc/nginx/nginx.conf" ]
```

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

### Sample docker file

```
FROM nginx:latest
MAINTAINER Albert van t Hart <avthart@gmail.com>

ADD https://github.com/kreuzwerker/envplate/releases/download/v0.0.7/ep-linux /bin/ep
RUN chmod +x /bin/ep

EXPOSE 80 443

CMD [ "/bin/ep", "-v", "/etc/nginx/*.conf", "--", "/usr/sbin/nginx", "-g", "daemon off;" ]
```
Source: https://github.com/avthart/docker-nginx-env