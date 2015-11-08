package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Booking struct {
	BookedFrom    *time.Time             `json:"booked_from,omitempty"`
	BookedTo      *time.Time             `json:"booked_to,omitempty"`
	CustomData    map[string]interface{} `json:"custom_data,omitempty"`
	Count         int                    `json:"count,omitempty"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
	ExternalID    string                 `json:"external_id,omitempty"`
	ID            int                    `json:"id,omitempty"`
	Notes         string                 `json:"notes,omitempty"`
	PersonID      int                    `json:"person_id,omitempty"`
	ResourceID    int                    `json:"resource_id,omitempty"`
	ServiceID     int                    `json:"service_id,omitempty"`
	PublicBooking bool                   `json:"public_booking,omitempty"`
	State         string                 `json:"state,omitempty"`
	CreatedAt     *time.Time             `json:"created_at,omitempty"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty"`
}

type wrapBooking struct {
	Booking Booking `json:"booking"`
}

var (
	// BookingURL defines the entrypoint for all bookings
	BookingURL = "/bookings/"
	// ErrBookingCapacityLimit is thrown when bookings have reached the
	// available capacity of the connected resource.
	ErrBookingCapacityLimit = errors.New("error resource_id: Not available for booking at this timeerror count: More than maximum count per booking")
)

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

func (c *Client) MakeBooking(b Booking) (ret Booking, err error) {
	// Verify time slot is available
	// b.PublicBooking = true
	bs, err := json.Marshal(wrapBooking{Booking: b})
	if err != nil {
		return
	}

	bs, err = c.Do("POST", BookingURL, bytes.NewBuffer(bs))
	if err != nil {
		return
	}
	var wrap wrapBooking
	err = json.Unmarshal(bs, &wrap)
	ret = wrap.Booking
	return
}

func (c *Client) BookingDelete(id int) (ret Booking, err error) {
	sid := strconv.Itoa(id)
	bs, err := c.Do("DELETE", BookingURL+sid, nil)
	if err != nil {
		return
	}

	wrap := wrapBooking{}
	err = json.Unmarshal(bs, &wrap)
	ret = wrap.Booking
	return
}
