package makeplans

import (
	"encoding/json"
	"time"
)

type Booking struct {
	BookedFrom time.Time `json:"booked_from,omitempty"`
	BookedTo   time.Time `json:"booked_to,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	CustomData struct {
	} `json:"custom_data,omitempty"`
	Count      int         `json:"count,omitempty"`
	ExpiresAt  interface{} `json:"expires_at,omitempty"`
	ExternalID interface{} `json:"external_id,omitempty"`
	ID         int         `json:"id,omitempty"`
	Notes      string      `json:"notes,omitempty"`
	PersonID   int         `json:"person_id,omitempty"`
	ResourceID int         `json:"resource_id,omitempty"`
	ServiceID  int         `json:"service_id,omitempty"`
	State      string      `json:"state,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty"`
}

type wrapBooking struct {
	Booking Booking `json:"booking"`
}

var BookingURL = "/bookings"

// Booking will return all active bookings
func (c *Client) Booking() ([]Booking, error) {
	bs, err := c.Do("GET", BookingURL, nil)
	if err != nil {
		return nil, err
	}
	var wrap []wrapBooking
	err = json.Unmarshal(bs, &wrap)
	if err != nil {
		return nil, err
	}

	books := make([]Booking, len(wrap))
	for i, book := range wrap {
		books[i] = book.Booking
	}
	return books, nil
}

var BookingAllURL = "/bookings/all"

// BookingAll will return all bookings of all states (including declined, cancelled, expired and deleted
func (c *Client) BookingAll() ([]Booking, error) {
	bs, err := c.Do("GET", BookingAllURL, nil)
	if err != nil {
		return nil, err
	}
	var wrap []wrapBooking
	err = json.Unmarshal(bs, &wrap)
	if err != nil {
		return nil, err
	}

	books := make([]Booking, len(wrap))
	for i, book := range wrap {
		books[i] = book.Booking
	}
	return books, nil
}
