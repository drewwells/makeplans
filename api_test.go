package makeplans

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type account struct {
	Name  string
	Token string
}

var ac account

var realClient *Client

func init() {
	bs, err := ioutil.ReadFile("cli/account.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bs, &ac)
	if err != nil {
		log.Fatal(err)
	}
	realClient = New(ac.Name, ac.Token)
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

func apiError(err error) []byte {
	return []byte(fmt.Sprintf(`{"error":{"description":"%s"}}`,
		err.Error()))
}

func mockServerClient(t *testing.T) (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case BookingURL + "410369":
			switch r.Method {
			case "DELETE":
				w.Write(bookingDeleteResponse)
			}
		case BookingURL + "410369" + "/cancel":
			w.Write(bookingCancelResponse)
		case PersonURL:
			switch r.Method {
			case "GET":
				w.Write(peopleResponse)
			case "POST":
				w.Write(personResponse)
			}
		case PersonURL + "12380":
			switch r.Method {
			case "PUT":
				w.Write(personResponse)
			}
		case ProvidersURL:
			if r.Method == "POST" {
				w.Write(testProviderCreate)
			} else if r.Method == "GET" {
				w.Write(testProviders)
			}
		case ProvidersURL + "2044912817":
			if r.Method == "DELETE" {
				w.Write(testProviderDelete)
			}
		case ServiceURL:
			if r.Method == "GET" {
				w.Write(testServices)
			} else if r.Method == "POST" {
				bs, _ := json.Marshal(mockServiceWrapCreate)
				w.Write(bs)
			} else if r.Method == "DELETE" {
				del := mockServiceWrapCreate
				del.Service.Active = false
				bs, _ := json.Marshal(del)
				w.Write(bs)
			}
		case "/services/320/next_available_date":
			w.Write(testSlotNext)
		case "/services/running/slots?from=2015-11-07&to=2015-11-07":
			w.Write(testSlots)
		case "/resources/501?from=2015-11-07&to=2015-08-01":
			w.Write(resourceOpeningResponse)
		case ResourceURL + "/":
			w.Write(testResources)
		case ResourceURL + "/484":
			w.Write(testResource)
		case ResourceURL + "/100":
			w.Write(apiError(ErrNotFound))
		case BookingURL:
			switch r.Method {
			case "POST":
				// Inject bad request if 400 is detected
				if fakeBookingCapacityFailure {
					w.Write([]byte(`{"error":{"description":"error resource_id: Not available for booking at this timeerror count: More than maximum count per booking"}}`))
				} else {
					w.Write(testBookingSuccess)
				}
			case "GET":
				w.Write(testBookings)
			}
		case EventsURL:
			w.Write(testEvents)
		default:
			pan := fmt.Sprintf("Not implemented %s: %s", r.Method,
				r.URL.String())
			panic(pan)
		}

	}))

	return ts, &Client{
		URL:      ts.URL,
		Resolver: testResolver,
	}
}

var testResolver = func(urlTmpl string, accountName string) string {
	return urlTmpl
}

var mockServiceWrapCreate = serviceWrap{
	Service: Service{
		Active:          true,
		BookingCapacity: 1,
		Interval:        30,
		MaxSlots:        100,
		Price:           "115.0",
		SameDay:         true,
		Title:           "Chiropractor",
		BookingTypeID:   1,
	},
}

func TestService_create(t *testing.T) {
	_, client := mockServerClient(t)

	svc := mockServiceWrapCreate.Service
	svc, err := client.ServiceCreate(svc)
	if err != nil {
		t.Fatal(err)
	}
	if e := true; svc.Active != e {
		t.Errorf("got: %t wanted: %t", svc.Active, e)
	}

	if e := 1; svc.BookingTypeID != e {
		t.Errorf("got: %d wanted: %d", svc.BookingTypeID, e)
	}
}

func TestService_delete(t *testing.T) {
	client := New(ac.Name, ac.Token)
	svc, err := client.ServiceDelete(405)
	if err != nil {
		t.Fatal(err)
	}

	if e := false; svc.Active != e {
		t.Errorf("got: %t wanted: %t", svc.Active, e)
	}
}

func TestService_list(t *testing.T) {

	_, client := mockServerClient(t)

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
	_, client := mockServerClient(t)

	slots, err := client.ServiceSlot("running", time.Now(), time.Now())
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

var testSlotNext = []byte(`[
    {
        "available_date": "2013-04-03"
    }
]`)

func TestSlot_next(t *testing.T) {
	_, client := mockServerClient(t)

	slots, err := client.SlotNextDate("320")
	if err != nil {
		t.Error(err)
	}

	if e := 1; len(slots) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(slots), e)
	}
}

var testEvents = []byte(`[
    {
        "event": {
            "capacity": 10,
            "created_at": "2012-09-20T15:34:16+02:00",
            "custom_data": {},
            "description": null,
            "end": "2015-08-10T11:30:00+02:00",
            "id": 1,
            "resource_id": 1,
            "published": true,
            "start": "2015-08-10T10:00:00+02:00",
            "service_id": 1,
            "title": "Super fun event",
            "updated_at": "2012-09-20T15:34:16+02:00"
        }
    }
]`)

func TestEvent_list(t *testing.T) {
	_, client := mockServerClient(t)

	evts, err := client.Events()
	if err != nil {
		t.Error(err)
	}

	if e := 1; len(evts) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(evts), e)
	}

	evt := evts[0]
	if e := 1; evt.ID != e {
		t.Fatalf("got: %d wanted: %d", evt.ID, e)
	}
	if e := 10; evt.Capacity != e {
		t.Fatalf("got: %d wanted: %d", evt.Capacity, e)
	}

}
