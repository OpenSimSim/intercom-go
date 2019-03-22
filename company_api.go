package intercom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// CompanyRepository defines the interface for working with Companies through the API.
type CompanyRepository interface {
	find(context.Context, CompanyIdentifiers) (Company, error)
	list(context.Context, companyListParams) (CompanyList, error)
	listUsers(context.Context, string, companyUserListParams) (UserList, error)
	scroll(context.Context, string) (CompanyList, error)
	save(context.Context, *Company) (Company, error)
}

// CompanyAPI implements CompanyRepository
type CompanyAPI struct {
	httpClient interfaces.HTTPClient
}

type requestCompany struct {
	ID               string                 `json:"id,omitempty"`
	CompanyID        string                 `json:"company_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	RemoteCreatedAt  int64                  `json:"remote_created_at,omitempty"`
	MonthlySpend     int64                  `json:"monthly_spend,omitempty"`
	Plan             string                 `json:"plan,omitempty"`
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
}

func (api CompanyAPI) find(ctx context.Context, params CompanyIdentifiers) (Company, error) {
	company := Company{}
	data, err := api.getClientForFind(ctx, params)
	if err != nil {
		return company, err
	}
	err = json.Unmarshal(data, &company)
	return company, err
}

func (api CompanyAPI) getClientForFind(ctx context.Context, params CompanyIdentifiers) ([]byte, error) {
	switch {
	case params.ID != "":
		return api.httpClient.Get(ctx, fmt.Sprintf("/companies/%s", params.ID), nil)
	case params.CompanyID != "", params.Name != "":
		return api.httpClient.Get(ctx, "/companies", params)
	}
	return nil, errors.New("Missing Company Identifier")
}

func (api CompanyAPI) list(ctx context.Context, params companyListParams) (CompanyList, error) {
	companyList := CompanyList{}
	data, err := api.httpClient.Get(ctx, "/companies", params)
	if err != nil {
		return companyList, err
	}
	err = json.Unmarshal(data, &companyList)
	return companyList, err
}

func (api CompanyAPI) listUsers(ctx context.Context, id string, params companyUserListParams) (UserList, error) {
	companyUserList := UserList{}
	data, err := api.getClientForListUsers(ctx, id, params)
	if err != nil {
		return companyUserList, err
	}
	err = json.Unmarshal(data, &companyUserList)
	return companyUserList, err
}

func (api CompanyAPI) getClientForListUsers(ctx context.Context, id string, params companyUserListParams) ([]byte, error) {
	switch {
	case id != "":
		return api.httpClient.Get(ctx, fmt.Sprintf("/companies/%s/users", id), params)
	case params.CompanyID != "", params.Type == "user":
		return api.httpClient.Get(ctx, "/companies", params)
	}
	return nil, errors.New("Missing Company Identifier")
}

func (api CompanyAPI) scroll(ctx context.Context, scrollParam string) (CompanyList, error) {
	companyList := CompanyList{}
	params := scrollParams{ScrollParam: scrollParam}
	data, err := api.httpClient.Get(ctx, "/companies/scroll", params)
	if err != nil {
		return companyList, err
	}
	err = json.Unmarshal(data, &companyList)
	return companyList, err
}

func (api CompanyAPI) save(ctx context.Context, company *Company) (Company, error) {
	requestCompany := requestCompany{
		ID:               company.ID,
		Name:             company.Name,
		CompanyID:        company.CompanyID,
		RemoteCreatedAt:  company.RemoteCreatedAt,
		MonthlySpend:     company.MonthlySpend,
		Plan:             api.getPlanName(company),
		CustomAttributes: company.CustomAttributes,
	}

	savedCompany := Company{}
	data, err := api.httpClient.Post(ctx, "/companies", &requestCompany)
	if err != nil {
		return savedCompany, err
	}
	err = json.Unmarshal(data, &savedCompany)
	return savedCompany, err
}

func (api CompanyAPI) getPlanName(company *Company) string {
	if company.Plan == nil {
		return ""
	}
	return company.Plan.Name
}
