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
	if err != nil {
		return nil, err
	}
	return bs, parseError(bs)
}

type FieldError map[string][]string

func (fe FieldError) Error() string {
	var msg string
	for field, errors := range fe {
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
		case ErrBookingCapacityLimit.Error():
			return ErrBookingCapacityLimit
		case ErrNotFound.Error():
			return ErrNotFound
		default:
			fmt.Println("default error", e.Error.Description)
			return errors.New(e.Error.Description)
		}
	}

	// Try again with FieldError
	var fe FieldError
	err = json.Unmarshal(bs, &fe)
	if err != nil {
		return nil
	}
	if len(fe) > 0 {
		return fe
	}
	return nil
}
