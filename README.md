# Envplate

[![Build Status](https://travis-ci.org/yawn/envplate.svg)](https://travis-ci.org/yawn/envplate)

Trivial templating for configuration files using environment keys:

1. Keys are defined as `${key}` in various configuration files (glob patterns work)
* Definitions without a value result in an error
* envplate can optionally
	* create backups using the `-b` flag, appending a `.bak` extension to backup copies
	* output to stdout instead of replacing values in files using the `-d` flag
	* be verbose about it's operations by using the `-v` flag
  * `exec()` another command by passing `--` after all arguments to `ep`, the path to the binary and arguments to it, e.g. `/usr/local/ep -v *.conf -- /usr/sbin/nginx -c /etc/nginx.conf`; this can be used to use envplate to parse configs and execute the container process using Dockers `CMD`

Example:

```
$ cat /etc/foo.conf
Database=${FOO_DATABASE}
Mode=fancy

$ export FOO_DATABASE=db.example.com

$ ep /etc/f*.conf

$ cat /etc/foo.conf
Database=db.example.com
Mode=fancy
```

## Why?

For apps running Docker which rely (fully or partially) on configuration files instead of being purely configured through environment variables.

You can directly download envplate binaries into your Dockerfile using Github releases like this:

```
RUN curl -sLo /usr/local/bin/ep https://github.com/yawn/envplate/releases/download/v0.0.2/ep-linux && chmod +x /usr/local/bin/ep

...

CMD /usr/local/bin/ep -v /etc/nginx/nginx.conf -- /usr/sbin/nginx -c /etc/nginx/conf
```
