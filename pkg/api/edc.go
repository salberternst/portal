package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type Criterion struct {
	Type         string      `json:"@type"`
	OperandLeft  string      `json:"operandLeft"`
	OperandRight interface{} `json:"operandRight"`
	Operator     string      `json:"operator"`
}

type QuerySpec struct {
	Context map[string]string `json:"@context"`
	Type    string            `json:"@type"`
	Offset  uint              `json:"offset,omitempty"`
	Limit   uint              `json:"limit,omitempty"`
	// SortOrder        string            `json:"sortOrder"`
	// SortField        string            `json:"sortField"`
	FilterExpression []Criterion `json:"filterExpression"`
}

type Constraint struct {
	Edctype string `json:"edctype"`
}

type Action struct {
	Constraint *Constraint `json:"constraint,omitempty"`
	IncludedIn string      `json:"includedIn,omitempty"`
	Type_      string      `json:"type,omitempty"`
	Id         string      `json:"@id,omitempty"`
}

type Prohibition struct {
	Action      *Action      `json:"action,omitempty"`
	Assignee    string       `json:"assignee,omitempty"`
	Assigner    string       `json:"assigner,omitempty"`
	Constraints []Constraint `json:"constraints,omitempty"`
	Target      string       `json:"target,omitempty"`
	Uid         string       `json:"uid,omitempty"`
}

type Duty struct {
	Action           *Action      `json:"action,omitempty"`
	Assignee         string       `json:"assignee,omitempty"`
	Assigner         string       `json:"assigner,omitempty"`
	Consequence      *Duty        `json:"consequence,omitempty"`
	Constraints      []Constraint `json:"constraints,omitempty"`
	ParentPermission *Permission  `json:"parentPermission,omitempty"`
	Target           string       `json:"target,omitempty"`
	Uid              string       `json:"uid,omitempty"`
}

type Permission struct {
	Action      *Action      `json:"odrl:action,omitempty"`
	Assignee    string       `json:"assignee,omitempty"`
	Assigner    string       `json:"assigner,omitempty"`
	Constraints []Constraint `json:"constraints,omitempty"`
	Duties      []Duty       `json:"duties,omitempty"`
	Target      string       `json:"target,omitempty"`
	Uid         string       `json:"uid,omitempty"`
}

type Policy struct {
	Context              *interface{}           `json:"@context,omitempty"`
	Type                 string                 `json:"@type,omitempty"`
	Assignee             string                 `json:"odrl:assignee,omitempty"`
	Assigner             string                 `json:"odrl:assigner,omitempty"`
	ExtensibleProperties map[string]interface{} `json:"odrl:extensibleProperties,omitempty"`
	InheritsFrom         string                 `json:"odrl:inheritsFrom,omitempty"`
	Obligations          Duties                 `json:"odrl:obligation,omitempty"`
	Permissions          Permissions            `json:"odrl:permission,omitempty"`
	Prohibitions         Prohibitions           `json:"odrl:prohibition,omitempty"`
	Target               interface{}            `json:"odrl:target,omitempty"`
}

type Permissions []Permission

func (a *Permissions) UnmarshalJSON(data []byte) error {
	var single Permission
	if err := json.Unmarshal(data, &single); err == nil {
		*a = Permissions{single}
		return nil
	}

	var multiple []Permission
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("failed to unmarshal Permissions")
}

type Duties []Duty

func (a *Duties) UnmarshalJSON(data []byte) error {
	var single Duty
	if err := json.Unmarshal(data, &single); err == nil {
		*a = Duties{single}
		return nil
	}

	var multiple []Duty
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("failed to unmarshal Duties")
}

type Prohibitions []Prohibition

func (a *Prohibitions) UnmarshalJSON(data []byte) error {
	var single Prohibition
	if err := json.Unmarshal(data, &single); err == nil {
		*a = Prohibitions{single}
		return nil
	}

	var multiple []Prohibition
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("failed to unmarshal Prohibitions")
}

type PolicyDefinition struct {
	Context              *interface{}           `json:"@context"`
	CreatedAt            uint                   `json:"createdAt"`
	ID                   string                 `json:"@id"`
	Type                 string                 `json:"@type"`
	Policy               Policy                 `json:"policy"`
	PrivateProperties    map[string]string      `json:"privateProperties,omitempty"`
	ExtensibleProperties map[string]interface{} `json:"extensibleProperties,omitempty"`
}

