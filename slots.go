package makeplans

import (
	"encoding/json"
	"fmt"
	"time"
)

type Slot struct {
	Timestamp             time.Time `json:"timestamp"`
	TimestampEnd          time.Time `json:"timestamp_end"`
	FormattedTimestamp    string    `json:"formatted_timestamp"`
	FormattedTimestampEnd string    `json:"formatted_timestamp_end"`
	Free                  int       `json:"free"`
	OpenResources         []int     `json:"open_resources"`
	AvailableResources    []int     `json:"available_resources"`
}

type slotWrap struct {
	Slot Slot `json:"slot"`
}

var SlotURL = "/services/%s/slots" // service_id

// Slots shows all available slots for a service
func (c *Client) Slots(serviceID string) ([]Slot, error) {
	path := fmt.Sprintf(SlotURL, serviceID)
	bs, err := c.Do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	// unwrap data structure provided
	wr := []slotWrap{}
	err = json.Unmarshal(bs, &wr)

	// Assign to a proper struct
	slots := make([]Slot, len(wr))
	for i, w := range wr {
		slots[i] = w.Slot
	}
	return slots, err
}

var SlotNextDateURL = "/services/%s/next_available_date" // service_id

// SlotNext is the next available slot time for a specified service
// This doesn't appear to work properly, only one service is ever returned
func (c *Client) SlotNextDate(serviceID string) ([]Slot, error) {
	path := fmt.Sprintf(SlotNextDateURL, serviceID)
	bs, err := c.Do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	wrap := []struct {
		AvailableDate string `json:"available_date"`
	}{}
	// unwrap data structure provided
	err = json.Unmarshal(bs, &wrap)

	layout := "2006-01-02"
	slots := make([]Slot, len(wrap))
	for i, w := range wrap {
		// json can't unmarshal to ISO8601 shortform, so do it manually
		t, _ := time.Parse(layout, w.AvailableDate)
		t = t.AddDate(0, 1, 0)
		slots[i].Timestamp = t
	}

	return slots, err
}
