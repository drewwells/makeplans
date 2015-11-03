package makeplans

import (
	"encoding/json"
	"time"
)

// [
// 	{
// 		"provider": {
// 			"created_at": "2015-08-27T00:55:03-05:00",
// 			"id": 2044912747,
// 			"resource_id": 484,
// 			"service_id": 395,
// 			"updated_at": "2015-08-27T00:55:03-05:00"
// 		}
// 	},
// 	{
// 		"provider": {
// 			"created_at": "2015-08-27T00:55:03-05:00",
// 			"id": 2044912746,
// 			"resource_id": 484,
// 			"service_id": 394,
// 			"updated_at": "2015-08-27T00:55:03-05:00"
// 		}
// 	},
// 	{
// 		"provider": {
// 			"created_at": "2015-08-27T00:55:03-05:00",
// 			"id": 2044912745,
// 			"resource_id": 484,
// 			"service_id": 393,
// 			"updated_at": "2015-08-27T00:55:03-05:00"
// 		}
// 	}
// ]

type Provider struct {
	ID         int       `json:"id"`
	ServiceID  int       `json:"service_id,omitempty"`
	ResourceID int       `json:"resource_id,omitempty"`
	CreatedAt  time.Time `json:"updated_at,omitempty"`
	UpdatedAt  time.Time `json:"created_at,omitempty"`
}

type providerWrap struct {
	Provider Provider `json:"provider"`
}

var ProviderURL = "/providers"

func (c *Client) Providers() ([]Provider, error) {
	bs, err := c.Do("GET", ProviderURL+"/", nil)
	if err != nil {
		return nil, err
	}
	var wp []providerWrap
	err = json.Unmarshal(bs, &wp)
	if err != nil {
		return nil, err
	}
	ress := make([]Provider, len(wp))
	for i, w := range wp {
		ress[i] = w.Provider
	}
	return ress, err
}
