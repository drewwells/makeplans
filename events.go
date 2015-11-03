package makeplans

import (
	"encoding/json"
	"time"
)

type Event struct {
	Capacity    int         `json:"capacity,omitempty"`
	CreatedAt   *time.Time  `json:"created_at,omitempty"`
	CustomData  interface{} `json:"custom_data,omitempty"`
	Description interface{} `json:"description,omitempty"`
	End         *time.Time  `json:"end,omitempty"`
	ID          int         `json:"id,omitempty"`
	ResourceID  int         `json:"resource_id,omitempty"`
	Published   bool        `json:"published,omitempty"`
	Start       *time.Time  `json:"start,omitempty"`
	ServiceID   int         `json:"service_id,omitempty"`
	Title       string      `json:"title,omitempty"`
	UpdatedAt   *time.Time  `json:"updated_at,omitempty"`
}

type wrapEvent struct {
	Event Event `json:"event"`
}

var EventsURL = "/events"

// Events list of events
func (c *Client) Events() ([]Event, error) {
	bs, err := c.Do("GET", EventsURL, nil)
	if err != nil {
		return nil, err
	}
	var wraps []wrapEvent
	err = json.Unmarshal(bs, &wraps)
	if err != nil {
		return nil, err
	}

	evts := make([]Event, len(wraps))
	for i, wrap := range wraps {
		evts[i] = wrap.Event
	}
	return evts, nil
}
