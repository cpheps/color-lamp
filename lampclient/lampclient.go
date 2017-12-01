package lampclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jmoiron/jsonq"
)

const colorEndpoint = "/color"
const clusterEndpoint = "/cluster"

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
func (lc LampClient) GetClusterColor(clusterID *string) (*int32, error) {
	body := fmt.Sprintf("{\"id\": \"%s\"}", *clusterID)

	response, err := lc.makeRequest(http.MethodGet, colorEndpoint, &body)
	if err != nil {
		return nil, err
	} else if response.StatusCode >= 400 {
		errorMessage, err := getErrorMessage(response)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Returned with status code %d, and error message %s", response.StatusCode, *errorMessage)
	}

	jq, err := parseJSON(response)
	if err != nil {
		return nil, err
	}

	color, err := jq.Int("color")
	if err != nil {
		return nil, errors.New("Property 'color' not found in JSON response")
	}

	color32 := int32(color)
	return &color32, nil
}

//SetClusterColor sets the color of the given cluster
func (lc LampClient) SetClusterColor(clusterID *string, color *int32) error {
	body := fmt.Sprintf("{\"id\": \"%s\",\"color\": %d}", *clusterID, *color)

	response, err := lc.makeRequest(http.MethodPut, colorEndpoint, &body)
	if err != nil {
		return err
	} else if response.StatusCode >= 400 {
		errorMessage, err := getErrorMessage(response)
		if err != nil {
			return err
		}
		return fmt.Errorf("Returned with status code %d, and error message %s", response.StatusCode, *errorMessage)
	}

	return nil
}

func (lc LampClient) getServerAddress(endpoint string) string {
	return fmt.Sprintf("http://%s:%s%s", lc.serverAddress, lc.port, endpoint)
}

func formatBody(content *string) io.Reader {
	if content == nil {
		return nil
	}

	var buffer bytes.Buffer
	buffer.WriteString(*content)
	return bytes.NewReader(buffer.Bytes())
}

func (lc LampClient) makeRequest(method, endpoint string, body *string) (*http.Response, error) {
	request, err := http.NewRequest(method, lc.getServerAddress(endpoint), formatBody(body))
	if err != nil {
		return nil, fmt.Errorf("Unable to make %s %s request. Error: %s", method, endpoint, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			err := ctx.Err()

			if err.Error() == "context deadline exceeded" {
				fmt.Printf("Timeout of %s %s request.\n", method, endpoint) // prints "context deadline exceeded"
			}
		}
	}()

	request = request.WithContext(ctx)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Unable to make %s %s request. Error: %s", method, endpoint, err.Error())
	}

	return resp, nil
}

func parseJSON(response *http.Response) (*jsonq.JsonQuery, error) {

	resp := make(map[string]interface{})
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON. Error: %s", err.Error())
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON. Error: %s", err.Error())
	}

	jq := jsonq.NewQuery(resp)

	return jq, nil
}

func getErrorMessage(response *http.Response) (*string, error) {
	jq, err := parseJSON(response)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse error: Error %s", err.Error())
	}

	errorMessage, err := jq.String("error")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse error: Error %s", err.Error())
	}

	return &errorMessage, nil
}