type Asset struct {
	Context           *interface{}      `json:"@context"`
	Id                string            `json:"@id,omitempty"`
	Type              string            `json:"@type,omitempty"`
	DataAddress       map[string]string `json:"dataAddress"`
	PrivateProperties map[string]string `json:"privateProperties,omitempty"`
	Properties        map[string]string `json:"properties"`
}

type ContractDefinition struct {
	Context           *interface{}      `json:"@context"`
	Id                string            `json:"@id,omitempty"`
	Type              string            `json:"@type,omitempty"`
	AccessPolicyId    string            `json:"accessPolicyId"`
	AssetsSelector    AssetSelector     `json:"assetsSelector,omitempty"`
	ContractPolicyId  string            `json:"contractPolicyId"`
	PrivateProperties map[string]string `json:"privateProperties,omitempty"`
}

type AssetSelector []Criterion

func (a *AssetSelector) UnmarshalJSON(data []byte) error {
	var single Criterion
	if err := json.Unmarshal(data, &single); err == nil {
		*a = AssetSelector{single}
		return nil
	}

	var multiple []Criterion
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("failed to unmarshal AssetSelector")
}

type CallbackAddress struct {
	Type          string   `json:"@type,omitempty"`
	AuthCodeId    string   `json:"authCodeId,omitempty"`
	AuthKey       string   `json:"authKey,omitempty"`
	Events        []string `json:"events,omitempty"`
	Transactional bool     `json:"transactional,omitempty"`
	Uri           string   `json:"uri,omitempty"`
}

type CallbackAddressSelector []CallbackAddress

func (a *CallbackAddressSelector) UnmarshalJSON(data []byte) error {
	var single CallbackAddress
	if err := json.Unmarshal(data, &single); err == nil {
		*a = CallbackAddressSelector{single}
		return nil
	}

	var multiple []CallbackAddress
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("failed to unmarshal CallbackAddressSelector")
}

type ContractNegotiation struct {
	Context             *interface{}            `json:"@context"`
	Id                  string                  `json:"@id,omitempty"`
	Type_               string                  `json:"@type,omitempty"`
	Type                string                  `json:"type,omitempty"`
	CallbackAddresses   CallbackAddressSelector `json:"callbackAddresses,omitempty"`
	ContractAgreementId string                  `json:"contractAgreementId,omitempty"`
	CounterPartyAddress string                  `json:"counterPartyAddress,omitempty"`
	CounterPartyId      string                  `json:"counterPartyId,omitempty"`
	ErrorDetail         string                  `json:"errorDetail,omitempty"`
	Protocol            string                  `json:"protocol,omitempty"`
	State               string                  `json:"state,omitempty"`
}

type ContractOfferDescription struct {
	Type    string  `json:"@type,omitempty"`
	AssetId string  `json:"assetId,omitempty"`
	OfferId string  `json:"offerId,omitempty"`
	Policy  *Policy `json:"policy,omitempty"`
}

type Offer struct {
	Context  *interface{} `json:"@context"`
	Id       string       `json:"@id"`
	Type_    string       `json:"@type,omitempty"`
	Assigner string       `json:"assigner"`
	Target   string       `json:"target"`
}

type ContractRequest struct {
	Context             *interface{}              `json:"@context"`
	Type                string                    `json:"@type,omitempty"`
	CallbackAddresses   CallbackAddressSelector   `json:"callbackAddresses,omitempty"`
	ConnectorAddress    string                    `json:"connectorAddress,omitempty"`
	CounterPartyAddress string                    `json:"counterPartyAddress"`
	Offer               *ContractOfferDescription `json:"offer,omitempty"`
	Policy              *Offer                    `json:"policy"`
	Protocol            string                    `json:"protocol"`
	ProviderId          string                    `json:"providerId,omitempty"`
}

