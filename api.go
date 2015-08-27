package makeplans

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var DefaultURL = "https://%s.test.makeplans.net/api/v1"

type Client struct {
	URL         string
	AccountName string
	Token       string
}

func New(account string, token string) *Client {
	return &Client{
		URL:         DefaultURL,
		Token:       token,
		AccountName: account,
	}
}

var tokenURL func(string, string) string

func init() {
	tokenURL = func(urlTmpl string, accountName string) string {
		return fmt.Sprintf(urlTmpl, accountName)
	}
}

func (c *Client) do(method string, path string, body io.Reader) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpCli := &http.Client{Transport: tr}

	req, err := http.NewRequest(method,
		tokenURL(c.URL, c.AccountName)+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "https://github.com/drewwells/makeplans")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Token, "")

	return httpCli.Do(req)
}

func (c *Client) Do(method string, path string, body io.Reader) ([]byte, error) {
	r, err := c.do(method, path, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bs, parseError(bs)
}

var ServiceURL = "/services"

type Service struct {
	Active                bool        `json:"active"`
	BookingCapacity       int         `json:"booking_capacity"`
	BookingTypeID         int         `json:"booking_type_id"`
	CustomData            interface{} `json:"custom_data"`
	DayBookingSpecifyTime interface{} `json:"day_booking_specify_time"`
	Description           string      `json:"description"`
	HasDayBooking         bool        `json:"has_day_booking"`
	ID                    int         `json:"id"`
	Interval              int         `json:"interval"`
	MailConfirmation      interface{} `json:"mail_confirmation"`
	MailVerification      interface{} `json:"mail_verification"`
	MaxSlots              int         `json:"max_slots"`
	Price                 string      `json:"price"`
	SameDay               bool        `json:"same_day"`
	SmsConfirmation       interface{} `json:"sms_confirmation"`
	SmsReminder           interface{} `json:"sms_reminder"`
	SmsVerification       interface{} `json:"sms_verification"`
	Template              interface{} `json:"template"`
	Title                 string      `json:"title"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

type serviceWrap struct {
	Service Service `json:"service"`
}

func (c *Client) Services() ([]Service, error) {
	bs, err := c.Do("GET", ServiceURL, nil)
	if err != nil {
		return nil, err
	}
	// unwrap data structure provided
	wr := []serviceWrap{}
	err = json.Unmarshal(bs, &wr)

	// Assign to a proper struct
	svcs := make([]Service, len(wr))
	for i, w := range wr {
		svcs[i] = w.Service
	}
	return svcs, err
}

func (c *Client) ServiceSave(svc Service) error {
	id := strconv.Itoa(svc.ID)
	payload, err := json.Marshal(serviceWrap{Service: svc})
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(payload)
	_, err = c.Do("PUT", ServiceURL+"/"+id, buf)
	return err
}

type E struct {
	Error struct {
		Description string
	}
}

func parseError(bs []byte) error {
	e := E{}
	err := json.Unmarshal(bs, &e)
	if err != nil {
		return nil
	}
	if len(e.Error.Description) > 0 {
		return errors.New(e.Error.Description)
	}
	return nil
}

func (c *Client) ServiceCreate(svc Service) error {
	return nil
}

func (c *Client) ServiceDelete(id int) error {
	return nil
}

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
