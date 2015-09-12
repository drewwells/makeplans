package makeplans

import (
	"encoding/json"
	"time"
)

type Booking struct {
	BookedFrom time.Time `json:"booked_from"`
	BookedTo   time.Time `json:"booked_to"`
	CreatedAt  time.Time `json:"created_at"`
	CustomData struct {
	} `json:"custom_data"`
	Count      int         `json:"count"`
	ExpiresAt  interface{} `json:"expires_at"`
	ExternalID interface{} `json:"external_id"`
	ID         int         `json:"id"`
	Notes      string      `json:"notes"`
	PersonID   int         `json:"person_id"`
	ResourceID int         `json:"resource_id"`
	ServiceID  int         `json:"service_id"`
	State      string      `json:"state"`
	UpdatedAt  time.Time   `json:"updated_at"`
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
