package intercom

import (
	"context"

	"github.com/opensimsim/intercom-go/interfaces"
)

// EventRepository defines the interface for working with Events through the API.
type EventRepository interface {
	save(context.Context, *Event) error
}

// EventAPI implements EventRepository
type EventAPI struct {
	httpClient interfaces.HTTPClient
}

func (api EventAPI) save(ctx context.Context, event *Event) error {
	_, err := api.httpClient.Post(ctx, "/events", event)
	return err
}
