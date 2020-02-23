package webapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// AttachmentField struct
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short string `json:"short,omitempty"`
}

// MessageAttachment struct
type MessageAttachment struct {
	Fallback   string `json:"fallback"`
	Color      string `json:"color,omitempty"`
	Pretext    string `json:"pretext,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`
	Fields     []AttachmentField
	ImageURL   string `json:"image_url,omitempty"`
	ThumbURL   string `json:"thumb_url,omitempty"`
}

// PostMessage struct
type PostMessage struct {
	Channel     string
	Text        string
	Parse       string
	LinkNames   int
	Attachments []*MessageAttachment
	UnfurlLinks bool
	UnfurlMedia bool
	UserName    string
	AsUser      bool
	IconURL     string
	IconEmoji   string
}

// ToURLValues method
func (message *PostMessage) ToURLValues() url.Values {
	values := url.Values{}
	values.Add("channel", message.Channel)
	values.Add("text", message.Text)
	values.Add("parse", message.Parse)
	values.Add("link_names", string(message.LinkNames))
	values.Add("unfurl_links", strconv.FormatBool(message.UnfurlLinks))
	values.Add("unfurl_media", strconv.FormatBool(message.UnfurlMedia))
	values.Add("as_user", strconv.FormatBool(message.AsUser))
	if message.UserName != "" {
		values.Add("user_name", message.UserName)
	}
	if message.IconURL != "" {
		values.Add("icon_url", message.IconURL)
	}
	if message.IconEmoji != "" {
		values.Add("icon_emoji", message.IconEmoji)
	}
	if message.Attachments != nil {
		s, _ := json.Marshal(message.Attachments)
		values.Add("attachments", string(s))
	}

	return values
}

// NewPostMessage creates new PostMessage
func NewPostMessage(channel string, text string) *PostMessage {
	return &PostMessage{
		Channel:     channel,
		Text:        text,
		LinkNames:   1,
		Parse:       "full",
		UnfurlLinks: true,
		UnfurlMedia: true,
	}
}