type ContractAgreement struct {
	Id                  string                  `json:"@id,omitempty"`
	Type                string                  `json:"@type,omitempty"`
	AssetId             string                  `json:"assetId,omitempty"`
	ConsumerId          string                  `json:"consumerId,omitempty"`
	ContractSigningDate int64                   `json:"contractSigningDate,omitempty"`
	Policy              *Policy                 `json:"policy,omitempty"`
	ProviderId          string                  `json:"providerId,omitempty"`
	PrivateProperties   map[string]string       `json:"privateProperties,omitempty"`
	CallbackAddresses   CallbackAddressSelector `json:"callbackAddresses,omitempty"`
}

type CatalogRequest struct {
	Context             *interface{} `json:"@context"`
	Type                string       `json:"@type,omitempty"`
	CounterPartyAddress string       `json:"counterPartyAddress"`
	CounterPartyId      string       `json:"counterPartyId,omitempty"`
	Protocol            string       `json:"protocol"`
	QuerySpec           *QuerySpec   `json:"querySpec,omitempty"`
}

type DatasetRequest struct {
	Context             *interface{} `json:"@context"`
	Type                string       `json:"@type,omitempty"`
	Id                  string       `json:"@id"`
	CounterPartyAddress string       `json:"counterPartyAddress"`
	CounterPartyId      string       `json:"counterPartyId,omitempty"`
	Protocol            string       `json:"protocol,omitempty"`
	QuerySpec           *QuerySpec   `json:"querySpec,omitempty"`
}

type DataAddress struct {
	Context       *interface{} `json:"@context"`
	Type          string       `json:"@type,omitempty"`
	Type_         string       `json:"type,omitempty"`
	BaseUrl       string       `json:"baseUrl,omitempty"`
	Endpoint      string       `json:"endpoint,omitempty"`
	EndpointType  string       `json:"endpointType,omitempty"`
	AuthType      string       `json:"authType,omitempty"`
	Authorization string       `json:"authorization,omitempty"`
}

type TransferRequest struct {
	Context             *interface{}            `json:"@context"`
	Type_               string                  `json:"@type,omitempty"`
	AssetId             string                  `json:"assetId"`
	CallbackAddresses   CallbackAddressSelector `json:"callbackAddresses,omitempty"`
	ConnectorAddress    string                  `json:"connectorAddress,omitempty"`
	ConnectorId         string                  `json:"connectorId,omitempty"`
	ContractId          string                  `json:"contractId"`
	CounterPartyAddress string                  `json:"counterPartyAddress"`
	DataDestination     *map[string]interface{} `json:"dataDestination,omitempty"`
	PrivateProperties   map[string]string       `json:"privateProperties,omitempty"`
	Protocol            string                  `json:"protocol"`
	TransferType        string                  `json:"transferType"`
}

type TransferProcess struct {
	Id                  string                  `json:"@id,omitempty"`
	Type_               string                  `json:"@type,omitempty"`
	CallbackAddresses   CallbackAddressSelector `json:"callbackAddresses,omitempty"`
	ContractAgreementId string                  `json:"contractAgreementId,omitempty"`
	CounterPartyAddress string                  `json:"counterPartyAddress,omitempty"`
	CounterPartyId      string                  `json:"counterPartyId,omitempty"`
	DataDestination     *DataAddress            `json:"dataDestination,omitempty"`
	ErrorDetail         string                  `json:"errorDetail,omitempty"`
	PrivateProperties   map[string]string       `json:"privateProperties,omitempty"`
	Protocol            string                  `json:"protocol,omitempty"`
	State               string                  `json:"state,omitempty"`
	StateTimestamp      int64                   `json:"stateTimestamp,omitempty"`
	Type                string                  `json:"type,omitempty"`
	CorrelationId       string                  `json:"correlationId,omitempty"`
	AssetId             string                  `json:"assetId,omitempty"`
	ContractId          string                  `json:"contractId,omitempty"`
	TransferType        string                  `json:"transferType,omitempty"`
}

type TerminateNegotiation struct {
	Context *interface{} `json:"@context"`
	Type    string       `json:"@type"`
	Id      string       `json:"@id"`
	Reason  string       `json:"reason"`
}

type TerminateTransferProcess struct {
	Context *interface{} `json:"@context"`
	Type    string       `json:"@type"`
	Id      string       `json:"@id"`
	Reason  string       `json:"reason"`
}

