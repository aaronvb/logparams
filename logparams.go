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

// Helper methods

// checkForFormParams checks for form params in the request.
func (lp LogParams) checkForFormParams() bool {
	err := lp.Request.ParseForm()
	if err != nil {
		return false
	}

	if len(lp.Request.PostForm) == 0 {
		return false
	}

	return true
}

// checkForQueryParams checks for query params in the request.
func (lp LogParams) checkForQueryParams() bool {
	if len(lp.Request.URL.Query()) == 0 {
		return false
	}

	return true
}

// checkForJSON checks for content-type application/json in the header.
func (lp LogParams) checkForJSON() bool {
	if lp.Request.Header.Get("Content-Type") != "application/json" {
		return false
	}

	return true
}

// parseParams will check the type of param in the request and call the correct parser.
func (lp LogParams) parseParams() string {
	if lp.checkForFormParams() == true {
		return fmt.Sprintf("{%s}", lp.parseFormParams())
	} else if lp.checkForQueryParams() == true {
		return fmt.Sprintf("{%s}", lp.parseQueryParams())
	} else if lp.checkForJSON() == true {
		return lp.parseJSONBody()
	}

	return ""
}

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

