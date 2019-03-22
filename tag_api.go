package intercom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// TagRepository defines the interface for working with Tags through the API.
type TagRepository interface {
	list(context.Context) (TagList, error)
	save(context.Context, *Tag) (Tag, error)
	delete(context.Context, string) error
	tag(context.Context, *TaggingList) (Tag, error)
}

// TagAPI implements TagRepository
type TagAPI struct {
	httpClient interfaces.HTTPClient
}

func (api TagAPI) list(ctx context.Context) (TagList, error) {
	tagList := TagList{}
	data, err := api.httpClient.Get(ctx, "/tags", nil)
	if err != nil {
		return tagList, err
	}
	err = json.Unmarshal(data, &tagList)
	return tagList, err
}

func (api TagAPI) save(ctx context.Context, tag *Tag) (Tag, error) {
	savedTag := Tag{}
	data, err := api.httpClient.Post(ctx, "/tags", tag)
	if err != nil {
		return savedTag, err
	}
	err = json.Unmarshal(data, &savedTag)
	return savedTag, err
}

func (api TagAPI) delete(ctx context.Context, id string) error {
	_, err := api.httpClient.Delete(ctx, fmt.Sprintf("/tags/%s", id), nil)
	return err
}

func (api TagAPI) tag(ctx context.Context, taggingList *TaggingList) (Tag, error) {
	savedTag := Tag{}
	data, err := api.httpClient.Post(ctx, "/tags", taggingList)
	if err != nil {
		return savedTag, err
	}
	err = json.Unmarshal(data, &savedTag)
	return savedTag, err
}
