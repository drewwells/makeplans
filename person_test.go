package makeplans

import (
	"encoding/json"
	"log"
	"testing"
)

var peopleResponse = []byte(`[
  {
    "person": {
      "city": null,
      "country_code": null,
      "created_at": "2015-09-15T19:11:38-05:00",
      "custom_data": {},
      "date_of_birth": null,
      "email": "espen@makeplans.no",
      "external_id": null,
      "id": 12316,
      "name": "Espen Antonsen",
      "national_id_no": null,
      "notes": null,
      "phonenumber": null,
      "postal_code": null,
      "state": null,
      "street": null,
      "updated_at": "2015-11-07T09:49:41-06:00",
      "phone_number": null,
      "phone_number_formatted": null
    }
  }
]`)

func TestPerson_list(t *testing.T) {
	_, client := mockServerClient(t)
	ppl, err := client.People()
	if err != nil {
		t.Fatal(err)
	}
	if len(ppl) == 0 {
		t.Fatal("expected more than 1")
	}
	if ppl[0].ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

var personResponse = []byte(`{
  "person": {
    "city": "Your City",
    "country_code": null,
    "created_at": "2015-11-07T11:21:49-06:00",
    "custom_data": {},
    "date_of_birth": null,
    "email": "test@mail.com",
    "external_id": null,
    "id": 12380,
    "name": "Chad McChad",
    "national_id_no": null,
    "notes": null,
    "phonenumber": null,
    "postal_code": null,
    "state": "Your State",
    "street": "Your Street",
    "updated_at": "2015-11-07T11:35:06-06:00",
    "phone_number": null,
    "phone_number_formatted": null
  }
}`)

func getTestPerson() Person {
	var tPpl personWrap
	err := json.Unmarshal(personResponse, &tPpl)
	if err != nil {
		log.Fatal("Failed to decode testPerson", err)
	}
	p := tPpl.Person
	return p
}

func TestPerson_make(t *testing.T) {
	_, client := mockServerClient(t)
	p := getTestPerson()

	pp, err := client.MakePerson(p)
	if err != nil {
		t.Fatal(err)
	}

	if e := 12380; e != pp.ID {
		t.Errorf("got: %d wanted: %d", pp.ID, e)
	}

	if e := "test@mail.com"; e != pp.Email {
		t.Errorf("got: %s wanted: %s", pp.Email, e)
	}
}

func TestPerson_update(t *testing.T) {

	_, client := mockServerClient(t)
	p := getTestPerson()
	p.Name = "Chad McChad"
	p.Street = "Your Street"
	p.City = "Your City"
	p.State = "Your State"
	pp, err := client.UpdatePerson(p)
	if err != nil {
		t.Fatal(err)
	}

	if e := 12380; e != pp.ID {
		t.Errorf("got: %d wanted: %d", pp.ID, e)
	}

	if e := "test@mail.com"; e != pp.Email {
		t.Errorf("got: %s wanted: %s", pp.Email, e)
	}

	if pp.Name != p.Name {
		t.Errorf("got: %s wanted: %s", pp.Name, p.Name)
	}
}
