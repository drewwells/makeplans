package makeplans

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

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

// ServiceCreate creates a new service. Not all fields are required, but
// errors are thrown if required fields are missing
func (c *Client) ServiceCreate(new Service) (svc Service, err error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(serviceWrap{new})
	if err != nil {
		return
	}
	resp, err := c.Do("POST", ServiceURL, &buf)
	if err != nil {
		return
	}
	wrap := serviceWrap{}
	err = json.Unmarshal(resp, &wrap)
	if err != nil {
		// It is likely a field error has been returned,
		// attempt to unmarshal it
		var fe FieldError
		err = json.Unmarshal(resp, &fe)
		if err != nil {
			// We tried, just give up
			err = errors.New(string(resp))
			return
		}
		err = fe
		return
	}
	svc = wrap.Service
	return
}

// ServiceDelete sets the service as inactive and it no longer appears
// in lists.
func (c *Client) ServiceDelete(id int) (svc Service, err error) {
	sid := strconv.Itoa(id)
	resp, err := c.Do("DELETE", ServiceURL+"/"+sid, nil)
	if err != nil {
		return
	}
	wrap := serviceWrap{}
	err = json.Unmarshal(resp, &wrap)
	if err != nil {
		err = errors.New(string(resp))
		return
	}
	svc = wrap.Service
	return
}
