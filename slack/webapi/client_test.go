// Docs
// https://api.slack.com/web
package webapi

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

type TestResponse struct {
	APIResponse
	Test string
}

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder, _ := httpmock.NewJsonResponder(
		200,
		&TestResponse{
			APIResponse: APIResponse{OK: true},
			Test:        "test",
		},
	)

	httpmock.RegisterResponder(
		"GET",
		"https://slack.com/api/this.is.test",
		responder,
	)

	client := NewClient("testing")
	response := &TestResponse{}
	err := client.Get("this.is.test", nil, response)

	if err != nil {
		t.Errorf("something went wrong %#v", err)
	}

	if response.OK != true {
		t.Errorf("we need positive OK status %#v,", response)
	}

	if response.Test != "test" {
		t.Errorf("we should recieve test value %#v", response)
	}
}

func TestGetStatusError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	statusCode := 404
	responder := httpmock.NewStringResponder(statusCode, "testing")

	httpmock.RegisterResponder(
		"GET",
		"https://slack.com/api/test",
		responder,
	)

	client := NewClient("123123")
	response := &TestResponse{}

	err := client.Get("test", nil, response)
	switch e := err.(type) {
	case nil:
		t.Errorf("error should return on %d status", statusCode)
	case *ResponseError:
		// OK
		if e.Response.StatusCode != statusCode {
			t.Errorf("error instance includes wrond status code %d. expected %d", e.Response.StatusCode, statusCode)
		}
	default:
		t.Errorf("%#v unhendeled error", err)
	}
}

func TestGetJSONError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://slack.com/api/test",
		httpmock.NewStringResponder(200, "invalid json"),
	)

	client := NewClient("123123")
	response := &TestResponse{}

	err := client.Get("test", nil, response)
	if err == nil {
		t.Error("there should be error")
	}
}

func TestRtmStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	testURL := "https://localhost/test"

	responder, _ := httpmock.NewJsonResponder(200,
		&RtmStart{
			APIResponse: APIResponse{OK: true},
			URL:         testURL,
			Self:        nil,
		},
	)

	httpmock.RegisterResponder(
		"GET",
		"https://slack.com/api/rtm.start",
		responder,
	)

	client := NewClient("12312321")

	rtmStart, err := client.RtmStart()

	if err != nil {
		t.Errorf("something went wrong %#v", err)
	}

	if rtmStart.URL != testURL {
		t.Errorf("URL is not returned properly %#v", rtmStart)
	}
}
