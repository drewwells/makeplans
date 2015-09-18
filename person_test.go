package makeplans

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestPerson_list(t *testing.T) {
	client := New(ac.Name, ac.Token)
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

var testPerson = []byte(`[
    {
        "person": {
            "created_at": "2012-09-20T15:34:16+02:00",
            "custom_data": {},
            "date_of_birth": null,
            "email": "espen@makeplans.no",
            "external_id": null,
            "id": 1,
            "name": "Espen Antonsen",
            "national_id_no": null,
            "notes": null,
            "phonenumber": "",
            "updated_at": "2012-09-20T15:34:16+02:00",
            "phone_number_formatted": null
        }
    }
]`)

func getTestPerson() Person {
	var tPpl []personWrap
	err := json.Unmarshal(testPerson, &tPpl)
	if err != nil {
		log.Fatal("Failed to decode testPerson")
	}
	p := tPpl[0].Person
	return p
}

func TestPerson_make(t *testing.T) {
	client := New(ac.Name, ac.Token)
	p := getTestPerson()
	fmt.Printf("% #v\n", p)
	err := client.MakePerson(p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerson_update(t *testing.T) {
	client := New(ac.Name, ac.Token)
	p := getTestPerson()
	p.ID = 12316
	pp, err := client.UpdatePerson(p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("% #v\n", pp)
}
