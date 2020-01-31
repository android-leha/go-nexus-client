package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testScriptsAPIEndpoint = "/service/rest/v1/script"
)

func testScript(name string) *Script {
	return &Script{
		Name:    name,
		Content: fmt.Sprintf("log.info('Hello, %s!')", name),
		Type:    "groovy",
	}
}

func testHTTPHeader(t *testing.T, method string, path string, req *http.Request) {
	assert.Equal(t, req.Method, method)
	assert.Equal(t, path, req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	user, pass, ok := req.BasicAuth()
	assert.Equal(t, "admin", user)
	assert.Equal(t, "admin123", pass)
	assert.True(t, ok)
}

func TestScriptList(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "GET", testScriptsAPIEndpoint, req)

		// Send response
		rw.Write([]byte(`[ {
			"name" : "helloWorld",
			"content" : "log.info('Hello, World!')",
			"type" : "groovy"
		  },
		  {
			"name" : "HelloTest",
			"content" : "log.info('Hello, Test!')",
			"type" : "groovy"
		  }]`))
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	scripts, err := client.ScriptLists()
	assert.Nil(t, err)
	expectedScripts := []Script{
		Script{
			Name:    "helloWorld",
			Content: "log.info('Hello, World!')",
			Type:    "groovy",
		},
		Script{
			Name:    "HelloTest",
			Content: "log.info('Hello, Test!')",
			Type:    "groovy",
		},
	}
	assert.Equal(t, &expectedScripts, scripts)
}

func TestScriptRead(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "GET", fmt.Sprintf("%s/helloWorld", testScriptsAPIEndpoint), req)

		// Send response
		rw.Write([]byte(`{
			"name" : "helloWorld",
			"content" : "log.info('Hello, World!')",
			"type" : "groovy"
		  }`))
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	script, err := client.ScriptRead("helloWorld")
	assert.Nil(t, err)
	expectedScript := Script{
		Name:    "helloWorld",
		Content: "log.info('Hello, World!')",
		Type:    "groovy",
	}
	assert.Equal(t, &expectedScript, script)
}

func TestScriptCreate(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "POST", testScriptsAPIEndpoint, req)

		expectedBody := `{"name":"test-script-create","content":"log.info('Hello, test-script-create!')","type":"groovy"}`
		bodyBytes, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, expectedBody, string(bodyBytes))

		// Send response
		rw.WriteHeader(204)
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	testScript := testScript("test-script-create")
	err := client.ScriptCreate(testScript)
	assert.Nil(t, err)
}

func TestScriptUpdate(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "PUT", fmt.Sprintf("%s/test-script-update", testScriptsAPIEndpoint), req)

		expectedBody := `{"name":"test-script-update","content":"log.info('Hello, test-script-update!')","type":"groovy"}`
		bodyBytes, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, expectedBody, string(bodyBytes))

		// Send response
		rw.WriteHeader(204)
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	testScript := testScript("test-script-update")
	err := client.ScriptUpdate(testScript)
	assert.Nil(t, err)
}

func TestScriptRUN(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "POST", fmt.Sprintf("%s/test-script-update/run", testScriptsAPIEndpoint), req)

		// Send response
		rw.Write([]byte(`{
			"name" : "test-script-update",
			"result" : "null"
		  }`))
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	err := client.ScriptRun("test-script-update")
	assert.Nil(t, err)
}

func TestScriptDelete(t *testing.T) {
	// Start a local HTTP server
	testserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testHTTPHeader(t, "DELETE", fmt.Sprintf("%s/test-script-delete", testScriptsAPIEndpoint), req)

		// Send response
		rw.WriteHeader(204)
	}))
	// Close the server when test finishes
	defer testserver.Close()

	client := NewClient(Config{
		URL:      testserver.URL,
		Username: "admin",
		Password: "admin123",
	})
	err := client.ScriptDelete("test-script-delete")
	assert.Nil(t, err)
}
