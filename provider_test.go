package makeplans

import (
	"fmt"
	"testing"
)

var testProviders = []byte(`[
	{
		"provider": {
			"created_at": "2015-08-27T00:55:03-05:00",
			"id": 2044912747,
			"resource_id": 484,
			"service_id": 395,
			"updated_at": "2015-08-27T00:55:03-05:00"
		}
	},
	{
		"provider": {
			"created_at": "2015-08-27T00:55:03-05:00",
			"id": 2044912746,
			"resource_id": 484,
			"service_id": 394,
			"updated_at": "2015-08-27T00:55:03-05:00"
		}
	},
	{
		"provider": {
			"created_at": "2015-08-27T00:55:03-05:00",
			"id": 2044912745,
			"resource_id": 484,
			"service_id": 393,
			"updated_at": "2015-08-27T00:55:03-05:00"
		}
	}
]`)

func TestProvider_list(t *testing.T) {
	_, client := mockServerClient(t)
	ress, err := client.Providers()
	if err != nil {
		t.Fatal(err)
	}
	if e := 3; len(ress) != e {
		t.Fatalf("got: %d wanted: %d", len(ress), e)
	}

	res := ress[2]
	if e := 2044912745; res.ID != e {
		t.Errorf("got: %d wanted: %d", res.ID, e)
	}

	if e := 484; res.ResourceID != e {
		t.Errorf("got: %d wanted: %d", res.ResourceID, e)
	}

	if e := 393; res.ServiceID != e {
		t.Errorf("got: %d wanted: %d", res.ServiceID, e)
	}

	if res.CreatedAt.IsZero() || res.UpdatedAt.IsZero() {
		t.Fatal("got zero creation or update time")
	}

}

var testProviderCreate = []byte(`{"provider":{"created_at":"2015-11-03T04:07:39-06:00","id":2044912816,"resource_id":503,"service_id":394,"updated_at":"2015-11-03T04:07:39-06:00"}}`)

func TestProvider_crud(t *testing.T) {
	_, client := mockServerClient(t)
	in := Provider{
		ResourceID: 503,
		ServiceID:  394,
	}
	res, err := client.MakeProvider(in)
	if err != nil {
		t.Fatal(err)
	}

	if e := 394; res.ServiceID != e {
		t.Errorf("got: %d wanted: %d", res.ServiceID, e)
	}

	if e := 503; res.ResourceID != e {
		t.Errorf("got: %d wanted: %d", res.ResourceID, e)
	}

	if e := 2044912816; res.ID != e {
		t.Errorf("got: %d wanted: %d", res.ID, e)
	}

	if res.CreatedAt.IsZero() || res.UpdatedAt.IsZero() {
		t.Fatal("timestamp nil")
	}
}

func TestProvider_update(t *testing.T) {
	client := New(ac.Name, ac.Token)
	in := Provider{
		ID:         2044912816,
		ResourceID: 484,
		ServiceID:  399,
	}
	res, err := client.ProviderUpdate(in)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("% #v\n", res)
}

var testProviderDelete = []byte(`{"provider":{"created_at":"2015-11-03T04:49:22-06:00","id":2044912817,"resource_id":484,"service_id":395,"updated_at":"2015-11-03T04:49:22-06:00"}}`)

func TestProvider_delete(t *testing.T) {
	_, client := mockServerClient(t)
	p, err := client.ProviderDelete(2044912817)
	if err != nil {
		t.Fatal(err)
	}

	if e := 2044912817; e != p.ID {
		t.Errorf("got: %d wanted: %d", p.ID, e)
	}

	if e := 484; e != p.ResourceID {
		t.Errorf("got: %d wanted: %d", p.ResourceID, e)
	}

	if e := 395; e != p.ServiceID {
		t.Errorf("got: %d wanted: %d", p.ServiceID, e)
	}
}
