package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type ServiceGetter interface {
	Service() ServiceInterface
}

type ServiceInterface interface {
	List(params url.Values) (*adminv1.ServiceList, error)
	Get(name string) (*adminv1.Service, *APIResponse)
	Create(service *adminv1.Service) (*adminv1.Service, *APIResponse)
	Delete(name string) error
}

type serviceAPI struct {
	client APIInterface
}

func (a *serviceAPI) Create(service *adminv1.Service) (*adminv1.Service, *APIResponse) {
	out := &adminv1.Service{}
	err := a.client.Create(service, out)
	return out, err
}

func (a *serviceAPI) Get(name string) (*adminv1.Service, *APIResponse) {
	out := &adminv1.Service{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *serviceAPI) List(params url.Values) (*adminv1.ServiceList, error) {
	ServiceList := &adminv1.ServiceList{}
	request := a.client.RestClient().Get().Resource("services")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, ServiceList); err != nil {
		return nil, err
	}

	if len(ServiceList.NextPage) > 0 {
		params.Add("offset", ServiceList.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		ServiceList.Items = append(ServiceList.Items, result.Items...)
	}

	return ServiceList, err
}

func (a *serviceAPI) Delete(id string) error {
	return a.client.Delete(id)
}
