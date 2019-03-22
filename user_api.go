package intercom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// UserRepository defines the interface for working with Users through the API.
type UserRepository interface {
	find(context.Context, UserIdentifiers) (User, error)
	list(context.Context, userListParams) (UserList, error)
	scroll(context.Context, string) (UserList, error)
	save(context.Context, *User) (User, error)
	delete(context.Context, string) (User, error)
}

// UserAPI implements UserRepository
type UserAPI struct {
	httpClient interfaces.HTTPClient
}

type requestScroll struct {
	ScrollParam string `json:"scroll_param,omitempty"`
}
type requestUser struct {
	ID                     string                 `json:"id,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	SignedUpAt             int64                  `json:"signed_up_at,omitempty"`
	RemoteCreatedAt        int64                  `json:"remote_created_at,omitempty"`
	LastRequestAt          int64                  `json:"last_request_at,omitempty"`
	LastSeenIP             string                 `json:"last_seen_ip,omitempty"`
	UnsubscribedFromEmails *bool                  `json:"unsubscribed_from_emails,omitempty"`
	Companies              []UserCompany          `json:"companies,omitempty"`
	CustomAttributes       map[string]interface{} `json:"custom_attributes,omitempty"`
	UpdateLastRequestAt    *bool                  `json:"update_last_request_at,omitempty"`
	NewSession             *bool                  `json:"new_session,omitempty"`
	LastSeenUserAgent      string                 `json:"last_seen_user_agent,omitempty"`
}

func (api UserAPI) find(ctx context.Context, params UserIdentifiers) (User, error) {
	return unmarshalToUser(api.getClientForFind(ctx, params))
}

func (api UserAPI) getClientForFind(ctx context.Context, params UserIdentifiers) ([]byte, error) {
	switch {
	case params.ID != "":
		return api.httpClient.Get(ctx, fmt.Sprintf("/users/%s", params.ID), nil)
	case params.UserID != "", params.Email != "":
		return api.httpClient.Get(ctx, "/users", params)
	}
	return nil, errors.New("Missing User Identifier")
}

func (api UserAPI) list(ctx context.Context, params userListParams) (UserList, error) {
	userList := UserList{}
	data, err := api.httpClient.Get(ctx, "/users", params)
	if err != nil {
		return userList, err
	}
	err = json.Unmarshal(data, &userList)
	return userList, err
}

func (api UserAPI) scroll(ctx context.Context, scrollParam string) (UserList, error) {
	userList := UserList{}

	url := "/users/scroll"
	params := scrollParams{ScrollParam: scrollParam}
	data, err := api.httpClient.Get(ctx, url, params)

	if err != nil {
		return userList, err
	}
	err = json.Unmarshal(data, &userList)
	return userList, err
}

func (api UserAPI) save(ctx context.Context, user *User) (User, error) {
	return unmarshalToUser(api.httpClient.Post(ctx, "/users", RequestUserMapper{}.ConvertUser(user)))
}

func unmarshalToUser(data []byte, err error) (User, error) {
	savedUser := User{}
	if err != nil {
		return savedUser, err
	}
	err = json.Unmarshal(data, &savedUser)
	return savedUser, err
}

func (api UserAPI) delete(ctx context.Context, id string) (User, error) {
	user := User{}
	data, err := api.httpClient.Delete(ctx, fmt.Sprintf("/users/%s", id), nil)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	return user, err
}
