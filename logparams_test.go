package logparams

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Post Form parameters

func TestPostFormToString(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	params := url.Values{}
	params.Set("foo", "bar")

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestPostFormToStringIsEmpty(t *testing.T) {
	expectedResults := ""

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, ShowEmpty: false}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	params := url.Values{}

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestPostFormToStringShowEmpty(t *testing.T) {
	expectedResults := "Parameters: "

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, ShowEmpty: true}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	params := url.Values{}

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestPostFormToLogger(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, expectedResults)
		}
	}))

	defer server.Close()

	params := url.Values{}
	params.Set("foo", "bar")

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestPostFormToLoggerEmpty(t *testing.T) {
	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, ShowEmpty: false}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != "" {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, "")
		}
	}))

	defer server.Close()

	params := url.Values{}

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestPostFormToLoggerShowEmpty(t *testing.T) {
	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, ShowEmpty: true}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != "Parameters: " {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, "Parameters: ")
		}
	}))

	defer server.Close()

	params := url.Values{}

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

// Query parameters

func TestParseQueryParamsToString(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	_, err := http.Get(server.URL + "?foo=bar")
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestParseQueryParamsToLogger(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, expectedResults)
		}
	}))

	defer server.Close()

	_, err := http.Get(server.URL + "?foo=bar")
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

// JSON body

func TestParseJSONBodyToString(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`{"foo":"bar"}`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestParseJSONBodyToLogger(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, expectedResults)
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`{"foo":"bar"}`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestParseJSONArrayBodyToString(t *testing.T) {
	expectedResults := "Parameters: [{\"foo\" => \"bar\"}]"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`[{"foo":"bar"}]`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestParseJSONArrayBodyToLogger(t *testing.T) {
	expectedResults := "Parameters: [{\"foo\" => \"bar\"}]"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, expectedResults)
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`[{"foo":"bar"}]`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestFormPasswordIsFilteredByDefault(t *testing.T) {
	expectedResults := "Parameters: {\"password\" => \"[FILTERED]\"}"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		result := strings.TrimSuffix(str.String(), "\n")
		if result != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", result, expectedResults)
		}
		if lp.Request.PostForm.Get("password") != "foo" {
			t.Errorf("Expected attribute was incorrect, got %s, want: %s", lp.Request.PostForm.Get("password"), "foo")
		}
	}))

	defer server.Close()

	params := url.Values{}
	params.Set("password", "foo")

	_, err := http.PostForm(server.URL, params)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestJSONBodyPasswordIsFilteredByDefault(t *testing.T) {
	expectedResults := "Parameters: {\"password\" => \"[FILTERED]\"}"

	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		lp.ToLogger(&logger)
		logResult := strings.TrimSuffix(str.String(), "\n")
		if logResult != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", logResult, expectedResults)
		}

		var actualResult map[string]interface{}
		body, _ := ioutil.ReadAll(lp.Request.Body)
		lp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err := json.Unmarshal(body, &actualResult)
		if err != nil {
			t.Error("Error parsing JSON body")
		}
		if actualResult["password"] != "foobar" {
			t.Errorf("Expected attribute was incorrect, got %s, want: %s", actualResult["password"], "foobar")
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`{"password":"foobar"}`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestJSONArrayPasswordIsFilteredByDefault(t *testing.T) {
	expectedResults := "Parameters: [{\"password\" => \"[FILTERED]\"}]"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}

		var actualResultArray []map[string]interface{}
		body, _ := ioutil.ReadAll(lp.Request.Body)
		lp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err := json.Unmarshal(body, &actualResultArray)
		if err != nil {
			t.Error("Error parsing JSON body")
		}
		for _, v := range actualResultArray {
			if v["password"] != "foobar" {
				t.Errorf("Expected attribute was incorrect, got %s, want: %s", v["password"], "foobar")
			}
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`[{"password":"foobar"}]`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}

func TestParseJSONBodyCreatesNewReqBuffer(t *testing.T) {
	expectedResults := "Parameters: {\"foo\" => \"bar\"}"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}

		lp = LogParams{Request: r}
		if lp.ToString() != expectedResults {
			t.Errorf("Expected string was incorrect, got %s, want: %s", lp.ToString(), expectedResults)
		}
	}))

	defer server.Close()

	var jsonStr = []byte(`{"foo":"bar"}`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Errorf("Error POST to httptest server")
	}
}
