package api

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type ThingsboardId struct {
	ID         string `json:"id"`
	EntityType string `json:"entityType"`
}

type ThingsboardCustomer struct {
	Id       *ThingsboardId `json:"id,omitempty"`
	Title    string         `json:"title"`
	Name     string         `json:"name"`
	Email    string         `json:"email,omitempty"`
	Phone    string         `json:"phone,omitempty"`
	Country  string         `json:"country,omitempty"`
	State    string         `json:"state,omitempty"`
	City     string         `json:"city,omitempty"`
	Address  string         `json:"address,omitempty"`
	Address2 string         `json:"address2,omitempty"`
	Zip      string         `json:"zip,omitempty"`
}

type ThingsboardDevice struct {
	Name           string                 `json:"name"`
	AdditionalInfo map[string]interface{} `json:"additionalInfo"`
}

type ExchangeTokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"type,omitempty"`
}

type ThingsboardAPI struct {
	client      *resty.Client
	url         string
	exchangeUrl string
}

func NewThingsboardAPI() *ThingsboardAPI {
	return &ThingsboardAPI{
		client:      resty.New(),
		url:         os.Getenv("THINGSBOARD_URL"),
		exchangeUrl: os.Getenv("THINGSBOARD_EXCHANGE_TOKEN_URL"),
	}
}

func (tb *ThingsboardAPI) ExchangeToken(accessToken string) (string, error) {
	var exchangeTokenResponse ExchangeTokenResponse

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username": "oauth2-token", "password": "` + accessToken + `"}`).
		SetResult(&exchangeTokenResponse).
		Post(tb.exchangeUrl + "/api/auth/login")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("unable to exchange token: %s", resp.String())
	}

	return exchangeTokenResponse.Token, nil
}

func (tb *ThingsboardAPI) CreateCustomer(accessToken string, customer ThingsboardCustomer) (ThingsboardCustomer, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return ThingsboardCustomer{}, err
	}

	var createdCustomer ThingsboardCustomer

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetBody(customer).
		SetResult(&createdCustomer).
		Post(tb.url + "/api/customer")

	if err != nil {
		return ThingsboardCustomer{}, err
	}

	if resp.StatusCode() != 200 {
		return ThingsboardCustomer{}, fmt.Errorf("unable to create customer: %s", resp.String())
	}

	return createdCustomer, nil
}

func (tb *ThingsboardAPI) DeleteCustomer(accessToken string, customerID string) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		Delete(tb.url + "/api/customer/" + customerID)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to delete customer: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) GetCustomer(accessToken string, customerID string) (ThingsboardCustomer, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return ThingsboardCustomer{}, err
	}

	var customer ThingsboardCustomer

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&customer).
		Get(tb.url + "/api/customer/" + customerID)

	if err != nil {
		return ThingsboardCustomer{}, err
	}

	if resp.StatusCode() != 200 {
		return ThingsboardCustomer{}, fmt.Errorf("unable to get customer: %s", resp.String())
	}

	return customer, nil
}

func (tb *ThingsboardAPI) GetCustomerDevices(accessToken string, customerID string) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var devices map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&devices).
		SetQueryParams(map[string]string{
			"page":     "0",
			"pageSize": "1000",
		}).
		Get(tb.url + "/api/customer/" + customerID + "/devices")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get customer devices: %s", resp.String())
	}

	return devices, nil
}

func (tb *ThingsboardAPI) GetTenantDevices(accessToken string) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var devices map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&devices).
		SetQueryParams(map[string]string{
			"page":     "0",
			"pageSize": "1000",
		}).
		Get(tb.url + "/api/tenant/devices")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get devices: %s", resp.String())
	}

	return devices, nil
}

func (tb *ThingsboardAPI) GetDevice(accessToken string, deviceID string) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var device map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&device).
		Get(tb.url + "/api/device/" + deviceID)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get device: %s", resp.String())
	}

	return device, nil
}

func (tb *ThingsboardAPI) CreateDevice(accessToken string, device ThingsboardDevice) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var createdDevice map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetBody(device).
		SetResult(&createdDevice).
		Post(tb.url + "/api/device")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create device: %s", resp.String())
	}

	return createdDevice, nil
}

func (tb *ThingsboardAPI) GetProfiles(accessToken string) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var profiles map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&profiles).
		SetQueryParams(map[string]string{
			"page":     "0",
			"pageSize": "1000",
		}).
		Get(tb.url + "/api/deviceProfiles")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get profiles: %s", resp.String())
	}

	return profiles, nil
}

func (tb *ThingsboardAPI) DeleteDevice(accessToken string, deviceID string) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		Delete(tb.url + "/api/device/" + deviceID)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to delete device: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) GetDeviceAttributes(accessToken string, deviceID string) ([]map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var attributes []map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&attributes).
		Get(tb.url + "/api/plugins/telemetry/DEVICE/" + deviceID + "/values/attributes")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get attributes: %s", resp.String())
	}

	return attributes, nil
}

func (tb *ThingsboardAPI) CreateDeviceAttributes(accessToken string, deviceID string, attributes map[string]interface{}) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetBody(attributes).
		Post(tb.url + "/api/plugins/telemetry/DEVICE/" + deviceID + "/attributes/SERVER_SCOPE")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to create attributes: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) DeleteDeviceAttribute(accessToken string, deviceID string, attribute string) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		Delete(tb.url + "/api/plugins/telemetry/DEVICE/" + deviceID + "/SERVER_SCOPE?keys=" + attribute)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to delete attribute: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) GetDeviceCredentials(accessToken string, deviceID string) (map[string]interface{}, error) {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var credentials map[string]interface{}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetResult(&credentials).
		Get(tb.url + "/api/device/" + deviceID + "/credentials")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get credentials: %s", resp.String())
	}

	return credentials, nil
}

func (tb *ThingsboardAPI) UpdateDevice(accessToken string, deviceID string, device map[string]interface{}) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetBody(device).
		Post(tb.url + "/api/device")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to update device: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) CreateUser(accessToken string, email string, firstName string, lastName string, tenantId string, customerId string) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	user := map[string]interface{}{
		"email":     email,
		"firstName": firstName,
		"lastName":  lastName,
		"tenantId": map[string]string{
			"id":         tenantId,
			"entityType": "TENANT",
		},
		"customerId": map[string]string{
			"id":         customerId,
			"entityType": "CUSTOMER",
		},
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		SetBody(user).
		Post(tb.url + "/api/user")

	fmt.Println(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to create user: %s", resp.String())
	}

	return nil
}

func (tb *ThingsboardAPI) DeleteUser(accessToken string, userID string) error {
	thingsboardToken, err := tb.ExchangeToken(accessToken)
	if err != nil {
		return err
	}

	resp, err := tb.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Authorization", "Bearer "+thingsboardToken).
		Delete(tb.url + "/api/user/" + userID)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to delete user: %s", resp.String())
	}

	return nil
}
