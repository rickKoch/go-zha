package main

import (
	"fmt"

	"gitlab.com/kochevRisto/go-zha/slack/webapi"
)

// TestResponse test
type TestResponse struct {
	webapi.APIResponse
	Self     *webapi.Self      `json:"self,omitempty"`
	Team     *webapi.Team      `json:"team,omitempty"`
	Channels []*webapi.Channel `json:"channels,omitempty"`
}

func main() {

	client := webapi.NewClient("xoxb-953511947447-940297123539-R4O3Cxt6W2Od0UuEQrOTXqwt")

	response, err := client.PostMessage(&webapi.PostMessage{
		Channel: "CTPF3NZ8R",
		Text:    "testing 123123",
	})

	if err != nil {
		fmt.Printf("error: %#v", err)
	}

	fmt.Println(response)
	// rtmStart, err := client.RtmStart()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(rtmStart)

	// // response := &TestResponse{}
	// // params := &url.Values{}
	// // params.Add("channel", "CTPF3NZ8R")
	// // err := client.Get("conversations.history", nil, response)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }

	// // fmt.Println(response)

	// response := &TestResponse{}
	// // params := &url.Values{}
	// // params.Add("channel", "CTPF3NZ8R")
	// err := client.Get("channels.list", nil, response)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// for _, channel := range response.Channels {
	// 	fmt.Printf("\n%#v", channel)
	// }
}
