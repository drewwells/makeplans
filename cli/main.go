package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/kr/pretty"
)

type account struct {
	Name  string
	Token string
}

func main() {
	bs, err := ioutil.ReadFile("account.json")
	if err != nil {
		log.Fatal(err)
	}

	ac := account{}
	err = json.Unmarshal(bs, &ac)
	if err != nil {
		log.Fatal(err)
	}

	client := Client{
		Token:       ac.Token,
		AccountName: ac.Name,
	}

	svcs, err := client.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #v\n", pretty.Formatter(svcs))

	svc := svcs[0]
	svc.Price = "20.0"
	err = client.ServiceSave(svc)
	if err != nil {
		log.Fatal(err)
	}

	svcs, _ = client.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #v\n", pretty.Formatter(svcs[0]))
}

var BaseURL = "https://%s.test.makeplans.net/api/v1"

type Client struct {
	AccountName string
	Token       string
}

func (c *Client) do(method string, path string, body io.Reader) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpCli := &http.Client{Transport: tr}

	req, err := http.NewRequest(method,
		fmt.Sprintf(BaseURL, c.AccountName)+path, body)
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
	CreatedAt             string      `json:"created_at"`
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
	UpdatedAt             string      `json:"updated_at"`
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
	fmt.Println("preparing:", string(payload))
	buf := bytes.NewBuffer(payload)
	bs, err := c.Do("PUT", ServiceURL+"/"+id, buf)
	if err != nil {
		return err
	}
	fmt.Println("return", string(bs))
	return nil
}
