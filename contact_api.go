package intercom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// ContactRepository defines the interface for working with Contacts through the API.
type ContactRepository interface {
	find(context.Context, UserIdentifiers) (Contact, error)
	list(context.Context, contactListParams) (ContactList, error)
	scroll(context.Context, string) (ContactList, error)
	create(context.Context, *Contact) (Contact, error)
	update(context.Context, *Contact) (Contact, error)
	convert(context.Context, *Contact, *User) (User, error)
	delete(context.Context, string) (Contact, error)
}

// ContactAPI implements ContactRepository
type ContactAPI struct {
	httpClient interfaces.HTTPClient
}

func (api ContactAPI) find(ctx context.Context, params UserIdentifiers) (Contact, error) {
	return unmarshalToContact(api.getClientForFind(ctx, params))
}

func (api ContactAPI) getClientForFind(ctx context.Context, params UserIdentifiers) ([]byte, error) {
	switch {
	case params.ID != "":
		return api.httpClient.Get(ctx, fmt.Sprintf("/contacts/%s", params.ID), nil)
	case params.UserID != "":
		return api.httpClient.Get(ctx, "/contacts", params)
	}
	return nil, errors.New("Missing Contact Identifier")
}

func (api ContactAPI) list(ctx context.Context, params contactListParams) (ContactList, error) {
	contactList := ContactList{}
	data, err := api.httpClient.Get(ctx, "/contacts", params)
	if err != nil {
		return contactList, err
	}
	err = json.Unmarshal(data, &contactList)
	return contactList, err
}

func (api ContactAPI) scroll(ctx context.Context, scrollParam string) (ContactList, error) {
	contactList := ContactList{}
	params := scrollParams{ScrollParam: scrollParam}
	data, err := api.httpClient.Get(ctx, "/contacts/scroll", params)
	if err != nil {
		return contactList, err
	}
	err = json.Unmarshal(data, &contactList)
	return contactList, err
}

func (api ContactAPI) create(ctx context.Context, contact *Contact) (Contact, error) {
	requestContact := api.buildRequestContact(contact)
	return unmarshalToContact(api.httpClient.Post(ctx, "/contacts", &requestContact))
}

func (api ContactAPI) update(ctx context.Context, contact *Contact) (Contact, error) {
	requestContact := api.buildRequestContact(contact)
	return unmarshalToContact(api.httpClient.Post(ctx, "/contacts", &requestContact))
}

func (api ContactAPI) convert(ctx context.Context, contact *Contact, user *User) (User, error) {
	cr := convertRequest{Contact: api.buildRequestContact(contact), User: requestUser{
		ID:         user.ID,
		UserID:     user.UserID,
		Email:      user.Email,
		SignedUpAt: user.SignedUpAt,
	}}
	return unmarshalToUser(api.httpClient.Post(ctx, "/contacts/convert", &cr))
}

func (api ContactAPI) delete(ctx context.Context, id string) (Contact, error) {
	contact := Contact{}
	data, err := api.httpClient.Delete(ctx, fmt.Sprintf("/contacts/%s", id), nil)
	if err != nil {
		return contact, err
	}
	err = json.Unmarshal(data, &contact)
	return contact, err
}

type convertRequest struct {
	User    requestUser `json:"user"`
	Contact requestUser `json:"contact"`
}

func unmarshalToContact(data []byte, err error) (Contact, error) {
	savedContact := Contact{}
	if err != nil {
		return savedContact, err
	}
	err = json.Unmarshal(data, &savedContact)
	return savedContact, err
}

func (api ContactAPI) buildRequestContact(contact *Contact) requestUser {
	return requestUser{
		ID:                     contact.ID,
		Email:                  contact.Email,
		Phone:                  contact.Phone,
		UserID:                 contact.UserID,
		Name:                   contact.Name,
		LastRequestAt:          contact.LastRequestAt,
		LastSeenIP:             contact.LastSeenIP,
		UnsubscribedFromEmails: contact.UnsubscribedFromEmails,
		Companies:              api.getCompaniesToSendFromContact(contact),
		CustomAttributes:       contact.CustomAttributes,
		UpdateLastRequestAt:    contact.UpdateLastRequestAt,
		NewSession:             contact.NewSession,
	}
}

func (api ContactAPI) getCompaniesToSendFromContact(contact *Contact) []UserCompany {
	if contact.Companies == nil {
		return []UserCompany{}
	}
	return RequestUserMapper{}.MakeUserCompaniesFromCompanies(contact.Companies.Companies)
}
