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

func TestBooking_list(t *testing.T) {
	_, client := mockServerClient(t)

	books, err := client.Booking()
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

}

var testBookingSuccess = []byte(`{"booking":{"booked_from":"2015-11-10T08:00:00-06:00","booked_to":"2015-11-10T09:00:00-06:00","collection_id":null,"count":1,"created_at":"2015-11-07T09:09:32-06:00","custom_data":{},"event_id":null,"expires_at":null,"external_id":null,"id":410372,"notes":"Very handsome client","person_id":null,"resource_id":484,"service_id":394,"state":"confirmed","updated_at":"2015-11-07T09:09:32-06:00","resource":{"id":484,"title":"Calendar"},"service":{"id":394,"title":"Cross Fit Type"}}}`)

var fakeBookingCapacityFailure bool

func TestBooking_create(t *testing.T) {
	// client := New(ac.Name, ac.Token)
	_, client := mockServerClient(t)

	// "timestamp": "2015-11-10T08:00:00-06:00",
	//         "timestamp_end": "2015-11-10T09:00:00-06:00",
	//         "formatted_timestamp": "Tuesday, November 10, 2015, 8:00 AM",
	//         "formatted_timestamp_end": "Tuesday, November 10, 2015, 9:00 AM",
	//         "free": 4,
	//         "available_resources": [
	//             484
	//         ]

	start, _ := time.Parse(time.RFC3339, "2015-11-10T08:00:00-06:00")
	stop, _ := time.Parse(time.RFC3339, "2015-11-10T09:00:00-06:00")

	book, err := client.MakeBooking(Booking{
		// TBD "external_id": null,
		Notes: "Very handsome client",
		// PersonID:   1,
		ResourceID: 484,
		ServiceID:  394,
		Count:      100,
		BookedFrom: &start,
		BookedTo:   &stop,
		State:      "confirmed",
	})
	if err != nil {
		t.Fatal(err)
	}

	if book.ID == 0 {
		t.Errorf("unexpected nil book: % #v\n", book)
	}

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