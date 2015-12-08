package makeplans

import (
	"testing"
	"time"
)

var testBookings = []byte(`[
    {
        "booking": {
            "booked_from": "2012-09-29T07:00:00+02:00",
            "booked_to": "2012-09-29T08:00:00+02:00",
            "created_at": "2012-09-20T15:34:16+02:00",
            "custom_data": {},
            "count": 1,
            "expires_at": null,
            "external_id": null,
            "id": 1,
            "notes": "Very handsome client",
            "person_id": 1,
            "resource_id": 1,
            "service_id": 1,
            "state": "confirmed",
            "updated_at": "2012-09-20T15:34:16+02:00"
        }
    }
]`)

var testBookingsResourceFilter = []byte(`[
  {
    "booking": {
      "booked_from": "2015-12-14T12:00:00-06:00",
      "booked_to": "2015-12-14T13:00:00-06:00",
      "collection_id": null,
      "count": 1,
      "created_at": "2015-12-07T23:08:07-06:00",
      "custom_data": {},
      "event_id": null,
      "expires_at": null,
      "external_id": null,
      "id": 410577,
      "notes": "Finally a good booking",
      "person_id": 12389,
      "reminded_at": null,
      "reminder_at": null,
      "resource_id": 517,
      "service_id": 427,
      "state": "confirmed",
      "updated_at": "2015-12-07T23:08:07-06:00",
      "person": {
        "city": null,
        "country_code": null,
        "created_at": "2015-11-26T09:50:48-06:00",
        "custom_data": {},
        "date_of_birth": null,
        "email": null,
        "external_id": null,
        "id": 12389,
        "name": "Test User",
        "national_id_no": null,
        "notes": null,
        "phonenumber": null,
        "postal_code": null,
        "state": null,
        "street": null,
        "updated_at": "2015-11-26T09:50:48-06:00"
      },
      "resource": {
        "id": 517,
        "title": "jill bob"
      },
      "service": {
        "id": 427,
        "title": "Kickboxing"
      }
    }
  }
]`)

func TestBooking_list(t *testing.T) {
	_, client := mockServerClient(t)

	books, err := client.Booking(BookingParams{})
	if err != nil {
		t.Error(err)
	}

	if e := 1; len(books) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(books), e)
	}

	book := books[0]
	if e := 1; book.ID != e {
		t.Fatalf("got: %d wanted: %d", book.ID, e)
	}
	if e := 1; book.Count != e {
		t.Fatalf("got: %d wanted: %d", book.ID, e)
	}

	// client = New(ac.Name, ac.Token)
	params := BookingParams{
		ResourceID: 517,
	}

	books, err = client.Booking(params)
	if err != nil {
		t.Fatal(err)
	}
	if e := 1; len(books) != e {
		t.Fatalf("got: %d wanted: %d", len(books), e)
	}

	book = books[0]
	if e := 517; book.ResourceID != e {
		t.Errorf("got: %d wanted: %d", book.ResourceID, e)
	}
	if e := 12389; book.PersonID != e {
		t.Errorf("got: %d wanted: %d", book.PersonID, e)
	}
}

var testBookingSuccess = []byte(`{"booking":{"booked_from":"2015-11-10T08:00:00-06:00","booked_to":"2015-11-10T09:00:00-06:00","collection_id":null,"count":1,"created_at":"2015-11-07T09:09:32-06:00","custom_data":{},"event_id":null,"expires_at":null,"external_id":null,"id":410372,"notes":"Very handsome client","person_id":null,"resource_id":484,"service_id":394,"state":"confirmed","updated_at":"2015-11-07T09:09:32-06:00","resource":{"id":484,"title":"Calendar"},"service":{"id":394,"title":"Cross Fit Type"}}}`)

var fakeBookingCapacityFailure bool

func TestBooking_create(t *testing.T) {
	_, client := mockServerClient(t)

	start, _ := time.Parse(time.RFC3339, "2015-11-10T12:00:00-06:00")
	stop, _ := time.Parse(time.RFC3339, "2015-11-10T13:00:00-06:00")

	book, err := client.MakeBooking(Booking{
		// TBD "external_id": null,
		Notes: "Very handsome client",
		// PersonID:   1,
		ResourceID: 484,
		ServiceID:  394,
		Count:      1,
		CustomData: map[string]interface{}{
			"poop":   "shoot",
			"number": 5,
			"slice":  []string{"a", "b", "c"},
		},
		BookedFrom: &start,
		BookedTo:   &stop,
		State:      "confirmed",
	})
	if err != nil {
		t.Fatal(err)
	}

	if e := 410372; book.ID != e {
		t.Errorf("got: %d wanted: %d", book.ID, e)
	}
	return

	fakeBookingCapacityFailure = true
	defer func() { fakeBookingCapacityFailure = false }()

	book, err = client.MakeBooking(Booking{
		// TBD "external_id": null,
		Notes: "This will go over capacity",
		// PersonID:   1,
		ResourceID: 0,
		ServiceID:  400,
		Count:      100,
		BookedFrom: &start,
		BookedTo:   &stop,
		State:      "confirmed",
	})
	if err != ErrBookingCapacityLimit {
		t.Errorf("got: %s wanted: %s", err, ErrBookingCapacityLimit)
	}

	if book.ID != 0 {
		t.Fatal("expected nil Booking")
	}
}

var bookingDeleteResponse = []byte(`{"booking":{"booked_from":"2015-11-10T08:00:00-06:00","booked_to":"2015-11-10T09:00:00-06:00","collection_id":null,"count":1,"created_at":"2015-11-07T08:59:57-06:00","custom_data":{},"event_id":null,"expires_at":null,"external_id":null,"id":410369,"notes":"Very handsome client","person_id":null,"resource_id":484,"service_id":394,"state":"deleted","updated_at":"2015-11-08T10:14:54-06:00","resource":{"id":484,"title":"Calendar"},"service":{"id":394,"title":"Cross Fit Type"}}}`)

func TestBooking_delete(t *testing.T) {
	// TODO: repeated delete error
	// {"state":["cannot transition via \"remove\""]}
	_, client := mockServerClient(t)
	// client := realClient
	id := 410369
	book, err := client.BookingDelete(id)
	if err != nil {
		t.Fatal(err)
	}

	if e := 484; book.ResourceID != e {
		t.Errorf("got: %d wanted: %d", book.ResourceID, e)
	}

	if e := id; book.ID != e {
		t.Errorf("got: %d wanted: %d", book.ID, e)
	}
}

var bookingCancelResponse = []byte(`{"booking":{"custom_data":{"number":"5","poop":"shoot","slice":"[\"a\", \"b\", \"c\"]"},"count":1,"id":410369,"notes":"Very handsome client","resource_id":484,"service_id":394,"state":"cancelled"}}`)

func TestBooking_cancel(t *testing.T) {
	_, client := mockServerClient(t)
	id := 410369
	book, err := client.BookingCancel(id)
	if err != nil {
		t.Fatal(err)
	}

	if e := 484; book.ResourceID != e {
		t.Errorf("got: %d wanted: %d", book.ResourceID, e)
	}

	if e := id; book.ID != e {
		t.Errorf("got: %d wanted: %d", book.ID, e)
	}

	if e := "cancelled"; book.State != e {
		t.Errorf("got: %s wanted: %s", book.State, e)
	}

}
