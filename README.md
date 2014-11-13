# Envplate

Trivial templating for configuration files using environment keys:

1. Keys are defined as `${key}` in various configuration files (glob patterns work)
* Definitions without a value result in an error
* envplate can optionally
	* create backups using the `-b` flag, appending a `.bak` extension to backup copies
	* optionally output to stdout instead of replacing values in files using the `-d` flag

Example: `envplate /etc/nginx/*.conf`

## Why?

For apps running Docker which rely (fully or partially) on configuration files instead of being purely configured through environment variables.

You can directly download envplate binaries into your Dockerfile using Github releases like this:

```
RUN curl -sLo /usr/local/bin/ep https://github.com/yawn/envplate/releases/download/v0.0.1/ep-linux && chmod +x /usr/local/bin/ep
```
