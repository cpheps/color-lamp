package lampclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const colorEndpoint = "/v1/cluster/%s/color"

//LampClient handles connections to a server
type LampClient struct {
	serverAddress string
	port          string
}

//CreateLampClient creates a new LampClient
func CreateLampClient(serverAddress, port string) *LampClient {
	return &LampClient{
		serverAddress: serverAddress,
		port:          port,
	}
}

//GetClusterColor Retrieves the color from the cluster
func (lc LampClient) GetClusterColor(clusterName string) (*uint32, error) {
	req, err := lc.newRequest(lc.getServerAddress(fmt.Sprintf(colorEndpoint, clusterName)), http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	message := new(colorMessage)
	if _, err := lc.do(req, message); err != nil {
		return nil, err
	}

	return &message.Color, nil
}

//SetClusterColor sets the color of the given cluster
func (lc LampClient) SetClusterColor(clusterName string, color uint32) error {
	message := &colorMessage{
		Color: color,
	}

	req, err := lc.newRequest(lc.getServerAddress(fmt.Sprintf(colorEndpoint, clusterName)), http.MethodPut, message)
	if err != nil {
		return err
	}

	_, err = lc.do(req, nil)
	return err
}

func (lc LampClient) getServerAddress(endpoint string) string {
	return fmt.Sprintf("http://%s:%s%s", lc.serverAddress, lc.port, endpoint)
}

func (lc LampClient) newRequest(url, method string, body interface{}) (*http.Request, error) {
	var buff io.ReadWriter
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		buff = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do sends an API request and returns and API response. The API response
// is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error occurred.
func (lc LampClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		err = json.Unmarshal(data, v)
	}

	return resp, err
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	errorResponse := &errorResponse{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}
