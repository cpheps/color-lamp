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

const serverAddress = "localhost"
const port = "5000"

const colorEndpoint = "/color"
const clusterEndpoint = "/cluster"

//GetClusterColor Retrieves the color from the cluster
func GetClusterColor(clusterID *string) (*int32, error) {
	body := fmt.Sprintf("{\"id\": \"%s\"}", *clusterID)

	response, err := makeRequest(http.MethodGet, colorEndpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("Unable to make GET /color request. Error: %s", err.Error())
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

func getServerAddress(endpoint string) string {
	return fmt.Sprintf("http://%s:%s%s", serverAddress, port, endpoint)
}

func formatBody(content *string) io.Reader {
	if content == nil {
		return nil
	}

	var buffer bytes.Buffer
	buffer.WriteString(*content)
	return bytes.NewReader(buffer.Bytes())
}

func makeRequest(method, endpoint string, body *string) (*http.Response, error) {
	request, err := http.NewRequest(method, getServerAddress(endpoint), formatBody(body))
	if err != nil {
		return nil, fmt.Errorf("Unable to make %s %s request. Error: %s", method, endpoint, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			fmt.Printf("Timeout of %s %s request.\n", method, endpoint) // prints "context deadline exceeded"
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
