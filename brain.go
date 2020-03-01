package zha

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Brain is a event processor
type Brain struct {
	input          chan event
	logger         *zap.Logger
	handlerTimeout time.Duration
	handlers       map[reflect.Type][]eventHandler
}

type event struct {
	Data      interface{}
	callbacks []func(event)
}

type eventHandler func(context.Context, reflect.Value) error

// NewBrain returns new Brain
func NewBrain(logger *zap.Logger, timeout time.Duration) *Brain {
	return &Brain{
		logger:         logger.Named("brain"),
		input:          make(chan event, 10),
		handlers:       make(map[reflect.Type][]eventHandler),
		handlerTimeout: timeout,
	}
}

// RegisterHandler is a register of a handler functions
func (b *Brain) RegisterHandler(fun interface{}) {
	logErr := func(err error, fields ...zapcore.Field) {
		b.logger.Error("Failed to register a handler: "+err.Error(), fields...)
	}

	handler := reflect.ValueOf(fun)
	handlerType := handler.Type()

	if handler.Kind() != reflect.Func {
		logErr(errors.New("event handler is not a function"))
		return
	}

	evtType, withContext, err := b.checkHandlerParams(handlerType)
	if err != nil {
		logErr(err)
		return
	}

	returnsErr, err := b.checkHandlerReturnValues(handlerType)

	b.logger.Debug(
		"Registering new event handler",
		zap.String("event_type", evtType.Name()),
	)

	handlerFun := b.newHandlerFunc(handler, withContext, returnsErr)
	b.handlers[evtType] = append(b.handlers[evtType], handlerFun)
}

func (b *Brain) checkHandlerReturnValues(handlerFun reflect.Type) (returnsErr bool, err error) {
	switch handlerFun.NumOut() {
	case 0:
		return false, nil
	case 1:
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if !handlerFun.Out(0).Implements(errorInterface) {
			err = errors.New("event handler function return value must be of type error")
			return
		}

		return true, nil
	default:
		return false, errors.Errorf("event handler function has more than one return value")
	}
}

func (b *Brain) checkHandlerParams(handlerFun reflect.Type) (evtType reflect.Type, withContext bool, err error) {
	numParams := handlerFun.NumIn()
	if numParams == 0 || numParams > 2 {
		err = errors.New("event handler needs one or two arguments")
		return
	}

	evtType = handlerFun.In(numParams - 1)
	withContext = numParams == 2

	if evtType.Kind() != reflect.Struct {
		err = errors.New("last event handler argument should be of type struct")
		return
	}

	if withContext {
		contextInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
		if !handlerFun.In(0).Implements(contextInterface) {
			err = errors.New("event handler argument 1 is not of type Context")
			return
		}
	}

	return evtType, withContext, nil
}

func (b *Brain) newHandlerFunc(handler reflect.Value, withContext, returnsErr bool) eventHandler {
	return func(ctx context.Context, evt reflect.Value) (handlerErr error) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				handlerErr = errors.Errorf("handler panic: %#v", err)
			}
		}()

		var args []reflect.Value
		if withContext {
			args = []reflect.Value{
				reflect.ValueOf(ctx),
				evt,
			}
		} else {
			args = []reflect.Value{evt}
		}

		results := handler.Call(args)
		if returnsErr && !results[0].IsNil() {
			return results[0].Interface().(error)
		}

		return nil
	}
}

// Process is processing the incoming events and process the outgoing event
func (b *Brain) Process(ctx context.Context) {
	for {
		select {
		case evt := <-b.input:
			b.handle(ctx, evt)
		case <-ctx.Done():
			b.handle(ctx, event{Data: ShutdownEvent{}})
			return
		}
	}
}

func (b *Brain) handle(ctx context.Context, evt event) {
	event := reflect.ValueOf(evt.Data)
	typ := event.Type()

	b.logger.Debug(
		"Handling new event",
		zap.String("event_type", typ.Name()),
		zap.Int("handlers", len(b.handlers[typ])),
	)

	for _, handler := range b.handlers[typ] {
		err := b.executeHandler(ctx, handler, event)
		if err != nil {
			b.logger.Error("Event handler failed", zap.Error(err))
		}
	}

	for _, callback := range evt.callbacks {
		callback(evt)
	}
}

func (b *Brain) executeHandler(ctx context.Context, handler eventHandler, event reflect.Value) error {
	if b.handlerTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, b.handlerTimeout)
		defer cancel()
	}

	done := make(chan error)

	go func() {
		done <- handler(ctx, event)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Emit emits the new events
func (b *Brain) Emit(eventData interface{}, callbacks ...func(event)) {
	go func() {
		b.input <- event{Data: eventData, callbacks: callbacks}
	}()
}

// BotInput interface
type BotInput interface {
	GetSenderID() string
	GetMessage() string
	GetSentAt() time.Time
	GetRoomID() string
}
