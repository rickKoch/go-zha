package webapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiEndpoint = "https://slack.com/api/%s"
)

// Client is the client struct
type Client struct {
	token string
}

// NewClient returns new Client
func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}

// Get creates get request to the slack web api
// ex. https://api.slack.com/methods/conversations.history
func (c *Client) Get(method string, queryParams *url.Values, unmarshaledResponse interface{}) error {
	endpoint := c.endpointGenerator(method, queryParams)
	fmt.Println(endpoint.String())

	resp, err := http.Get(endpoint.String())
	if err != nil {
		fmt.Println(err, "error")
		switch e := err.(type) {
		case *url.Error:
			return e
		default:
			return fmt.Errorf("error on HTTP GET request. %#v", e)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewResponseError("response status error. status: %d", resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &unmarshaledResponse); err != nil {
		return err
	}

	return nil
}

// RtmStart begins a Real Time Messaging API session and
// reserves your application a specific URL with which to connect via websocket.
func (c *Client) RtmStart() (*RtmStart, error) {
	rtmStart := &RtmStart{}
	if err := c.Get("rtm.start", nil, &rtmStart); err != nil {
		return nil, err
	}

	return rtmStart, nil
}

// Post creates post request to the slack api
func (c *Client) Post(method string, body url.Values, response interface{}) error {
	endpoint := c.endpointGenerator(method, nil)

	resp, err := http.PostForm(endpoint.String(), body)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			return e
		default:
			panic(fmt.Sprintf("error on http post request %#v", e))
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return NewResponseError(fmt.Sprintf("response status error. status %d.", resp.StatusCode), resp)
	}

	readResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(readResp, response); err != nil {
		return err
	}

	return nil
}

// PostMessage creates post request to the chat.postMessage slack method
func (c *Client) PostMessage(message *PostMessage) (*APIResponse, error) {
	response := &APIResponse{}
	err := c.Post("chat.postMessage", message.ToURLValues(), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) endpointGenerator(method string, params *url.Values) *url.URL {
	if params == nil {
		params = &url.Values{}
	}
	params.Add("token", c.token)

	url, err := url.Parse(fmt.Sprintf(apiEndpoint, method))
	if err != nil {
		panic(err.Error())
	}
	url.RawQuery = params.Encode()

	return url
}
