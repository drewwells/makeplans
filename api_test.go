package makeplans

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type account struct {
	Name  string
	Token string
}

var testServices = []byte(`[
  {
    "service": {
      "active": true,
      "booking_capacity": 1,
      "booking_type_id": 1,
      "created_at": "2015-08-23T18:59:51-05:00",
      "custom_data": {},
      "day_booking_specify_time": null,
      "description": "<p>This is a description of the service.</p>",
      "has_day_booking": false,
      "id": 393,
      "interval": 60,
      "mail_confirmation": null,
      "mail_verification": null,
      "max_slots": 1,
      "price": "20.0",
      "same_day": false,
      "sms_confirmation": null,
      "sms_reminder": null,
      "sms_verification": null,
      "template": null,
      "title": "Cross Fit Session",
      "updated_at": "2015-08-24T22:28:17-05:00"
    }
  },
  {
    "service": {
      "active": true,
      "booking_capacity": 1,
      "booking_type_id": 2,
      "created_at": "2015-08-23T19:00:07-05:00",
      "custom_data": {},
      "day_booking_specify_time": null,
      "description": "<p>This is a cross fit type</p>",
      "has_day_booking": false,
      "id": 394,
      "interval": 60,
      "mail_confirmation": null,
      "mail_verification": null,
      "max_slots": 1,
      "price": null,
      "same_day": false,
      "sms_confirmation": null,
      "sms_reminder": null,
      "sms_verification": null,
      "template": null,
      "title": "Cross Fit Type",
      "updated_at": "2015-08-23T19:01:34-05:00"
    }
  },
  {
    "service": {
      "active": true,
      "booking_capacity": 10,
      "booking_type_id": 1,
      "created_at": "2015-08-23T19:37:47-05:00",
      "custom_data": {},
      "day_booking_specify_time": null,
      "description": null,
      "has_day_booking": false,
      "id": 395,
      "interval": 30,
      "mail_confirmation": null,
      "mail_verification": null,
      "max_slots": 1,
      "price": null,
      "same_day": false,
      "sms_confirmation": null,
      "sms_reminder": null,
      "sms_verification": null,
      "template": null,
      "title": "Running Session",
      "updated_at": "2015-08-23T19:37:47-05:00"
    }
  }
]`)

func mockServer(t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/services":
			w.Write(testServices)
		case "/running/slots":
			w.Write(testSlots)
		}
	}))

	return ts
}

func init() {
	tokenURL = func(urlTmpl string, accountName string) string {
		return urlTmpl
	}
}

func TestService_list(t *testing.T) {

	ts := mockServer(t)

	client := Client{
		URL: ts.URL,
	}

	svcs, err := client.Services()
	if err != nil {
		t.Error(err)
	}

	if e := 3; len(svcs) != e {
		t.Fatalf("wrong number of services returned got: %d wanted: %d",
			len(svcs), e)
	}

	if e := "20.0"; svcs[0].Price != e {
		t.Fatalf("got: %s wanted: %s", svcs[0].Price, e)
	}

	if e := "Cross Fit Session"; e != svcs[0].Title {
		t.Fatalf("got: %s wanted: %s", svcs[0].Title, e)
	}

}

var testSlots = []byte(`[
    {
        "slot": {
            "timestamp": "2013-03-08T10:00:00+00:00",
            "timestamp_end": "2013-03-08T10:15:00+00:00",
            "formatted_timestamp": "Friday, March  8, 2013, 10:00 AM",
            "formatted_timestamp_end": "Friday, March  8, 2013, 10:15 AM",
            "free": 1,
            "open_resources": [
                1
            ],
            "available_resources": [
                1,2
            ]
        }
    }
]`)

func TestSlot_list(t *testing.T) {
	ts := mockServer(t)

	client := Client{
		URL: ts.URL,
	}

	slots, err := client.Slots("running")
	if err != nil {
		t.Error(err)
	}

	if e := 1; len(slots) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(slots), e)
	}

	if e := 1; slots[0].Free != e {
		t.Fatalf("got: %d wanted: %d", slots[0].Free, e)
	}

	if e := "Friday, March  8, 2013, 10:00 AM"; e !=
		slots[0].FormattedTimestamp {
		t.Fatalf("got: %s wanted: %s", slots[0].FormattedTimestamp, e)
	}

}
