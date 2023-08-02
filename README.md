# logparams
[![go.dev Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/aaronvb/logparams) 
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/aaronvb/logparams/go.yml?label=tests)

This is a Go middleware log output that prints parameters if present in the HTTP request. Currently supports `PostForm`, `query params`, and `JSON` body.

The output can be a string or printed directly to the logger. Recommend using with middleware, see example below.

## Install
```sh
go get -u github.com/aaronvb/logparams
```

## Usage
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

Returning data in struct:
```go
lp := logparams.LogParams{Request: r}
lp.ToFields()
```
```go
type ParamFields struct {
	Form      map[string]string
	Query     map[string]string
	Json      map[string]interface{}
	JsonArray []map[string]interface{}
}
```


## Middleware Example (using [gorilla/mux](https://github.com/gorilla/mux))
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaronvb/logparams"

	"github.com/gorilla/mux"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", ":8080")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/foobar", app.foobar).Methods("POST")

	// Middleware
	r.Use(app.logRequest)
	r.Use(app.logParams)

	return r
}

func (app *application) foobar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

// Middleware

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s", r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) logParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lp := logparams.LogParams{Request: r}
		lp.ToLogger(app.infoLog)
		next.ServeHTTP(w, r)
	})
}
```

```sh
> go run main.go
INFO	2020/03/22 11:14:49 Starting server on :8080
INFO	2020/03/22 11:14:51 POST - /foobar
INFO	2020/03/22 11:15:18 Parameters: {"foo" => "bar"}
```

## Optional Values
- `ShowEmpty (bool)` will return an empty string, or not print to logger, if there are no parameters. Default is to false if struct arg is not passed.

- `ShowPassword (bool)` will show the `password` and `password_confirmation` parameters. Default is false if not explicitly passed(DO NOT RECOMMEND).

- `HidePrefix (bool)` will hide the `Parameters: ` prefix in the output. Default is to false if struct arg is not passed.
