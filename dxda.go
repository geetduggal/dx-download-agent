package dxda

// Some inspiration + code snippets taken from https://github.com/dnanexus/precision-fda/blob/master/go/pfda.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"     // required by go-retryablehttp
	"github.com/hashicorp/go-retryablehttp" // use http libraries from hashicorp for implement retry logic
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func urlFailure(requestType string, url string, status string) {
	log.Fatalln(fmt.Errorf("%s request to '%s' failed with status %s", requestType, url, status))
}

// Utilities to interact with the DNAnexus API
// TODO: Create automatic API wrappers for the dx toolkit
// e.g. via: https://github.com/dnanexus/dx-toolkit/tree/master/src/api_wrappers

// As much as I love Go, see https://mholt.github.io/json-to-go/
// for auto-generation

// DXConfig - Basic variables regarding DNAnexus environment config
type DXConfig struct {
	DXSECURITYCONTEXT    string `json:"DX_SECURITY_CONTEXT"`
	DXAPISERVERHOST      string `json:"DX_APISERVER_HOST"`
	DXPROJECTCONTEXTNAME string `json:"DX_PROJECT_CONTEXT_NAME"`
	DXPROJECTCONTEXTID   string `json:"DX_PROJECT_CONTEXT_ID"`
	DXAPISERVERPORT      string `json:"DX_APISERVER_PORT"`
	DXUSERNAME           string `json:"DX_USERNAME"`
	DXAPISERVERPROTOCOL  string `json:"DX_APISERVER_PROTOCOL"`
	DXCLIWD              string `json:"DX_CLI_WD"`
}

// DXAuthorization - Basic variables regarding DNAnexus authorization
type DXAuthorization struct {
	AuthToken     string `json:"auth_token"`
	AuthTokenType string `json:"auth_token_type"`
}

// GetToken - Get DNAnexus authentication token
/*
   Returns a pair of strings representing the authentication token and where it was received from
   If the environment variable DX_API_TOKEN is set, the token is obtained from it
   Otherwise, the token is obtained from the '~/.dnanexus_config/environment.json' file
   If no token can be obtained from these methods, a pair of empty strings is returned
   If the token was received from the 'DX_API_TOKEN' environment variable, the second variable in the pair
   will be the string 'environment'.  If it is obtained from a DNAnexus configuration file, the second variable
   in the pair will be '.dnanexus_config/environment.json'.
*/
func GetToken() (string, string) {
	envToken := os.Getenv("DX_API_TOKEN")
	envFile := fmt.Sprintf("%s/.dnanexus_config/environment.json", os.Getenv("HOME"))
	if envToken != "" {
		return envToken, "environment"
	}
	if _, err := os.Stat(envFile); err == nil {
		config, _ := ioutil.ReadFile(envFile)
		var dxconf DXConfig
		json.Unmarshal(config, &dxconf)
		var dxauth DXAuthorization
		json.Unmarshal([]byte(dxconf.DXSECURITYCONTEXT), &dxauth)
		return dxauth.AuthToken, "~/.dnanexus_config/environment.json"
	}
	return "", ""
}

func makeRequestWithHeadersFail(requestType string, url string, headers map[string]string, data []byte) (status string, body []byte) {
	const minRetryTime = 1  // seconds
	const maxRetryTime = 30 // seconds
	const maxRetryCount = 5
	const userAgent = "DNAnexus Download Agent (v. 0.1)"

	client := &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultClient(),
		Logger:       log.New(ioutil.Discard, "", 0), // Throw away retryablehttp internal logging
		RetryWaitMin: minRetryTime * time.Second,
		RetryWaitMax: maxRetryTime * time.Second,
		RetryMax:     maxRetryCount,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
	}

	req, err := retryablehttp.NewRequest(requestType, url, bytes.NewReader(data))
	check(err)
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	status = resp.Status
	body, _ = ioutil.ReadAll(resp.Body)

	if !strings.HasPrefix(status, "2") {
		urlFailure(requestType, url, status)
	}
	return status, body
}

// GenericDXAPI (WIP) - Function to wrap a generic API call to DNAnexus
func GenericDXAPI(token, api string, payload []byte) (status string, body []byte) {
	headers := map[string]string{
		"User-Agent":    "DNAnexus Download Client v0.1",
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
	url := fmt.Sprintf("https://api.dnanexus.com/%s", api)
	return makeRequestWithHeadersFail("POST", url, headers, payload)
}

// WhoAmI - TODO: Should the token be abstracted into a struct that is reused with other methods more like a class?
func WhoAmI(token string) string {
	_, body := GenericDXAPI(token, "system/whoami", []byte("{}"))
	return string(body)
}
