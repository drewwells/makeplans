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
	Active                bool        `json:"active,omitempty"`
	BookingCapacity       int         `json:"booking_capacity,omitempty"`
	BookingTypeID         int         `json:"booking_type_id,omitempty"`
	CustomData            interface{} `json:"custom_data,omitempty"`
	DayBookingSpecifyTime interface{} `json:"day_booking_specify_time,omitempty"`
	Description           string      `json:"description,omitempty"`
	HasDayBooking         bool        `json:"has_day_booking,omitempty"`
	ID                    int         `json:"id,omitempty"`
	Interval              int         `json:"interval,omitempty"`
	MailConfirmation      interface{} `json:"mail_confirmation,omitempty"`
	MailVerification      interface{} `json:"mail_verification,omitempty"`
	MaxSlots              int         `json:"max_slots,omitempty"`
	Price                 string      `json:"price,omitempty"`
	SameDay               bool        `json:"same_day,omitempty"`
	SmsConfirmation       interface{} `json:"sms_confirmation,omitempty"`
	SmsReminder           interface{} `json:"sms_reminder,omitempty"`
	SmsVerification       interface{} `json:"sms_verification,omitempty"`
	Template              interface{} `json:"template,omitempty"`
	Title                 string      `json:"title,omitempty"`
	CreatedAt             time.Time   `json:"created_at,omitempty"`
	UpdatedAt             time.Time   `json:"updated_at,omitempty"`
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

func (c *Client) ServiceSave(svc Service) (ret Service, err error) {
	id := strconv.Itoa(svc.ID)
	payload, err := json.Marshal(serviceWrap{Service: svc})
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(payload)
	bs, err := c.Do("PUT", ServiceURL+"/"+id, buf)
	var wrap serviceWrap
	err = json.Unmarshal(bs, &wrap)
	ret = wrap.Service
	return
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
