package logparams

import (
	"fmt"
	"log"
	"net/http"
)

// LogParams struct
type LogParams struct {
	Request   *http.Request
	HideEmpty bool
}

// ToString will return a string of all parameters within the http request.
func (lp LogParams) ToString() string {
	if lp.HideEmpty == true && lp.parseParams() == "" {
		return ""
	}

	return fmt.Sprintf("Parameters: {%s}", lp.parseParams())
}

// ToLogger will log print all parameters within the http request.
func (lp LogParams) ToLogger(logger *log.Logger) {
	if lp.HideEmpty == true && lp.parseParams() == "" {
		return
	}

	logger.Printf("Parameters: {%s}", lp.parseParams())
}

// parseParams will parse the form for values and return a string of parameters
func (lp LogParams) parseParams() string {
	var paramString string
	err := lp.Request.ParseForm()
	if err != nil {
		return paramString
	}

	if len(lp.Request.PostForm) > 0 {
		var paramCount = 0
		for k, _ := range lp.Request.PostForm {
			paramCount += 1
			paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, lp.Request.PostForm.Get(k))
			if paramCount != len(lp.Request.PostForm) {
				paramString += ", "
			}
		}
	}
	return paramString
}

// parseQueryParams will parse query parameters in the URL.
func (lp LogParams) parseQueryParams() string {
	var paramString string

	var paramCount = 0
	for k := range lp.Request.URL.Query() {
		paramCount += 1
		paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, lp.Request.URL.Query()[k][0])
		if paramCount != len(lp.Request.URL.Query()) {
			paramString += ", "
		}
	}

	return paramString
}

