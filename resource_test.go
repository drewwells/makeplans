package makeplans

import (
	"testing"
	"time"
)

var testResources = []byte(`[{"resource":{"capacity":4,"created_at":"2015-08-27T00:53:53-05:00","custom_data":{},"id":484,"opening_hours_fri":["08:00","20:00"],"opening_hours_mon":["08:00","20:00"],"opening_hours_sat":null,"opening_hours_sun":null,"opening_hours_thu":["08:00","20:00"],"opening_hours_tue":["08:00","20:00"],"opening_hours_wed":["08:00","20:00"],"title":"Calendar","updated_at":"2015-09-11T22:21:08-05:00","open_0":"08:00","close_0":"20:00","open_1":"08:00","close_1":"20:00","open_2":"08:00","close_2":"20:00","open_3":"08:00","close_3":"20:00","open_4":"08:00","close_4":"20:00","open_5":null,"close_5":null,"open_6":null,"close_6":null}}]`)

var testResource = []byte(`{
  "resource": {
    "capacity": 4,
    "created_at": "2015-08-27T00:53:53-05:00",
    "custom_data": {},
    "id": 484,
    "opening_hours_fri": [
      "08:00",
      "20:00"
    ],
    "opening_hours_mon": [
      "08:00",
      "20:00"
    ],
    "opening_hours_sat": null,
    "opening_hours_sun": null,
    "opening_hours_thu": [
      "08:00",
      "20:00"
    ],
    "opening_hours_tue": [
      "08:00",
      "20:00"
    ],
    "opening_hours_wed": [
      "08:00",
      "20:00"
    ],
    "title": "Calendar",
    "updated_at": "2015-09-11T22:21:08-05:00",
    "services": [
      {
        "active": false,
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
        "updated_at": "2015-08-27T09:22:06-05:00"
      },
      {
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
        "price": "20.0",
        "same_day": false,
        "sms_confirmation": null,
        "sms_reminder": null,
        "sms_verification": null,
        "template": null,
        "title": "Cross Fit Type",
        "updated_at": "2015-09-11T22:21:08-05:00"
      },
      {
        "active": false,
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
        "updated_at": "2015-08-27T09:21:53-05:00"
      }
    ]
  }
}`)

func TestResource_list(t *testing.T) {
	_, client := mockServerClient(t)
	ress, err := client.Resources()
	if err != nil {
		t.Fatal(err)
	}

	if len(ress) == 0 {
		t.Fatal("expected more than 0")
	}

	if e := 484; ress[0].ID != e {
		t.Errorf("got: %d wanted: %d", ress[0].ID, e)
	}
	up := ress[0]
	up.OpeningHoursMon = []string{"08:00", "10:00"}
	r, err := client.ResourceUpdate(up)
	if err != nil {
		t.Fatal(err)
	}

	if len(r.OpeningHoursMon) == 0 {
		t.Fatal("expected Mon Opening Hours")
	}

	if e := "08:00"; r.OpeningHoursMon[0] != e {
		t.Errorf("got: %s wanted: %s", r.OpeningHoursMon[0], e)
	}

}

func TestResource_get(t *testing.T) {
	_, client := mockServerClient(t)
	res, err := client.Resource(484)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID != 484 {
		t.Errorf("failed to retrieve ID: % #v\n", res)
	}

	res, err = client.Resource(100)
	if err == nil {
		t.Fatal("expected err")
	}

	if err != ErrNotFound {
		t.Errorf("invalid error returned: %s", err)
	}
}

var resourceOpeningResponse = []byte(`{"resource":{"capacity":1,"created_at":"2015-10-08T10:45:25-05:00","custom_data":{},"id":501,"opening_hours_fri":["12:00","14:00"],"opening_hours_mon":["12:00","14:00"],"opening_hours_sat":["12:00","14:00"],"opening_hours_sun":["12:00","14:00"],"opening_hours_thu":["12:00","14:00"],"opening_hours_tue":["12:00","14:00"],"opening_hours_wed":["12:00","14:00"],"title":"joe bob","updated_at":"2015-11-06T07:42:02-06:00","open_0":"12:00","close_0":"14:00","open_1":"12:00","close_1":"14:00","open_2":"12:00","close_2":"14:00","open_3":"12:00","close_3":"14:00","open_4":"12:00","close_4":"14:00","open_5":"12:00","close_5":"14:00","open_6":"12:00","close_6":"14:00","services":[{"active":true,"booking_capacity":10,"booking_type_id":1,"created_at":"2015-08-23T19:00:07-05:00","custom_data":{},"day_booking_specify_time":null,"description":"This is a cross fit type","has_day_booking":false,"id":394,"interval":60,"mail_confirmation":null,"mail_verification":null,"max_slots":75,"price":"20.0","same_day":false,"sms_confirmation":null,"sms_reminder":null,"sms_verification":null,"template":null,"title":"Cross Fit Type","updated_at":"2015-11-07T08:59:29-06:00"}]}}`)

func TestResource_opening(t *testing.T) {
	_, client := mockServerClient(t)
	layout := "2006-01-02"
	to, err := time.Parse(layout, "2015-08-01")
	if err != nil {
		t.Fatal(err)
	}
	e := 501
	rsc, err := client.ResourceOpening(e, to, to)
	if err != nil {
		t.Fatal(err)
	}

	if rsc.ID != e {
		t.Errorf("got: %d wanted: %d", rsc.ID, e)

	}

	if e := 1; len(rsc.Services) != e {
		t.Errorf("got: %d wanted: %d", len(rsc.Services), e)
	}
}
