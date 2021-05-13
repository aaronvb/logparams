package logparams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// LogParams struct
// Request is the http request
// HideEmpty will not log or return "" if param is empty.
// FilterPassword will filter password parameters (default true).
type LogParams struct {
	Request      *http.Request
	ShowEmpty    bool
	ShowPassword bool
	HidePrefix   bool
}

// ToString will return a string of all parameters within the http request.
func (lp LogParams) ToString() string {
	var str string

	if !lp.ShowEmpty && lp.parseParams() == "" {
		return ""
	}

	if lp.HidePrefix {
		str = lp.parseParams()
	} else {
		str = fmt.Sprintf("Parameters: %s", lp.parseParams())
	}

	return str
}

// ToLogger will log print all parameters within the http request.
func (lp LogParams) ToLogger(logger *log.Logger) {
	var str string

	if !lp.ShowEmpty && lp.parseParams() == "" {
		return
	}

	if lp.HidePrefix {
		str = lp.parseParams()
	} else {
		str = fmt.Sprintf("Parameters: %s", lp.parseParams())
	}

	logger.Printf(str)
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
	return len(lp.Request.URL.Query()) != 0
}

// checkForJSON checks for content-type application/json in the header.
func (lp LogParams) checkForJSON() bool {
	matched, _ := regexp.MatchString(`application\/json`, lp.Request.Header.Get("Content-Type"))
	return matched
}

// parseParams will check the type of param in the request and call the correct parser.
func (lp LogParams) parseParams() string {
	if lp.checkForFormParams() {
		return fmt.Sprintf("{%s}", lp.parseFormParams())
	} else if lp.checkForQueryParams() {
		return fmt.Sprintf("{%s}", lp.parseQueryParams())
	} else if lp.checkForJSON() {
		return lp.parseJSONBody()
	}

	return ""
}

// parseFormParams will parse the form for values and return a string of parameters
func (lp LogParams) parseFormParams() string {
	var paramString string

	err := lp.Request.ParseForm()
	if err != nil {
		return paramString
	}

	var paramCount = 0
	for k := range lp.Request.PostForm {
		if k == "password" || k == "password_confirmation" {
			paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, "[FILTERED]")
		} else {
			paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, lp.Request.PostForm.Get(k))
		}
		paramCount++
		if paramCount != len(lp.Request.PostForm) {
			paramString += ", "
		}
	}

	return paramString
}

// parseQueryParams will parse query parameters in the URL.
func (lp LogParams) parseQueryParams() string {
	var paramString string

	var paramCount = 0
	for k := range lp.Request.URL.Query() {
		paramCount++
		paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, lp.Request.URL.Query()[k][0])
		if paramCount != len(lp.Request.URL.Query()) {
			paramString += ", "
		}
	}

	return paramString
}

// parseJSONBody will parse the json in the body as parameters.
func (lp LogParams) parseJSONBody() string {
	var b []byte
	var result map[string]interface{}
	var resultArray []map[string]interface{}

	body, _ := ioutil.ReadAll(lp.Request.Body)
	lp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err := json.Unmarshal(body, &result)
	if err != nil {
		err := json.Unmarshal(body, &resultArray)
		if err != nil {
			return ""
		}
	}

	if len(result) != 0 {
		if !lp.ShowPassword {
			if result["password"] != nil {
				result["password"] = "[FILTERED]"
			}
			if result["password_confirmation"] != nil {
				result["password_confirmation"] = "[FILTERED]"
			}
		}
		b, err = json.Marshal(&result)
		if err != nil {
			return ""
		}
	} else if len(resultArray) != 0 {
		if !lp.ShowPassword {
			for _, v := range resultArray {
				if v["password"] != nil {
					v["password"] = "[FILTERED]"
				}
				if v["password_confirmation"] != nil {
					v["password_confirmation"] = "[FILTERED]"
				}
			}
		}
		b, err = json.Marshal(&resultArray)
		if err != nil {
			return ""
		}
	}

	str := string(b)
	str = strings.Replace(str, "\":\"", "\" => \"", -1)
	str = strings.Replace(str, "\":{", "\" => {", -1)
	str = strings.Replace(str, "\",\"", "\", \"", -1)
	str = strings.Replace(str, "},{", "}, {", -1)
	return fmt.Sprint(str)
}
