package intercom

import (
	"context"
	"encoding/json"

	"github.com/opensimsim/intercom-go/interfaces"
)

// AdminRepository defines the interface for working with Admins through the API.
type AdminRepository interface {
	list(context.Context) (AdminList, error)
}

// AdminAPI implements AdminRepository
type AdminAPI struct {
	httpClient interfaces.HTTPClient
}

func (api AdminAPI) list(ctx context.Context) (AdminList, error) {
	adminList := AdminList{}
	data, err := api.httpClient.Get(ctx, "/admins", nil)
	if err != nil {
		return adminList, err
	}
	err = json.Unmarshal(data, &adminList)
	return adminList, err
}
