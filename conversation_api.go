package intercom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// ConversationRepository defines the interface for working with Conversations through the API.
type ConversationRepository interface {
	find(context.Context, string) (Conversation, error)
	list(context.Context, ConversationListParams) (ConversationList, error)
	read(context.Context, string) (Conversation, error)
	reply(context.Context, string, *Reply) (Conversation, error)
}

// ConversationAPI implements ConversationRepository
type ConversationAPI struct {
	httpClient interfaces.HTTPClient
}

type conversationReadRequest struct {
	Read bool `json:"read"`
}

func (api ConversationAPI) list(ctx context.Context, params ConversationListParams) (ConversationList, error) {
	convoList := ConversationList{}
	data, err := api.httpClient.Get(ctx, "/conversations", params)
	if err != nil {
		return convoList, err
	}
	err = json.Unmarshal(data, &convoList)
	return convoList, err
}

func (api ConversationAPI) read(ctx context.Context, id string) (Conversation, error) {
	conversation := Conversation{}
	data, err := api.httpClient.Post(ctx, fmt.Sprintf("/conversations/%s", id), conversationReadRequest{Read: true})
	if err != nil {
		return conversation, err
	}
	err = json.Unmarshal(data, &conversation)
	return conversation, err
}

func (api ConversationAPI) reply(ctx context.Context, id string, reply *Reply) (Conversation, error) {
	conversation := Conversation{}
	data, err := api.httpClient.Post(ctx, fmt.Sprintf("/conversations/%s/reply", id), reply)
	if err != nil {
		return conversation, err
	}
	err = json.Unmarshal(data, &conversation)
	return conversation, nil
}

func (api ConversationAPI) find(ctx context.Context, id string) (Conversation, error) {
	conversation := Conversation{}
	data, err := api.httpClient.Get(ctx, fmt.Sprintf("/conversations/%s", id), nil)
	if err != nil {
		return conversation, err
	}
	err = json.Unmarshal(data, &conversation)
	return conversation, err
}
