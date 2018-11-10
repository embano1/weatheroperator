package openweather

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	openWeatherAPI = "https://api.openweathermap.org/data/2.5/weather"
	httpTimeout    = 3 * time.Second
)

// Client is a client to make request against the OpenWeather API
type Client struct {
	http  *http.Client
	appID string
}

// Get calls the OpenWeather API for the given city and metric ("celsius", "fahrenheit")
func (c *Client) Get(cty, metric string) (*Result, error) {
	req, err := c.newRequest(cty, metric)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get response from OpenWeather API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("expected HTTP 200, got %d: %s", resp.StatusCode, resp.Status))
	}

	var t = Result{}
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode data: %v", err)
	}

	return &t, nil
}

func (c *Client) newRequest(cty, metric string) (*http.Request, error) {
	switch strings.ToLower(metric) {
	case "celsius":
		metric = "metric"
	case "fahrenheit":
		metric = "imperial"
	default:
		return nil, errors.New("unsupported metric")
	}

	req, err := http.NewRequest("GET", openWeatherAPI, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", cty)
	q.Add("units", metric)
	q.Add("APPID", c.appID)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// NewClient creates a ready to use OpenWeather client
func NewClient(appID string, sslVerify bool) *Client {
	tr := &http.Transport{}

	if !sslVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Timeout:   httpTimeout,
		Transport: tr,
	}

	return &Client{
		http:  client,
		appID: appID,
	}
}
