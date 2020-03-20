# aaronvb/logparams
Package `aaronvb/logparams` implements an parameter output if present in the HTTP request.

The output can be a string or printed directly to the logger.

## Install
```sh
go get -u github.com/aaronvb/logparams
```

## Example
Using logger:
```go
var logger log.Logger{}
lp := logparams.LogParams{Request: r, HideEmpty: true}
lp.ToLogger(&logger)
```

Output to a string:
```go
lp := logparams.LogParams{Request: r, HideEmpty: true}
lp.ToString()
```

Result:
```sh
Parameters: {"foo" => "bar", "hello" => "world"}
```

## Optional Values
`HideEmpty (bool)` will return an empty string or not print to logger if there are no parameters.