type EndpointDataReferenceEntry struct {
	Context               *interface{} `json:"@context"`
	Type_                 string       `json:"@type"`
	ProviderId            string       `json:"providerId"`
	AssetId               string       `json:"assetId"`
	AgreementId           string       `json:"agreementId"`
	TransferProcessId     string       `json:"transferProcessId"`
	CreatedAt             int64        `json:"createdAt"`
	ContractNegotiationId string       `json:"contractNegotiationId"`
}

type EdcAPI struct {
	client *resty.Client
	url    string
	apiKey string
}

func NewEdcAPI() *EdcAPI {
	return &EdcAPI{
		client: resty.New(),
		url:    os.Getenv("EDC_URL"),
		apiKey: os.Getenv("EDC_API_KEY"),
	}
}

func (e *EdcAPI) GetPolicies(querySpec QuerySpec) ([]PolicyDefinition, error) {
	var policies []PolicyDefinition
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&policies).
		Post(e.url + "/api/management/v3/policydefinitions/request")

	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get policies: %s", resp.String())
	}

	return policies, nil
}

func (e *EdcAPI) GetPolicy(id string) (*PolicyDefinition, error) {
	policy := PolicyDefinition{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&policy).
		Get(e.url + "/api/management/v3/policydefinitions/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get policy: %s", resp.String())
	}

	return &policy, nil
}

func (e *EdcAPI) CreatePolicy(policy PolicyDefinition) (*PolicyDefinition, error) {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(policy).
		SetResult(&policy).
		Post(e.url + "/api/management/v3/policydefinitions")
	fmt.Println(resp)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create policy: %s", resp.String())
	}

	return &policy, nil
}

func (e *EdcAPI) DeletePolicy(id string) error {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		Delete(e.url + "/api/management/v3/policydefinitions/" + id)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("unable to delete policy: %s", resp.String())
	}

	return nil
}

func (e *EdcAPI) GetAssets(querySpec QuerySpec) ([]Asset, error) {
	var assets []Asset
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&assets).
		Post(e.url + "/api/management/v3/assets/request")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get assets: %s", resp.String())
	}

	return assets, nil
}

func (e *EdcAPI) GetAsset(id string) (*Asset, error) {
	asset := Asset{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&asset).
		Get(e.url + "/api/management/v3/assets/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get asset: %s", resp.String())
	}

	return &asset, nil
}

func (e *EdcAPI) CreateAsset(asset Asset) (*Asset, error) {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(asset).
		SetResult(&asset).
		Post(e.url + "/api/management/v3/assets")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create asset: %s", resp.String())
	}

	return &asset, nil
}

func (e *EdcAPI) DeleteAsset(id string) error {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		Delete(e.url + "/api/management/v3/assets/" + id)

	fmt.Println(resp)
	fmt.Println(err)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("unable to delete asset: %s", resp.String())
	}

	return nil
}

func (e *EdcAPI) GetContractDefinitions(querySpec QuerySpec) ([]ContractDefinition, error) {
	var contractsDefinitions []ContractDefinition
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&contractsDefinitions).
		Post(e.url + "/api/management/v3/contractdefinitions/request")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contracts: %s", resp.String())
	}

	return contractsDefinitions, nil
}

func (e *EdcAPI) GetContractDefinition(id string) (*ContractDefinition, error) {
	contractDefinition := ContractDefinition{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&contractDefinition).
		Get(e.url + "/api/management/v3/contractdefinitions/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contract: %s", resp.String())
	}

	return &contractDefinition, nil
}

func (e *EdcAPI) DeleteContractDefinition(id string) error {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		Delete(e.url + "/api/management/v3/contractdefinitions/" + id)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("unable to delete contract: %s", resp.String())
	}

	return nil
}

func (e *EdcAPI) CreateContractDefinition(contractDefinition ContractDefinition) (*ContractDefinition, error) {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(contractDefinition).
		SetResult(&contractDefinition).
		Post(e.url + "/api/management/v3/contractdefinitions")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create contract: %s", resp.String())
	}

	return &contractDefinition, nil
}

func (e *EdcAPI) CreateContractNegotiation(contractNegotiation ContractRequest) (*ContractAgreement, error) {
	var contractAgreement ContractAgreement

	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(contractNegotiation).
		SetResult(&contractAgreement).
		Post(e.url + "/api/management/v3/contractnegotiations")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create contract negotiation: %s", resp.String())
	}

	return &contractAgreement, nil
}

