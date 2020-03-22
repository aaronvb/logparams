package logparams

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestToString(t *testing.T) {
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

func TestToStringIsEmpty(t *testing.T) {
	expectedResults := ""

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, HideEmpty: true}
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

func TestToLogger(t *testing.T) {
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

func TestToLoggerEmpty(t *testing.T) {
	var str bytes.Buffer
	var logger = log.Logger{}
	logger.SetOutput(&str)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lp := LogParams{Request: r, HideEmpty: true}
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

