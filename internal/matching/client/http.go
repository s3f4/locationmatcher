package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/s3f4/locationmatcher/internal/matching/models"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type APIClient interface {
	FindNearest(url string, query *models.Query) (*apihelper.Response, error)
}

type apiClient struct {
	client *http.Client
}

var client *apiClient

func GetAPIClient() APIClient {
	if client == nil {
		client = new(apiClient)
		client.client = &http.Client{
			Timeout: time.Second * 15,
		}
	}

	return client
}

func (a *apiClient) FindNearest(url string, query *models.Query) (*apihelper.Response, error) {
	client := a.client

	newReq, err := json.Marshal(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(newReq)))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("X-REQUEST-FROM", "matching")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error %s", err)
		return nil, err
	}

	defer resp.Body.Close()
	var response apihelper.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
