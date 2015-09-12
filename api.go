package makeplans

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
