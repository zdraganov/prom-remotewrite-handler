## Prometheus remote_write handler ##

The example_remote_writer is a Go application, that parses the
Prometheus protobuf request and converts the timeseries to JSON. For
this example the response is just printed in the STDOUT.

Run the dev application using:
```
$ docker-compose up
```

If you need to rebuild, just pass the `--rebuild` flag to the
docker-compose.

:warning: &nbsp; NaN and +/-Inf are converted to nil, because JSON does not supports those types
