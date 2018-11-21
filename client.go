package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	BASEURL   string = "https://api.bitflyer.jp"
	SIGN             = "ACCESS-SIGN"
	TIMESTAMP        = "ACCESS-TIMESTAMP"
	KEY              = "ACCESS-KEY"
)

type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

type Request struct {
	Method   string
	Endpoint string
	Header   http.Header
	Query    url.Values
	Form     url.Values
	Body     io.Reader
	FullUrl  string
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    BASEURL,
		HTTPClient: http.DefaultClient,
	}
}

func (r *Request) validate() (err error) {
	if r.Query == nil {
		r.Query = url.Values{}
	}
	if r.Form == nil {
		r.Form = url.Values{}
	}
	return nil
}

func (c *Client) parseRequest(r *Request) error {
	r.validate()
	now := time.Now().String()
	fullUrl := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)
	queryString := r.Query.Encode()
	if queryString != "" {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, queryString)
	}
	body := new(bytes.Buffer)
	if r.Body != nil {
		body.ReadFrom(r.Body)
	}
	raw := fmt.Sprintf("%s%s%s%s", now, r.Method, r.Endpoint, body.String())
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, err := mac.Write([]byte(raw))
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set(KEY, c.APIKey)
	header.Set(TIMESTAMP, now)
	header.Set(SIGN, hex.EncodeToString(mac.Sum(nil)))
	r.Header = header
	r.FullUrl = fullUrl
	return nil
}

func (c *Client) call(r *Request) ([]byte, error) {
	err := c.parseRequest(r)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.Method, r.FullUrl, r.Body)
	if err != nil {
		return []byte{}, err
	}
	req.Header = r.Header
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return byteArray, nil
}
