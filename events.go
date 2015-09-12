package makeplans

import (
	"encoding/json"
	"time"
)

type Event struct {
	Capacity    int         `json:"capacity"`
	CreatedAt   time.Time   `json:"created_at"`
	CustomData  interface{} `json:"custom_data"`
	Description interface{} `json:"description"`
	End         time.Time   `json:"end"`
	ID          int         `json:"id"`
	ResourceID  int         `json:"resource_id"`
	Published   bool        `json:"published"`
	Start       time.Time   `json:"start"`
	ServiceID   int         `json:"service_id"`
	Title       string      `json:"title"`
	UpdatedAt   time.Time   `json:"updated_at"`
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
