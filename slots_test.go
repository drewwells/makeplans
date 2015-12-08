package makeplans

import (
	"testing"
	"time"
)

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
	layout := "2006-01-02"

	prev, err := time.Parse(layout, "2015-11-07")
	if err != nil {
		t.Fatal(err)
	}

	slots, err := client.ServiceSlot("running", prev, prev)
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
