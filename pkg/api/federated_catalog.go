package api

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type FederatedCatalogAPI struct {
	client *resty.Client
	url    string
	apiKey string
}

func NewFederatedCatalogAPI() *FederatedCatalogAPI {
	return &FederatedCatalogAPI{
		client: resty.New(),
		url:    os.Getenv("FEDERATED_CATALOG_URL"),
		apiKey: os.Getenv("FEDERATED_CATALOG_API_KEY"),
	}
}

func (e *FederatedCatalogAPI) Query(querySpec QuerySpec) ([]map[string]interface{}, error) {
	var federatedCatalog []map[string]interface{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&federatedCatalog).
		Post(e.url + "/api/catalog/v1alpha/catalog/query")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get federated catalog: %s", resp.String())
	}

	return federatedCatalog, nil
}
