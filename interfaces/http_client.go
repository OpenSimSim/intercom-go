package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

type HTTPClient interface {
	Get(context.Context, string, interface{}) ([]byte, error)
	Post(context.Context, string, interface{}) ([]byte, error)
	Patch(context.Context, string, interface{}) ([]byte, error)
	Delete(context.Context, string, interface{}) ([]byte, error)
}

type IntercomHTTPClient struct {
	*http.Client
	BaseURI       *string
	AppID         string
	APIKey        string
	ClientVersion *string
	Debug         *bool
}

func NewIntercomHTTPClient(appID, apiKey string, baseURI, clientVersion *string, debug *bool) IntercomHTTPClient {
	return IntercomHTTPClient{Client: &http.Client{}, AppID: appID, APIKey: apiKey, BaseURI: baseURI, ClientVersion: clientVersion, Debug: debug}
}

func (c IntercomHTTPClient) UserAgentHeader() string {
	return fmt.Sprintf("intercom-go/%s", *c.ClientVersion)
}

func (c IntercomHTTPClient) Get(ctx context.Context, url string, queryParams interface{}) ([]byte, error) {
	// Setup request
	req, _ := http.NewRequest("GET", *c.BaseURI+url, nil)
	req = req.WithContext(ctx)
	req.SetBasicAuth(c.AppID, c.APIKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgentHeader())
	addQueryParams(req, queryParams)
	if *c.Debug {
		fmt.Printf("%s %s\n", req.Method, req.URL)
	}

	// Do request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	data, err := c.readAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, c.parseResponseError(data, resp.StatusCode)
	}
	return data, err
}

func addQueryParams(req *http.Request, params interface{}) {
	v, _ := query.Values(params)
	req.URL.RawQuery = v.Encode()
}

func (c IntercomHTTPClient) Patch(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.postOrPatch(ctx, "PATCH", url, body)
}

func (c IntercomHTTPClient) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.postOrPatch(ctx, "POST", url, body)
}

func (c IntercomHTTPClient) postOrPatch(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	// Marshal our body
	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return nil, err
	}

	// Setup request
	req, err := http.NewRequest(method, *c.BaseURI+url, buffer)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.AppID, c.APIKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.UserAgentHeader())
	if *c.Debug {
		fmt.Printf("%s %s %s\n", req.Method, req.URL, buffer)
	}

	// Do request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	data, err := c.readAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, c.parseResponseError(data, resp.StatusCode)
	}
	return data, err
}

func (c IntercomHTTPClient) Delete(ctx context.Context, url string, queryParams interface{}) ([]byte, error) {
	// Setup request
	req, _ := http.NewRequest("DELETE", *c.BaseURI+url, nil)
	req = req.WithContext(ctx)
	req.SetBasicAuth(c.AppID, c.APIKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgentHeader())
	addQueryParams(req, queryParams)
	if *c.Debug {
		fmt.Printf("%s %s\n", req.Method, req.URL)
	}

	// Do request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	data, err := c.readAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, c.parseResponseError(data, resp.StatusCode)
	}
	return data, err
}

type IntercomError interface {
	Error() string
	GetStatusCode() int
	GetCode() string
	GetMessage() string
}

func (c IntercomHTTPClient) parseResponseError(data []byte, statusCode int) IntercomError {
	errorList := HTTPErrorList{}
	err := json.Unmarshal(data, &errorList)
	if err != nil {
		return NewUnknownHTTPError(statusCode)
	}
	if len(errorList.Errors) == 0 {
		return NewUnknownHTTPError(statusCode)
	}
	httpError := errorList.Errors[0]
	httpError.StatusCode = statusCode
	return httpError // only care about the first
}

func (c IntercomHTTPClient) readAll(body io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if *c.Debug {
		fmt.Println(string(b))
		fmt.Println("")
	}
	return b, err
}