func (e *EdcAPI) GetContractNegotiation(id string) (*ContractNegotiation, error) {
	contractNegotiation := ContractNegotiation{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&contractNegotiation).
		Get(e.url + "/api/management/v3/contractnegotiations/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contract negotiation: %s", resp.String())
	}

	return &contractNegotiation, nil
}

func (e *EdcAPI) TerminateContractNegotiation(terminateNegotiation TerminateNegotiation) error {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(terminateNegotiation).
		Post(e.url + "/api/management/v3/contractnegotiations/" + terminateNegotiation.Id + "/terminate")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("unable to terminate contract negotiation: %s", resp.String())
	}

	return nil
}

func (e *EdcAPI) GetContractAgreements(querySpec QuerySpec) ([]ContractAgreement, error) {
	var contractAgreements []ContractAgreement
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&contractAgreements).
		Post(e.url + "/api/management/v3/contractagreements/request")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contract agreements: %s", resp.String())
	}

	return contractAgreements, nil
}

func (e *EdcAPI) GetContractAgreement(id string) (*ContractAgreement, error) {
	contractAgreement := ContractAgreement{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&contractAgreement).
		Get(e.url + "/api/management/v3/contractagreements/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contract agreement: %s", resp.String())
	}

	return &contractAgreement, nil
}

func (e *EdcAPI) GetContractAgreementNegotiation(id string) (*ContractNegotiation, error) {
	contractNegotiation := ContractNegotiation{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&contractNegotiation).
		Get(e.url + "/api/management/v3/contractagreements/" + id + "/negotiation")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get contract agreement negotiation: %s", resp.String())
	}

	return &contractNegotiation, nil
}

func (e *EdcAPI) GetCatalog(catalogRequest CatalogRequest) (map[string]interface{}, error) {
	var catalog map[string]interface{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(catalogRequest).
		SetResult(&catalog).
		Post(e.url + "/api/management/v3/catalog/request")

	fmt.Println(resp)
	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get catalogs: %s", resp.String())
	}

	return catalog, nil
}

func (e *EdcAPI) GetCatalogDataset(datasetRequest DatasetRequest) (map[string]interface{}, error) {
	var dataset map[string]interface{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(datasetRequest).
		SetResult(&dataset).
		Post(e.url + "/api/management/v3/catalog/dataset/request")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get dataset: %s", resp.String())
	}

	return dataset, nil
}

func (e *EdcAPI) GetTransferProcesses(querySpec QuerySpec) ([]TransferProcess, error) {
	var transferProcesses []TransferProcess
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(querySpec).
		SetResult(&transferProcesses).
		Post(e.url + "/api/management/v3/transferprocesses/request")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get transfer processes: %s", resp.String())
	}

	return transferProcesses, nil
}

func (e *EdcAPI) GetTransferProcess(id string) (*TransferProcess, error) {
	transferProcess := TransferProcess{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&transferProcess).
		Get(e.url + "/api/management/v3/transferprocesses/" + id)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get transfer process: %s", resp.String())
	}

	return &transferProcess, nil
}

func (e *EdcAPI) CreateTransferProcess(transferRequest TransferRequest) (*TransferProcess, error) {
	var transferProcess TransferProcess

	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(transferRequest).
		SetResult(&transferProcess).
		Post(e.url + "/api/management/v3/transferprocesses")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to create transfer process: %s", resp.String())
	}

	return &transferProcess, nil
}

func (e *EdcAPI) TerminateTransferProcess(terminateTransferProcess TerminateTransferProcess) error {
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetBody(terminateTransferProcess).
		Post(e.url + "/api/management/v3/transferprocesses/" + terminateTransferProcess.Id + "/terminate")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("unable to terminate transfer process: %s", resp.String())
	}

	return nil
}

func (e *EdcAPI) GetEdrDataAddress(id string) (*DataAddress, error) {
	dataAddress := DataAddress{}
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", e.apiKey).
		SetResult(&dataAddress).
		Get(e.url + "/api/management/v3/edrs/" + id + "/dataaddress")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get edr data address: %s", resp.String())
	}

	return &dataAddress, nil
}
