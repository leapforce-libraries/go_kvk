package kvk

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"regexp"
)

const (
	apiName       string = "KvK"
	apiPathV1     string = "https://api.kvk.nl/api/v1"
	apiPathV2     string = "https://api.kvk.nl/api/v2"
	apiPathTest   string = "https://developers.kvk.nl/test/api/v1"
	regexPostcode string = `^[1-9]{1}[0-9]{3}[ ]{0,1}[a-zA-Z]{2}$`
)

type Service struct {
	apiKey        string
	isTest        bool
	httpService   *go_http.Service
	rPostcode     *regexp.Regexp
	errorResponse *ErrorResponse
}

type ServiceConfig struct {
	ApiKey string
	IsTest bool
}

func (service *Service) ValidatePostcode(postcode string) bool {
	return service.rPostcode.Match([]byte(postcode))
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ApiKey == "" {
		return nil, errortools.ErrorMessage("Service ApiKey not provided")
	}

	httpClient := &http.Client{}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{
		HttpClient: httpClient,
	})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiKey:      serviceConfig.ApiKey,
		isTest:      serviceConfig.IsTest,
		httpService: httpService,
		rPostcode:   regexp.MustCompile(regexPostcode),
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("apikey", service.apiKey)
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	service.errorResponse = &ErrorResponse{}
	(*requestConfig).ErrorModel = &service.errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if len(service.errorResponse.Fout) > 0 {
		e.SetMessage(service.errorResponse.Fout[0].Omschrijving)
	}

	return request, response, e
}

func (service *Service) urlV1(path string, values *url.Values) string {
	values_ := url.Values{}
	if values != nil {
		values_ = *values
	}

	apiPath_ := apiPathV1
	if service.isTest {
		apiPath_ = apiPathTest
	}

	return fmt.Sprintf("%s/%s?%s", apiPath_, path, values_.Encode())
}

func (service *Service) urlV2(path string, values *url.Values) string {
	values_ := url.Values{}
	if values != nil {
		values_ = *values
	}

	apiPath_ := apiPathV2
	if service.isTest {
		apiPath_ = apiPathTest
	}

	return fmt.Sprintf("%s/%s?%s", apiPath_, path, values_.Encode())
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}

func (service *Service) ErrorResponse() *ErrorResponse {
	return service.errorResponse
}
