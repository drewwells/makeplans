package makeplans

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	return ioutil.ReadAll(r.Body)
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
	bs, err := c.Do("PUT", ServiceURL+"/"+id, buf)
	if err != nil {
		return err
	}
	if err := parseError(bs); err != nil {
		fmt.Println("ERROR", err)
		return err
	}
	return nil
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
		log.Fatal("error unmarshalling", err)
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

var SlotURL = "/%s/slots" // service_id

func (c *Client) Slots(service_id string) ([]Slot, error) {
	path := fmt.Sprintf(SlotURL, service_id)
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
