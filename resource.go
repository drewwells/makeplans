package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// "resource": {
//            "capacity": 1,
//            "created_at": "2012-09-20T15:34:16+02:00",
//            "id": 1,
//            "opening_hours_mon": ["08:00", "16:00"],
//            "opening_hours_tue": ["08:00", "11:00", "13:00", "17:30"],
//            "opening_hours_wed": ["08:00", "16:00"],
//            "opening_hours_thu": ["08:00", "12:00", "14:00", "20:00"],
//            "opening_hours_fri": ["08:00", "12:00", "12:30", "17:30"],
//            "opening_hours_sat": null,
//            "opening_hours_sun": null,
//            "title": "Mr. Spine Twister",
//            "updated_at": "2012-09-20T15:34:16+02:00"
//        }

// swagger:model resource
type Resource struct {
	ID       int    `json:"id"`
	Capacity int    `json:"capacity"`
	Title    string `json:"title"`

	OpeningHoursMon []string `json:"opening_hours_mon"`
	OpeningHoursTue []string `json:"opening_hours_tue"`
	OpeningHoursWed []string `json:"opening_hours_wed"`
	OpeningHoursThu []string `json:"opening_hours_thu"`
	OpeningHoursFri []string `json:"opening_hours_fri"`
	OpeningHoursSat []string `json:"opening_hours_sat"`
	OpeningHoursSun []string `json:"opening_hours_sun"`

	Services   []Service   `json:"services"`
	CustomData interface{} `json:"custom_data"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type resourceWrap struct {
	Resource Resource `json:"resource"`
}

var ResourceURL = "/resources"

func (c *Client) Resources() ([]Resource, error) {
	bs, err := c.Do("GET", ResourceURL+"/", nil)
	if err != nil {
		return nil, err
	}

	var wp []resourceWrap
	err = json.Unmarshal(bs, &wp)
	if err != nil {
		return nil, err
	}
	ress := make([]Resource, len(wp))
	for i, w := range wp {
		ress[i] = w.Resource
	}
	return ress, err
}

func (c *Client) Resource(id int) (Resource, error) {
	bs, err := c.Do("GET", ResourceURL+"/"+strconv.Itoa(id), nil)
	if err != nil {
		return Resource{}, err
	}

	var wp resourceWrap
	err = json.Unmarshal(bs, &wp)
	if err != nil {
		return Resource{}, err
	}

	return wp.Resource, err
}

func (c *Client) ResourceUpdate(r Resource) (Resource, error) {
	var ret resourceWrap
	if r.ID == 0 {
		return ret.Resource, errors.New("id required")
	}
	var buf bytes.Buffer
	req := resourceWrap{
		Resource: r,
	}
	enc := json.NewEncoder(&buf)
	enc.Encode(req)
	u := ResourceURL + "/" + strconv.Itoa(r.ID)
	bs, err := c.Do("PUT", u, &buf)
	if err != nil {
		return ret.Resource, err
	}
	err = json.Unmarshal(bs, &ret)
	return ret.Resource, err
}

func (c *Client) ResourceDelete(id int) (r Resource, err error) {
	bs, err := c.Do("DELETE", ResourceURL+"/"+strconv.Itoa(id), nil)
	if err != nil {
		return
	}
	var wr resourceWrap
	err = json.Unmarshal(bs, &wr)
	r = wr.Resource
	return
}

func (c *Client) MakeResource(r Resource) (Resource, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(resourceWrap{r})
	if err != nil {
		return Resource{}, err
	}
	bs, err := c.Do("POST", ResourceURL, &buf)
	if err != nil {
		return Resource{}, err
	}
	var wr resourceWrap
	err = json.Unmarshal(bs, &wr)
	if err != nil {
		return Resource{}, err
	}
	return wr.Resource, nil
}
