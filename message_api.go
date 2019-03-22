package intercom

import (
	"context"
	"encoding/json"

	"github.com/opensimsim/intercom-go/interfaces"
)

// MessageRepository defines the interface for creating and updating Messages through the API.
type MessageRepository interface {
	save(context.Context, *MessageRequest) (MessageResponse, error)
}

// MessageAPI implements MessageRepository
type MessageAPI struct {
	httpClient interfaces.HTTPClient
}

func (api MessageAPI) save(ctx context.Context, message *MessageRequest) (MessageResponse, error) {
	data, err := api.httpClient.Post(ctx, "/messages", message)
	savedMessage := MessageResponse{}
	if err != nil {
		return savedMessage, err
	}
	err = json.Unmarshal(data, &savedMessage)
	return savedMessage, err
}
