package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Person struct {
	ID                int                    `json:"id"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Name              string                 `json:"name"`
	Email             string                 `json:"email"`
	PhoneNumber       string                 `json:"phonenumber"`
	PrettyPhoneNumber string                 `json:"phone_number_formatted"`
	ExternalID        string                 `json:"external_id"`
	CustomData        map[string]interface{} `json:"custom_data"`
	DOB               string                 `json:"date_of_birth"`
	NationalID        string                 `json:"national_id_no"`
	Street            string                 `json:"street"`
	City              string                 `json:"city"`
	PostalCode        string                 `json:"postal_code"`
	State             string                 `json:"state"`
	CountryCode       string                 `json:"country_code"`
	Notes             string                 `json:"notes"`
}

type personWrap struct {
	Person Person `json:"person"`
}

var PersonURL = "/people"

func (c *Client) People() ([]Person, error) {
	bs, err := c.Do("GET", PersonURL+"/", nil)
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

func (c *Client) MakePerson(p Person) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(personWrap{p})
	if err != nil {
		return err
	}
	bs, err := c.Do("POST", PersonURL, &buf)
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return nil
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
	bs, err := c.Do("PUT", PersonURL+"/"+strconv.Itoa(p.ID), &buf)
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
