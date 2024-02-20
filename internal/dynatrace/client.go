package dynatrace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ConfigurationApiPath string = "/api/config/v1"

type DynatraceClient struct {
	ApiToken string
	EnvUrl   string
	Client   *http.Client
}

type GetExtensionsResponse struct {
	Extensions   []DynatraceExtensionInfo `json:"extensions"`
	TotalResults int                      `json:"totalResults"`
	NextPageKey  string                   `json:"nextPageKey"`
}

type DynatraceExtensionInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewClient(envUrl *string, apiToken *string) (*DynatraceClient, error) {

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	client := DynatraceClient{
		Client: &httpClient,
	}

	if envUrl == nil {
		return nil, errors.New("environment URL is missing")
	}

	if apiToken == nil {
		return nil, errors.New("API Token is missing")
	}

	client.ApiToken = *apiToken
	client.EnvUrl = *envUrl

	err := client.validateConnection()
	if(err != nil) {
		return nil, err
	}

	return &client, nil
}

func (c *DynatraceClient) validateConnection() error {
	_, err := c.getExtensions(10, nil)
	return err
}

func (c *DynatraceClient) getExtensions(pageSize int, nextPageKey *string) (*GetExtensionsResponse, error) {
	req, err := c.getConfigurationApiRequest("GET", "/extensions")
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Api-Token %s", c.ApiToken))
	req.Header.Set("accept", "application/json; charset=utf-8")

	queryParams := req.URL.Query()
	queryParams.Set("pageSize", fmt.Sprintf("%d",pageSize))

	if(nextPageKey != nil) {
		queryParams.Set("nextPageKey", *nextPageKey)
	}

	req.URL.RawQuery = queryParams.Encode()

	response, err := c.Client.Do(req)

	if(err != nil) {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if(err != nil) {
		return nil, err
	}

	if(response.StatusCode != 200) {
		return nil, fmt.Errorf("can't retrieve extensions: reason %d %s", response.StatusCode, string(responseBody))
	}

	var extensions GetExtensionsResponse
	err = json.Unmarshal(responseBody, &extensions)
	if(err != nil) {
		return nil, err
	}

	return &extensions, nil

	
}

func (c *DynatraceClient) getConfigurationApiRequest(method string, path string) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s%s", c.EnvUrl, ConfigurationApiPath, path), nil)
}
