package intercom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensimsim/intercom-go/interfaces"
)

// JobRepository defines the interface for working with Jobs.
type JobRepository interface {
	save(context.Context, *JobRequest) (JobResponse, error)
	find(context.Context, string) (JobResponse, error)
}

// JobAPI implements TagRepository
type JobAPI struct {
	httpClient interfaces.HTTPClient
}

func (api JobAPI) save(ctx context.Context, job *JobRequest) (JobResponse, error) {
	for i := range job.Items {
		obj := job.Items[i].Data
		switch obj.(type) {
		case *User:
			user := obj.(*User)
			job.Items[i].Data = RequestUserMapper{}.ConvertUser(user)
		}
	}
	savedJob := JobResponse{}
	data, err := api.httpClient.Post(ctx, fmt.Sprintf("/bulk/%s", job.bulkType), job)
	if err != nil {
		return savedJob, err
	}
	err = json.Unmarshal(data, &savedJob)
	return savedJob, err
}

func (api JobAPI) find(ctx context.Context, id string) (JobResponse, error) {
	fetchedJob := JobResponse{}
	data, err := api.httpClient.Get(ctx, fmt.Sprintf("/jobs/%s", id), nil)
	if err != nil {
		return fetchedJob, err
	}
	err = json.Unmarshal(data, &fetchedJob)
	return fetchedJob, err
}
