# aaronvb/logparams
Package `aaronvb/logparams` implements an parameter output if present in the HTTP request. Currently supports `PostForm`, `query params`, and `JSON` body.

The output can be a string or printed directly to the logger.

## Install
```sh
go get -u github.com/aaronvb/logparams
```

## Example
Using logger:
```go
var logger log.Logger{}
lp := logparams.LogParams{Request: r}
lp.ToLogger(&logger)
```

Output to a string:
```go
lp := logparams.LogParams{Request: r}
lp.ToString()
```

Result:
```sh
Parameters: {"foo" => "bar", "hello" => "world"}
```

## Optional Values
`ShowEmpty (bool)` will return an empty string, or not print to logger, if there are no parameters. Default is to false if struct arg is not passed.
`ShowPassword (bool)` will show the `password` and `password_confirmation` parameters. Default is false if not explicitly passed(DO NOT RECOMMEND).
