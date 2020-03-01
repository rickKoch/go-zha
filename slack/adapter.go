package slack

import (
	"fmt"
	"io"
	"time"

	"gitlab.com/kochevRisto/go-zha"
	"gitlab.com/kochevRisto/go-zha/slack/retry"
	"gitlab.com/kochevRisto/go-zha/slack/rtmapi"
	"gitlab.com/kochevRisto/go-zha/slack/webapi"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

// Config is a slack config
type Config struct {
	Token  string
	Logger *zap.Logger
}

// Adapter struct
type Adapter struct {
	WebAPIClient        *webapi.Client
	RtmAPIClient        *rtmapi.Client
	tryPing             chan bool
	Events              chan rtmapi.DecodedEvent
	outgoingEventID     *rtmapi.OutgoingEventID
	OutgoingMessages    chan *rtmapi.TextMessage
	StartNewRtm         chan bool
	webSocketConnection *websocket.Conn
	Stopper             chan bool
	stopAll             chan bool
	logger              *zap.Logger
}

// NewAdapter generates new Adapter
func NewAdapter(token string, opts ...Option) zha.Option {
	return func(b *zha.Bot) error {
		conf := Config{Token: token}
		for _, opt := range opts {
			err := opt(&conf)
			if err != nil {
				return nil
			}
		}

		if conf.Logger == nil {
			conf.Logger = b.Logger
		}

		b.Adapter = NewSlackAdapter(&conf)

		return nil
	}
}

// NewSlackAdapter creates new Slack instnce
func NewSlackAdapter(config *Config) *Adapter {
	a := &Adapter{
		WebAPIClient:     webapi.NewClient(config.Token),
		RtmAPIClient:     rtmapi.NewClient(),
		tryPing:          make(chan bool),
		Events:           make(chan rtmapi.DecodedEvent, 100),
		outgoingEventID:  rtmapi.NewOutgoingEventID(),
		OutgoingMessages: make(chan *rtmapi.TextMessage, 100),
		StartNewRtm:      make(chan bool),
		Stopper:          make(chan bool),
		stopAll:          make(chan bool),
	}

	if a.logger == nil {
		a.logger = zap.NewNop()
	}

	return a
}

// Register starts slacker
func (s *Adapter) Register(b *zha.Brain) {
	go s.supervise()
	go s.sendEnqueuedMessage()
	go s.receiveEvent(b)

	s.StartNewRtm <- true
}

func (s *Adapter) supervise() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.StartNewRtm:
			s.disconnect()
			if err := s.connect(); err != nil {
				s.logger.Error("error on connect")
				s.Stopper <- true
			}
		case <-s.Stopper:
			close(s.stopAll)
			s.disconnect()
			return
		case <-ticker.C:
			s.checkConnection()
		case <-s.tryPing:
			s.checkConnection()
		}
	}
}

func (s *Adapter) connect() error {
	rtmInfo, err := s.fetchRtmInfo()
	if err != nil {
		return err
	}

	conn, err := s.connectRtm(rtmInfo)
	if err != nil {
		return err
	}

	s.webSocketConnection = conn
	return nil
}

func (s *Adapter) disconnect() {
	if s.webSocketConnection == nil {
		return
	}

	if err := s.webSocketConnection.Close(); err != nil {
		s.logger.Error(
			"error on connection close. type %T. value: %+v.",
			zap.Error(err),
			zap.Error(err),
		)
	}
}

func (s *Adapter) checkConnection() {
	s.logger.Debug("checking connection status with Ping payload.")
	ping := rtmapi.NewPing(s.outgoingEventID)
	if err := websocket.JSON.Send(s.webSocketConnection, ping); err != nil {
		s.logger.Error("failed sending Ping payload", zap.Error(err))
		s.StartNewRtm <- true
	}
}

func (s *Adapter) fetchRtmInfo() (*webapi.RtmStart, error) {
	var rtmStart *webapi.RtmStart
	err := retry.Interval(10, func() error {
		r, e := s.WebAPIClient.RtmStart()
		rtmStart = r
		return e
	}, 500*time.Microsecond)

	return rtmStart, err
}

func (s *Adapter) connectRtm(rtm *webapi.RtmStart) (*websocket.Conn, error) {
	var conn *websocket.Conn
	err := retry.Interval(10, func() error {
		c, e := s.RtmAPIClient.Connect(rtm.URL)
		conn = c
		return e
	}, 500*time.Microsecond)

	return conn, err
}

func (s *Adapter) sendEnqueuedMessage() {
	for {
		select {
		case <-s.stopAll:
			return
		case message := <-s.OutgoingMessages:
			if s.webSocketConnection == nil {
				continue
			}

			fmt.Printf("message: %#v", message)

			event := rtmapi.NewOutgoingMessage(s.outgoingEventID, message)
			if err := websocket.JSON.Send(s.webSocketConnection, event); err != nil {
				fmt.Println(err, "eveneeeet")
				s.logger.Error("failed to send event", zap.Any("event", event))
			}
		}
	}
}

func (s *Adapter) receiveEvent(b *zha.Brain) {
	for {
		select {
		case <-s.stopAll:
			s.logger.Error("stop receiving events due to stop queue")
			return
		default:
			if s.webSocketConnection == nil {
				continue
			}

			payload, err := rtmapi.ReceivePayload(s.webSocketConnection)
			if err == io.EOF {
				s.tryPing <- true
				continue
			}

			if err != nil {
				s.logger.Error("error on receiving payload", zap.Any("error", err.Error()))
				continue
			}

			event, err := s.RtmAPIClient.DecodePayload(payload)
			if err != nil {
				switch err.(type) {
				case *rtmapi.EventTypeError:
					s.logger.Warn("malformed payload was passed.", zap.Any("payload", payload))
				case *rtmapi.ReplyStatusError:
					s.logger.Error("something was wrong with previous posted message. %#v", zap.Any("error", err.Error()))
				default:
					s.logger.Error("unhandled error occured on payload decode. %#v", zap.Any("error", err.Error()))
				}
			}

			if event == nil {
				continue
			}

			s.logger.Debug("Received message", zap.Any("event", event))

			if botInput, ok := event.(zha.BotInput); ok {
				b.Emit(zha.ReciveMessageEvent{
					Text:     botInput.GetMessage(),
					ChannelD: botInput.GetRoomID(),
				})
			} else {
				// s.logger.Warn("unhandeled")
			}
		}
	}
}

// Close should shutdown the adapter
func (s *Adapter) Close() error {
	return nil
}

// Send sends message to slack
func (s *Adapter) Send(text, channelID string) error {
	s.logger.Info("Sending message to channel",
		zap.String("channel_id", channelID),
	)
	if s.webSocketConnection == nil {
		return nil
	}

	message := rtmapi.NewTextMessage(channelID, text)

	fmt.Printf("messsageee: %#v", message)

	event := rtmapi.NewOutgoingMessage(s.outgoingEventID, message)
	if err := websocket.JSON.Send(s.webSocketConnection, event); err != nil {
		s.logger.Error("failed to send event", zap.Any("error", err.Error()))
	}

	return nil
}
