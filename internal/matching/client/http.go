package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/s3f4/locationmatcher/internal/matching/models"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type APIClient interface {
	FindNearest(context.Context, string, *models.Query) (*apihelper.Response, error)
}

type httpClient struct {
	client *http.Client
}

var client *httpClient
var circuit Circuit

func NewAPIClient() APIClient {
	if client == nil {
		client = new(httpClient)
		client.client = &http.Client{
			Timeout: time.Second * 15,
		}

		circuit = Breaker(func(ctx context.Context, url string, reader io.Reader) (*http.Response, error) {
			req, err := http.NewRequest("POST", url, reader)
			if err != nil {
				return nil, err
			}

			req = req.WithContext(ctx)
			req.Header.Set("X-USER-AUTHENTICATED", "true")

			resp, err := client.client.Do(req)
			if err != nil {
				return nil, err
			}

			return resp, nil
		}, 5)
	}

	return client
}

func (a *httpClient) FindNearest(ctx context.Context, url string, query *models.Query) (*apihelper.Response, error) {
	newReq, err := json.Marshal(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	reader := io.NopCloser(bytes.NewReader(newReq))
	resp, err := circuit(ctx, url, reader)
	if err != nil {
		// log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var response apihelper.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
