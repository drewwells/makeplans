package makeplans

import (
	"testing"
	"time"
)

var testSlotsAll = []byte(`[
  {
    "slot": {
      "timestamp": "2015-12-14T12:00:00-06:00",
      "timestamp_end": "2015-12-14T13:00:00-06:00",
      "formatted_timestamp": "Monday, December 14, 2015, 12:00 PM",
      "formatted_timestamp_end": "Monday, December 14, 2015, 1:00 PM",
      "free": 1,
      "available_resources": [
        501
      ],
      "maximum_capacity": 1
    }
  },
  {
    "slot": {
      "timestamp": "2015-12-14T13:00:00-06:00",
      "timestamp_end": "2015-12-14T14:00:00-06:00",
      "formatted_timestamp": "Monday, December 14, 2015, 1:00 PM",
      "formatted_timestamp_end": "Monday, December 14, 2015, 2:00 PM",
      "free": 2,
      "available_resources": [
        501,
        517
      ],
      "maximum_capacity": 2
    }
  },
  {
    "slot": {
      "timestamp": "2015-12-15T12:00:00-06:00",
      "timestamp_end": "2015-12-15T13:00:00-06:00",
      "formatted_timestamp": "Tuesday, December 15, 2015, 12:00 PM",
      "formatted_timestamp_end": "Tuesday, December 15, 2015, 1:00 PM",
      "free": 1,
      "available_resources": [
        501
      ],
      "maximum_capacity": 1
    }
  },
  {
    "slot": {
      "timestamp": "2015-12-15T13:00:00-06:00",
      "timestamp_end": "2015-12-15T14:00:00-06:00",
      "formatted_timestamp": "Tuesday, December 15, 2015, 1:00 PM",
      "formatted_timestamp_end": "Tuesday, December 15, 2015, 2:00 PM",
      "free": 1,
      "available_resources": [
        501
      ],
      "maximum_capacity": 1
    }
  }
]`)

var testSlots501 = []byte(`[
  {
    "slot": {
      "timestamp": "2015-12-14T13:00:00-06:00",
      "timestamp_end": "2015-12-14T14:00:00-06:00",
      "formatted_timestamp": "Monday, December 14, 2015, 1:00 PM",
      "formatted_timestamp_end": "Monday, December 14, 2015, 2:00 PM",
      "free": 1,
      "available_resources": [
        517
      ],
      "maximum_capacity": 1
    }
  }
]`)

func TestSlot_list(t *testing.T) {
	_, client := mockServerClient(t)
	layout := "2006-01-02"

	prev, err := time.Parse(layout, "2015-12-14")
	if err != nil {
		t.Fatal(err)
	}
	to, _ := time.Parse(layout, "2015-12-15")
	_ = prev
	params := SlotParams{
		From:              prev,
		To:                to,
		SelectedResources: []int{},
	}
	slots, err := client.ServiceSlot(427, params)
	if err != nil {
		t.Error(err)
	}

	if e := 4; len(slots) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(slots), e)
	}
}

func TestSlot_filtered(t *testing.T) {
	_, client := mockServerClient(t)
	// client := New(ac.Name, ac.Token)
	layout := "2006-01-02"

	prev, err := time.Parse(layout, "2015-12-14")
	if err != nil {
		t.Fatal(err)
	}
	to, _ := time.Parse(layout, "2015-12-15")
	_ = prev
	params := SlotParams{
		From:              prev,
		To:                to,
		SelectedResources: []int{501, 517},
	}
	slots, err := client.ServiceSlot(427, params)
	if err != nil {
		t.Error(err)
	}

	if e := 4; len(slots) != e {
		t.Fatalf("wrong number of slots returned got: %d wanted: %d",
			len(slots), e)
	}

	if e := 1; slots[0].Free != e {
		t.Fatalf("got: %d wanted: %d", slots[0].Free, e)
	}

	if e := "Monday, December 14, 2015, 12:00 PM"; e !=
		slots[0].FormattedTimestamp {
		t.Fatalf("got: %s wanted: %s", slots[0].FormattedTimestamp, e)
	}

	params = SlotParams{
		From:              prev,
		To:                to,
		SelectedResources: []int{517},
	}

	slots, err = client.ServiceSlot(427, params)
	if err != nil {
		t.Fatal(err)
	}

	if e := 1; len(slots) != e {
		t.Errorf("got: %d wanted: %d", len(slots), e)
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
