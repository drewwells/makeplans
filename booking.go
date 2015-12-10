package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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

type BookingParams struct {
	ServiceID    int
	EventID      int
	ResourceID   int
	PersonID     int
	ExternalID   string
	Start        time.Time
	End          time.Time
	Since        time.Time
	CollectionID string
}

// Booking returns just booking matching the passed id
func (c *Client) Booking(bookingID int) (Booking, error) {
	var ret Booking
	path := BookingURL + strconv.Itoa(bookingID)
	bs, err := c.Do("GET", path, nil)
	if err != nil {
		return ret, err
	}
	var wrap wrapBooking
	err = json.Unmarshal(bs, &wrap)
	return wrap.Booking, err
}

// Bookings will return all active bookings with applied filters
func (c *Client) Bookings(params BookingParams) ([]Booking, error) {
	path := BookingURL
	var qs string
	v := url.Values{}
	if params.ServiceID > 0 {
		v.Add("service_id", strconv.Itoa(params.ServiceID))
	}
	if params.ResourceID > 0 {
		v.Add("resource_id", strconv.Itoa(params.ResourceID))
	}
	if params.PersonID > 0 {
		v.Add("person_id", strconv.Itoa(params.PersonID))
	}
	if params.EventID > 0 {
		v.Add("event_id", strconv.Itoa(params.EventID))
	}
	if len(params.ExternalID) > 0 {
		v.Add("external_id", params.ExternalID)
	}
	layout := "2006-01-02"
	if !params.Start.IsZero() {
		v.Set("start", params.Start.Format(layout))
	}
	if !params.End.IsZero() {
		v.Set("end", params.End.Format(layout))
	}

	if enc := v.Encode(); len(enc) > 0 {
		qs = "?" + enc
	}
	bs, err := c.Do("GET", path+qs, nil)

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

// BookingAll will return all bookings of all states (including declined,
// cancelled, expired and deleted
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
	b.PublicBooking = true
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

func (c *Client) mutateBooking(action string, id int) (ret Booking, err error) {
	var bs []byte
	sid := strconv.Itoa(id)
	switch action {
	case "delete":
		bs, err = c.Do("DELETE", BookingURL+sid, nil)
	case "cancel":
		bs, err = c.Do("PUT", BookingURL+sid+"/cancel", nil)
	case "verify":
		bs, err = c.Do("PUT", BookingURL+sid+"/verify", nil)
	case "confirm":
		bs, err = c.Do("PUT", BookingURL+sid+"/confirm", nil)
	case "decline":
		bs, err = c.Do("PUT", BookingURL+sid+"/decline", nil)
	default:
		err = fmt.Errorf("action %s not implemented", action)
	}
	if err != nil {
		return
	}
	wrap := wrapBooking{}
	err = json.Unmarshal(bs, &wrap)
	ret = wrap.Booking
	return
}

func (c *Client) BookingUpdate(b Booking) (Booking, error) {

	bs, err := json.Marshal(wrapBooking{Booking: b})
	if err != nil {
		return Booking{}, err
	}

	bs, err = c.Do("PUT", BookingURL+"/"+strconv.Itoa(b.ID), bytes.NewBuffer(bs))
	wrap := wrapBooking{}
	err = json.Unmarshal(bs, &wrap)
	return wrap.Booking, err
}

func (c *Client) BookingDelete(id int) (Booking, error) {
	return c.mutateBooking("delete", id)
}

func (c *Client) BookingCancel(id int) (Booking, error) {
	return c.mutateBooking("cancel", id)
}

func (c *Client) BookingVerify(id int) (Booking, error) {
	return c.mutateBooking("verify", id)
}

func (c *Client) BookingConfirm(id int) (Booking, error) {
	return c.mutateBooking("confirm", id)
}

func (c *Client) BookingDecline(id int) (Booking, error) {
	return c.mutateBooking("decline", id)
}
