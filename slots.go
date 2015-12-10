package makeplans

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Slot describes an available time for a booking. It combines a trainer
// with the service they provide so a booking can be made.
//
type Slot struct {
	Timestamp             *time.Time `json:"timestamp,omitempty"`
	TimestampEnd          *time.Time `json:"timestamp_end,omitempty"`
	FormattedTimestamp    string     `json:"formatted_timestamp,omitempty"`
	FormattedTimestampEnd string     `json:"formatted_timestamp_end,omitempty"`
	Free                  int        `json:"free,omitempty"`
	OpenResources         []int      `json:"open_resources,omitempty"`
	AvailableResources    []int      `json:"available_resources,omitempty"`
}

type slotWrap struct {
	Slot Slot `json:"slot"`
}

// SlotParams are used for finding available slots for service slots and
// next available data
type SlotParams struct {
	From              time.Time
	To                time.Time
	SelectedResources []int
	OnlyFree          bool
}

var SlotURL = "/services/%d/slots" // service_id

// ServiceSlot shows all available slots for a service
func (c *Client) ServiceSlot(serviceID int, params SlotParams) ([]Slot, error) {
	path := fmt.Sprintf(SlotURL, serviceID)
	v := url.Values{}
	layout := "2006-01-02"
	from := params.From
	to := params.To
	if !from.IsZero() {
		v.Set("from", from.Format(layout))
	}
	if !to.IsZero() {
		v.Set("to", to.Format(layout))
	}
	if len(params.SelectedResources) > 0 {
		var s []string
		for _, res := range params.SelectedResources {
			s = append(s, strconv.Itoa(res))
		}
		v.Add("selected_resources", strings.Join(s, ","))
	}

	bs, err := c.Do("GET", path+"?"+v.Encode(), nil)
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
		slots[i].Timestamp = &t
	}

	return slots, err
}
