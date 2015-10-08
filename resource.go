package makeplans

import (
	"bytes"
	"encoding/json"
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

// func (c *Client) MakePerson(p Person) error {
// 	var buf bytes.Buffer
// 	enc := json.NewEncoder(&buf)
// 	err := enc.Encode(personWrap{p})
// 	if err != nil {
// 		return err
// 	}
// 	bs, err := c.Do("POST", PersonURL, &buf)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(string(bs))
// 	return nil
// }

// func (c *Client) UpdatePerson(p Person) (Person, error) {
// 	var buf bytes.Buffer
// 	enc := json.NewEncoder(&buf)
// 	err := enc.Encode(personWrap{p})
// 	if err != nil {
// 		return Person{}, err
// 	}
// 	if p.ID == 0 {
// 		return Person{}, errors.New("ID is required")
// 	}
// 	bs, err := c.Do("PUT", PersonURL+"/"+strconv.Itoa(p.ID), &buf)
// 	if err != nil {
// 		return Person{}, err
// 	}

// 	var pw personWrap
// 	err = json.Unmarshal(bs, &pw)
// 	if err != nil {
// 		return Person{}, err
// 	}
// 	return pw.Person, err
// }

// func (c *Client) DeletePerson(p Person) error {
// 	return errors.New("not implemented https://github.com/makeplans/makeplans-api/#delete-person")
// }
