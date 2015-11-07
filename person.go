package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Person struct {
	ID                int                    `json:"id,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Email             string                 `json:"email,omitempty"`
	PhoneNumber       string                 `json:"phonenumber,omitempty"`
	PrettyPhoneNumber string                 `json:"phone_number_formatted,omitempty"`
	ExternalID        string                 `json:"external_id,omitempty"`
	CustomData        map[string]interface{} `json:"custom_data,omitempty"`
	DOB               string                 `json:"date_of_birth,omitempty,omitempty"`
	NationalID        string                 `json:"national_id_no,omitempty"`
	Street            string                 `json:"street,omitempty"`
	City              string                 `json:"city,omitempty"`
	PostalCode        string                 `json:"postal_code,omitempty"`
	State             string                 `json:"state,omitempty"`
	CountryCode       string                 `json:"country_code,omitempty"`
	Notes             string                 `json:"notes,omitempty"`
}

type personWrap struct {
	Person Person `json:"person"`
}

// PersonURL is the base path for people endpoints
var PersonURL = "/people/"

func (c *Client) People() ([]Person, error) {
	bs, err := c.Do("GET", PersonURL, nil)
	if err != nil {
		return nil, err
	}
	var wp []personWrap
	err = json.Unmarshal(bs, &wp)
	if err != nil {
		return nil, err
	}
	ppl := make([]Person, len(wp))
	for i, w := range wp {
		ppl[i] = w.Person
	}
	return ppl, err
}

func (c *Client) MakePerson(p Person) (ret Person, err error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(personWrap{p})
	if err != nil {
		return
	}
	bs, err := c.Do("POST", PersonURL, &buf)
	if err != nil {
		return
	}

	var wrap personWrap
	err = json.Unmarshal(bs, &wrap)
	ret = wrap.Person
	return
}

func (c *Client) UpdatePerson(p Person) (Person, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(personWrap{p})
	if err != nil {
		return Person{}, err
	}
	if p.ID == 0 {
		return Person{}, errors.New("ID is required")
	}
	bs, err := c.Do("PUT", PersonURL+strconv.Itoa(p.ID), &buf)
	if err != nil {
		return Person{}, err
	}

	var pw personWrap
	err = json.Unmarshal(bs, &pw)
	if err != nil {
		return Person{}, err
	}
	return pw.Person, err
}

func (c *Client) DeletePerson(p Person) error {
	return errors.New("not implemented https://github.com/makeplans/makeplans-api/#delete-person")
}
