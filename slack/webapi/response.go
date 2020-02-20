package webapi

// SlackTimeStamp represents the timestamp
type SlackTimeStamp int64

// APIResponse provides common fields shared by all API response.
type APIResponse struct {
	OK bool `json:"ok"`
}

// Self property contains details on the authenticated user.
type Self struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Created        SlackTimeStamp `json:"created"`
	ManualPresence string         `json:"manual_presence"`
}

// UserProfile information
type UserProfile struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	RealName           string `json:"real_name"`
	RealNameNormalized string `json:"real_name_normalized"`
	Email              string `json:"email"`
	Skype              string `json:"skype"`
	Phone              string `json:"phone"`
	Image24            string `json:"image_24"`
	Image32            string `json:"image_32"`
	Image48            string `json:"image_48"`
	Image72            string `json:"image_72"`
	Image192           string `json:"image_192"`
	ImageOriginal      string `json:"image_original"`
	Title              string `json:"title"`
}

// User property contains a list of user object
type User struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Deleted           bool        `json:"deleted"`
	Color             string      `json:"color"`
	RealName          string      `json:"real_name"`
	TZ                string      `json:"tz,omitempty"`
	TZLabel           string      `json:"tz_label"`
	TZOffset          int         `json:"tz_offset"`
	Profile           UserProfile `json:"profile"`
	IsBot             bool        `json:"is_bot"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	Has2FA            bool        `json:"has_2fa"`
	HasFiles          bool        `json:"has_files"`
	Presence          string      `json:"presence"`
}

// Team provides information about your team.
type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// Topic provides information about your topic.
type Topic struct {
	Value   string         `json:"value"`
	Creator string         `json:"creator"`
	LastSet SlackTimeStamp `json:"last_set"`
}

// Purpose provides information about your purpose.
type Purpose struct {
	Value   string         `json:"value"`
	Creator string         `json:"creator"`
	LastSet SlackTimeStamp `json:"last_set"`
}

// Message provides information about your message.
type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

// Channel provides information about your channel.
type Channel struct {
	ID                 string         `json:"id"`
	Created            SlackTimeStamp `json:"created"`
	IsOpen             bool           `json:"is_open"`
	LastRead           string         `json:"last_read,omitempty"`
	Latest             *Message       `json:"latest,omitempty"`
	UnreadCount        int            `json:"unread_count,omitempty"`
	UnreadCountDisplay int            `json:"unread_count_display,omitempty"`
	Name               string         `json:"name"`
	Creator            string         `json:"creator"`
	IsArchived         bool           `json:"is_archived"`
	Members            []string       `json:"members"`
	NumMembers         int            `json:"num_members,omitempty"`
	Topic              Topic          `json:"topic"`
	Purpose            Purpose        `json:"purpose"`
	IsChannel          bool           `json:"is_channel"`
	IsGeneral          bool           `json:"is_general"`
	IsMember           bool           `json:"is_member"`
}

// Group provides information about your group.
type Group struct {
	ID                 string         `json:"id"`
	Created            SlackTimeStamp `json:"created"`
	IsOpen             bool           `json:"is_open"`
	LastRead           string         `json:"last_read,omitempty"`
	Latest             *Message       `json:"latest,omitempty"`
	UnreadCount        int            `json:"unread_count,omitempty"`
	UnreadCountDisplay int            `json:"unread_count_display,omitempty"`
	Name               string         `json:"name"`
	Creator            string         `json:"creator"`
	IsArchived         bool           `json:"is_archived"`
	Members            []string       `json:"members"`
	NumMembers         int            `json:"num_members,omitempty"`
	Topic              Topic          `json:"topic"`
	Purpose            Purpose        `json:"purpose"`
	IsGroup            bool           `json:"is_group"`
}

// Icons provides information about your icons.
type Icons struct {
	Image48 string `json:"image_48"`
}

// Bot provides information about your bot.
type Bot struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Icons   Icons  `json:"icons"`
}

// IM struct
type IM struct {
	ID                 string         `json:"id"`
	Created            SlackTimeStamp `json:"created"`
	IsOpen             bool           `json:"is_open"`
	LastRead           string         `json:"last_read,omitempty"`
	Latest             *Message       `json:"latest,omitempty"`
	UnreadCount        int            `json:"unread_count,omitempty"`
	UnreadCountDisplay int            `json:"unread_count_display,omitempty"`
	IsIM               bool           `json:"is_im"`
	User               string         `json:"user"`
	IsUserDeleted      bool           `json:"is_user_deleted"`
}

// RtmStart begins a Real Time Messaging API session and
// reserves your application a specific URL with which to connect via websocket.
type RtmStart struct {
	APIResponse

	// TODO consider net/url
	URL string `json:"url,omitempty"`

	Self     *Self     `json:"self,omitempty"`
	Team     *Team     `json:"team,omitempty"`
	Users    []User    `json:"users,omitempty"`
	Channels []Channel `json:"channels,omitempty"`
	Groups   []Group   `json:"groups,omitempty"`
	Bots     []Bot     `json:"bots,omitempty"`
	IMs      []IM      `json:"ims,omitempty"`
}
