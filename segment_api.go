package intercom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// SegmentRepository defines the interface for working with Segments through the API.
type SegmentRepository interface {
	list(context.Context) (SegmentList, error)
	find(context.Context, string) (Segment, error)
}

// SegmentAPI implements SegmentRepository
type SegmentAPI struct {
	httpClient interfaces.HTTPClient
}

func (api SegmentAPI) list(ctx context.Context) (SegmentList, error) {
	segmentList := SegmentList{}
	data, err := api.httpClient.Get(ctx, "/segments", nil)
	if err != nil {
		return segmentList, err
	}
	err = json.Unmarshal(data, &segmentList)
	return segmentList, err
}

func (api SegmentAPI) find(ctx context.Context, id string) (Segment, error) {
	segment := Segment{}
	data, err := api.httpClient.Get(ctx, fmt.Sprintf("/segments/%s", id), nil)
	if err != nil {
		return segment, err
	}
	err = json.Unmarshal(data, &segment)
	return segment, err
}
