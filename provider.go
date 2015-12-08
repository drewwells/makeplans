package makeplans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
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
	ID         int        `json:"id,omitempty"`
	ServiceID  int        `json:"service_id,omitempty"`
	ResourceID int        `json:"resource_id,omitempty"`
	CreatedAt  *time.Time `json:"updated_at,omitempty"`
	UpdatedAt  *time.Time `json:"created_at,omitempty"`
}

type providerWrap struct {
	Provider Provider `json:"provider"`
}

var ProvidersURL = "/providers/"

func (c *Client) Providers() ([]Provider, error) {
	bs, err := c.Do("GET", ProvidersURL, nil)
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

func (c *Client) MakeProvider(in Provider) (p Provider, err error) {
	bs, err := json.Marshal(providerWrap{Provider: in})
	if err != nil {
		return
	}
	fmt.Println("sending", string(bs))
	buf := bytes.NewBuffer(bs)
	bs, err = c.Do("POST", ProvidersURL, buf)
	if err != nil {
		return
	}
	var wp providerWrap
	err = json.Unmarshal(bs, &wp)
	p = wp.Provider
	return
}

func (c *Client) ProviderUpdate(in Provider) (p Provider, err error) {
	sid := strconv.Itoa(in.ID)
	in.ID = 0
	bs, err := json.Marshal(providerWrap{Provider: in})
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(bs)
	fmt.Println("put", ProvidersURL+sid)
	fmt.Println(string(bs))
	bs, err = c.Do("PUT", ProvidersURL+sid, buf)
	if err != nil {
		return
	}
	fmt.Println("found", string(bs))
	var pw providerWrap
	err = json.Unmarshal(bs, &pw)
	p = pw.Provider
	return
}

func (c *Client) ProviderDelete(id int) (p Provider, err error) {
	sid := strconv.Itoa(id)
	bs, err := c.Do("DELETE", ProvidersURL+sid, nil)
	if err != nil {
		return
	}
	var pw providerWrap
	err = json.Unmarshal(bs, &pw)
	p = pw.Provider
	return
}
