package makeplans

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var DefaultURL = "http://%s.test.makeplans.net/api/v1"

// DefaultResolver replaces the API url with the account name specified.
// This can be overridden to use a different mechanism
var DefaultResolver = func(urlTmpl string, accountName string) string {
	return fmt.Sprintf(urlTmpl, accountName)
}

type Client struct {
	URL         string
	AccountName string
	Token       string
	// annoying patch for appengine
	Client   *http.Client
	Resolver func(string, string) string
}

func New(account string, token string) *Client {
	return &Client{
		URL:         DefaultURL,
		Token:       token,
		AccountName: account,
		Resolver:    DefaultResolver,
	}
}

var showResponse = false

func (c *Client) do(method string, path string, body io.Reader) (*http.Response, error) {
	var httpCli *http.Client

	if c.Client != nil {
		httpCli = c.Client
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpCli = &http.Client{Transport: tr}
	}
	u := c.Resolver(c.URL, c.AccountName) + path

	req, err := http.NewRequest(method, u, body)
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
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	// fmt.Println(string(bs))
	if err != nil {
		return nil, err
	}
	// FIXME: parseError should happen AFTER endpoints attempt to unmarshal
	return bs, parseError(bs)
}

type FieldError map[string][]string

func (fe FieldError) Error() string {
	var msg string
	virgin := true
	for field, errors := range fe {
		if !virgin {
			msg += "\n"
		}
		virgin = false
		msg += "error " + field + ": " + strings.Join(errors, ", ")
	}

	return msg
}

type E struct {
	Error struct {
		Description string
	}
}

var (
	// ErrNotFound is a generic error returned by Makeplans. It sometimes
	// indicates an ID is invalid.
	ErrNotFound      = errors.New("Not found")
	ErrEmptyResponse = errors.New("empty response")

	// ErrEmailTaken indicates the person already exists
	ErrEmailTaken = errors.New("error email: has already been taken")
)

func parseError(bs []byte) error {
	if len(bs) == 0 {
		return ErrEmptyResponse
	}
	e := E{}
	err := json.Unmarshal(bs, &e)
	if err != nil {
		return nil
	}
	if len(e.Error.Description) > 0 {
		// Produce real errors for known http errors
		desc := e.Error.Description
		switch desc {
		case ErrEmailTaken.Error():
			return ErrEmailTaken
		case ErrBookingCapacityLimit.Error():
			return ErrBookingCapacityLimit
		case ErrNotFound.Error():
			return ErrNotFound
		default:
			fmt.Println("default error", e.Error.Description)
			return errors.New(e.Error.Description)
		}
	}

	// False positive error
	// if fal := `{"resource_id":["can't be blank"],"count":["More than maximum count per booking"]}`; fal == string(bs) {
	// 	log.Println("bypassed false error")
	// 	return nil
	// }

	// Try again with FieldError
	var fe FieldError
	err = json.Unmarshal(bs, &fe)
	if err != nil {
		return nil
	}

	// FIXME: this is weird, there's a better way to report up errors
	if msgs, ok := fe["email"]; ok {
		if len(msgs) == 1 && "error email: "+msgs[0] == ErrEmailTaken.Error() {
			return ErrEmailTaken
		}
	}

	if len(fe) > 0 {
		return fe
	}
	return nil
}
