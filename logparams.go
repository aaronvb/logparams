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

type ParamFields struct {
	Form      map[string]string
	Query     map[string]string
	Json      map[string]interface{}
	JsonArray []map[string]interface{}
}

// ToString will return a string of all parameters within the http request.
func (lp *LogParams) ToString() string {
	paramsString, _ := lp.parseParams()
	if !lp.ShowEmpty && paramsString == "" {
		return ""
	}

	var str string
	if lp.HidePrefix {
		str = paramsString
	} else {
		str = fmt.Sprintf("Parameters: %s", paramsString)
	}

	return str
}

// ToLogger will log print all parameters within the http request.
func (lp *LogParams) ToLogger(logger *log.Logger) {
	paramsString, _ := lp.parseParams()
	if !lp.ShowEmpty && paramsString == "" {
		return
	}

	var str string
	if lp.HidePrefix {
		str = paramsString
	} else {
		str = fmt.Sprintf("Parameters: %s", paramsString)
	}

	logger.Printf(str)
}

// ToLogger will log print all parameters within the http request.
func (lp *LogParams) ToFields() ParamFields {
	paramsString, fields := lp.parseParams()
	if !lp.ShowEmpty && paramsString == "" {
		return ParamFields{}
	}

	return fields
}

// Helper methods

// checkForFormParams checks for form params in the request.
func (lp *LogParams) checkForFormParams() bool {
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
func (lp *LogParams) checkForQueryParams() bool {
	return len(lp.Request.URL.Query()) != 0
}

// checkForJSON checks for content-type application/json in the header.
func (lp *LogParams) checkForJSON() bool {
	matched, _ := regexp.MatchString(`application\/json`, lp.Request.Header.Get("Content-Type"))
	return matched
}

// parseParams will check the type of param in the request and call the correct parser.
func (lp *LogParams) parseParams() (string, ParamFields) {
	if lp.checkForFormParams() {
		str, fields := lp.parseFormParams()
		str = fmt.Sprintf("{%s}", str)
		return str, fields
	} else if lp.checkForQueryParams() {
		str, fields := lp.parseQueryParams()
		str = fmt.Sprintf("{%s}", str)
		return str, fields
	} else if lp.checkForJSON() {
		str, fields := lp.parseJSONBody()
		return str, fields
	}

	return "", ParamFields{}
}

// parseFormParams will parse the form for values and return a string of parameters
func (lp *LogParams) parseFormParams() (string, ParamFields) {
	var paramString string

	err := lp.Request.ParseForm()
	if err != nil {
		return paramString, ParamFields{}
	}

	var paramCount = 0
	formFields := ParamFields{Form: make(map[string]string, len(lp.Request.PostForm))}
	for k := range lp.Request.PostForm {
		if k == "password" || k == "password_confirmation" {
			formFields.Form[k] = "[FILTERED]"
			paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, "[FILTERED]")
		} else {
			formValue := lp.Request.PostForm.Get(k)
			formFields.Form[k] = formValue
			paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, formValue)
		}
		paramCount++
		if paramCount != len(lp.Request.PostForm) {
			paramString += ", "
		}
	}

	return paramString, formFields
}

// parseQueryParams will parse query parameters in the URL.
func (lp *LogParams) parseQueryParams() (string, ParamFields) {
	var paramString string

	var paramCount = 0
	formFields := ParamFields{Query: make(map[string]string, len(lp.Request.URL.Query()))}
	for k := range lp.Request.URL.Query() {
		paramValue := lp.Request.URL.Query()[k][0]
		formFields.Query[k] = paramValue
		paramString += fmt.Sprintf("\"%s\" => \"%s\"", k, paramValue)
		paramCount++
		if paramCount != len(lp.Request.URL.Query()) {
			paramString += ", "
		}
	}

	return paramString, formFields
}

// parseJSONBody will parse the json in the body as parameters.
func (lp *LogParams) parseJSONBody() (string, ParamFields) {
	var result map[string]interface{}
	var resultArray []map[string]interface{}

	body, _ := ioutil.ReadAll(lp.Request.Body)
	lp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err := json.Unmarshal(body, &result)
	if err != nil {
		err := json.Unmarshal(body, &resultArray)
		if err != nil {
			return "", ParamFields{}
		}
	}

	var b []byte
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
			return "", ParamFields{}
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
			return "", ParamFields{}
		}
	}

	str := string(b)
	str = strings.Replace(str, "\":\"", "\" => \"", -1)
	str = strings.Replace(str, "\":{", "\" => {", -1)
	str = strings.Replace(str, "\",\"", "\", \"", -1)
	str = strings.Replace(str, "},{", "}, {", -1)
	fields := ParamFields{Json: result, JsonArray: resultArray}
	return fmt.Sprint(str), fields
}
