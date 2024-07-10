package api

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type EdcReceiverEndpointApi struct {
	client *resty.Client
	url    string
}

type DataRequest struct {
	Id         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	ContractId string                 `json:"contractId,omitempty"`
	Endpoint   string                 `json:"endpoint,omitempty"`
	AuthKey    string                 `json:"authKey,omitempty"`
	AuthCode   string                 `json:"authCode,omitempty"`
}

func NewEdcReceiverEndpointApi() *EdcReceiverEndpointApi {
	return &EdcReceiverEndpointApi{
		client: resty.New(),
		url:    os.Getenv("EDC_RECEIVER_ENDPOINT_URL"),
	}
}

func (e *EdcReceiverEndpointApi) GetDataRequest(id string) (DataRequest, error) {
	var dataRequest DataRequest

	resp, err := e.client.R().
		SetResult(&dataRequest).
		Get(e.url + "/" + id)

	fmt.Println(e.url + "/" + id)
	fmt.Println(resp)

	if err != nil {
		return DataRequest{}, err
	}

	if resp.StatusCode() != 200 {
		return DataRequest{}, err
	}

	return dataRequest, nil
}
