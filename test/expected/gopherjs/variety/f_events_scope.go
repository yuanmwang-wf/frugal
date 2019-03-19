// Autogenerated by Frugal Compiler (3.0.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package variety

import (
	"fmt"

	"github.com/Workiva/frugal/lib/gopherjs/frugal"
	"github.com/Workiva/frugal/lib/gopherjs/thrift"
)

const delimiter = "."

// This docstring gets added to the generated code because it has
// the @ sign. Prefix specifies topic prefix tokens, which can be static or
// variable.
type EventsPublisher interface {
	Open() error
	Close() error
	PublishEventCreated(ctx frugal.FContext, user string, req *Event) error
	PublishSomeInt(ctx frugal.FContext, user string, req int64) error
	PublishSomeStr(ctx frugal.FContext, user string, req string) error
	PublishSomeList(ctx frugal.FContext, user string, req []map[ID]*Event) error
}

type eventsPublisher struct {
	transport       frugal.FPublisherTransport
	protocolFactory *frugal.FProtocolFactory
	methods         map[string]*frugal.Method
}

func NewEventsPublisher(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) EventsPublisher {
	transport, protocolFactory := provider.NewPublisher()
	methods := make(map[string]*frugal.Method)
	publisher := &eventsPublisher{
		transport:       transport,
		protocolFactory: protocolFactory,
		methods:         methods,
	}
	middleware = append(middleware, provider.GetMiddleware()...)
	methods["publishEventCreated"] = frugal.NewMethod(publisher, publisher.publishEventCreated, "publishEventCreated", middleware)
	methods["publishSomeInt"] = frugal.NewMethod(publisher, publisher.publishSomeInt, "publishSomeInt", middleware)
	methods["publishSomeStr"] = frugal.NewMethod(publisher, publisher.publishSomeStr, "publishSomeStr", middleware)
	methods["publishSomeList"] = frugal.NewMethod(publisher, publisher.publishSomeList, "publishSomeList", middleware)
	return publisher
}

func (p *eventsPublisher) Open() error {
	return p.transport.Open()
}

func (p *eventsPublisher) Close() error {
	return p.transport.Close()
}

// This is a docstring.
func (p *eventsPublisher) PublishEventCreated(ctx frugal.FContext, user string, req *Event) error {
	ret := p.methods["publishEventCreated"].Invoke([]interface{}{ctx, user, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *eventsPublisher) publishEventCreated(ctx frugal.FContext, user string, req *Event) error {
	ctx.AddRequestHeader("_topic_user", user)
	op := "EventCreated"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := req.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", req), err)
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return p.transport.Publish(topic, buffer.Bytes())
}

func (p *eventsPublisher) PublishSomeInt(ctx frugal.FContext, user string, req int64) error {
	ret := p.methods["publishSomeInt"].Invoke([]interface{}{ctx, user, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *eventsPublisher) publishSomeInt(ctx frugal.FContext, user string, req int64) error {
	ctx.AddRequestHeader("_topic_user", user)
	op := "SomeInt"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteI64(int64(req)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return p.transport.Publish(topic, buffer.Bytes())
}

func (p *eventsPublisher) PublishSomeStr(ctx frugal.FContext, user string, req string) error {
	ret := p.methods["publishSomeStr"].Invoke([]interface{}{ctx, user, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *eventsPublisher) publishSomeStr(ctx frugal.FContext, user string, req string) error {
	ctx.AddRequestHeader("_topic_user", user)
	op := "SomeStr"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteString(string(req)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return p.transport.Publish(topic, buffer.Bytes())
}

func (p *eventsPublisher) PublishSomeList(ctx frugal.FContext, user string, req []map[ID]*Event) error {
	ret := p.methods["publishSomeList"].Invoke([]interface{}{ctx, user, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *eventsPublisher) publishSomeList(ctx frugal.FContext, user string, req []map[ID]*Event) error {
	ctx.AddRequestHeader("_topic_user", user)
	op := "SomeList"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteListBegin(thrift.MAP, len(req)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range req {
		if err := oprot.WriteMapBegin(thrift.I64, thrift.STRUCT, len(v)); err != nil {
			return thrift.PrependError("error writing map begin: ", err)
		}
		for k, v := range v {
			if err := oprot.WriteI64(int64(k)); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteMapEnd(); err != nil {
			return thrift.PrependError("error writing map end: ", err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return p.transport.Publish(topic, buffer.Bytes())
}

// This docstring gets added to the generated code because it has
// the @ sign. Prefix specifies topic prefix tokens, which can be static or
// variable.
type EventsSubscriber interface {
	SubscribeEventCreated(user string, handler func(frugal.FContext, *Event)) (*frugal.FSubscription, error)
	SubscribeSomeInt(user string, handler func(frugal.FContext, int64)) (*frugal.FSubscription, error)
	SubscribeSomeStr(user string, handler func(frugal.FContext, string)) (*frugal.FSubscription, error)
	SubscribeSomeList(user string, handler func(frugal.FContext, []map[ID]*Event)) (*frugal.FSubscription, error)
}

// This docstring gets added to the generated code because it has
// the @ sign. Prefix specifies topic prefix tokens, which can be static or
// variable.
type EventsErrorableSubscriber interface {
	SubscribeEventCreatedErrorable(user string, handler func(frugal.FContext, *Event) error) (*frugal.FSubscription, error)
	SubscribeSomeIntErrorable(user string, handler func(frugal.FContext, int64) error) (*frugal.FSubscription, error)
	SubscribeSomeStrErrorable(user string, handler func(frugal.FContext, string) error) (*frugal.FSubscription, error)
	SubscribeSomeListErrorable(user string, handler func(frugal.FContext, []map[ID]*Event) error) (*frugal.FSubscription, error)
}

type eventsSubscriber struct {
	provider   *frugal.FScopeProvider
	middleware []frugal.ServiceMiddleware
}

func NewEventsSubscriber(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) EventsSubscriber {
	middleware = append(middleware, provider.GetMiddleware()...)
	return &eventsSubscriber{provider: provider, middleware: middleware}
}

func NewEventsErrorableSubscriber(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) EventsErrorableSubscriber {
	middleware = append(middleware, provider.GetMiddleware()...)
	return &eventsSubscriber{provider: provider, middleware: middleware}
}

// This is a docstring.
func (l *eventsSubscriber) SubscribeEventCreated(user string, handler func(frugal.FContext, *Event)) (*frugal.FSubscription, error) {
	return l.SubscribeEventCreatedErrorable(user, func(fctx frugal.FContext, arg *Event) error {
		handler(fctx, arg)
		return nil
	})
}

// This is a docstring.
func (l *eventsSubscriber) SubscribeEventCreatedErrorable(user string, handler func(frugal.FContext, *Event) error) (*frugal.FSubscription, error) {
	op := "EventCreated"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvEventCreated(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *eventsSubscriber) recvEventCreated(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, *Event) error) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeEventCreated", l.middleware)
	return func(transport thrift.TTransport) error {
		iprot := pf.GetProtocol(transport)
		ctx, err := iprot.ReadRequestHeader()
		if err != nil {
			return err
		}

		name, _, _, err := iprot.ReadMessageBegin()
		if err != nil {
			return err
		}

		if name != op {
			iprot.Skip(thrift.STRUCT)
			iprot.ReadMessageEnd()
			return thrift.NewTApplicationException(frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, "Unknown function"+name)
		}
		req := NewEvent()
		if err := req.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", req), err)
		}
		iprot.ReadMessageEnd()

		return method.Invoke([]interface{}{ctx, req}).Error()
	}
}

func (l *eventsSubscriber) SubscribeSomeInt(user string, handler func(frugal.FContext, int64)) (*frugal.FSubscription, error) {
	return l.SubscribeSomeIntErrorable(user, func(fctx frugal.FContext, arg int64) error {
		handler(fctx, arg)
		return nil
	})
}

func (l *eventsSubscriber) SubscribeSomeIntErrorable(user string, handler func(frugal.FContext, int64) error) (*frugal.FSubscription, error) {
	op := "SomeInt"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvSomeInt(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *eventsSubscriber) recvSomeInt(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, int64) error) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeSomeInt", l.middleware)
	return func(transport thrift.TTransport) error {
		iprot := pf.GetProtocol(transport)
		ctx, err := iprot.ReadRequestHeader()
		if err != nil {
			return err
		}

		name, _, _, err := iprot.ReadMessageBegin()
		if err != nil {
			return err
		}

		if name != op {
			iprot.Skip(thrift.STRUCT)
			iprot.ReadMessageEnd()
			return thrift.NewTApplicationException(frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, "Unknown function"+name)
		}
		var req int64
		if v, err := iprot.ReadI64(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			req = v
		}
		iprot.ReadMessageEnd()

		return method.Invoke([]interface{}{ctx, req}).Error()
	}
}

func (l *eventsSubscriber) SubscribeSomeStr(user string, handler func(frugal.FContext, string)) (*frugal.FSubscription, error) {
	return l.SubscribeSomeStrErrorable(user, func(fctx frugal.FContext, arg string) error {
		handler(fctx, arg)
		return nil
	})
}

func (l *eventsSubscriber) SubscribeSomeStrErrorable(user string, handler func(frugal.FContext, string) error) (*frugal.FSubscription, error) {
	op := "SomeStr"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvSomeStr(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *eventsSubscriber) recvSomeStr(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, string) error) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeSomeStr", l.middleware)
	return func(transport thrift.TTransport) error {
		iprot := pf.GetProtocol(transport)
		ctx, err := iprot.ReadRequestHeader()
		if err != nil {
			return err
		}

		name, _, _, err := iprot.ReadMessageBegin()
		if err != nil {
			return err
		}

		if name != op {
			iprot.Skip(thrift.STRUCT)
			iprot.ReadMessageEnd()
			return thrift.NewTApplicationException(frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, "Unknown function"+name)
		}
		var req string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			req = v
		}
		iprot.ReadMessageEnd()

		return method.Invoke([]interface{}{ctx, req}).Error()
	}
}

func (l *eventsSubscriber) SubscribeSomeList(user string, handler func(frugal.FContext, []map[ID]*Event)) (*frugal.FSubscription, error) {
	return l.SubscribeSomeListErrorable(user, func(fctx frugal.FContext, arg []map[ID]*Event) error {
		handler(fctx, arg)
		return nil
	})
}

func (l *eventsSubscriber) SubscribeSomeListErrorable(user string, handler func(frugal.FContext, []map[ID]*Event) error) (*frugal.FSubscription, error) {
	op := "SomeList"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvSomeList(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *eventsSubscriber) recvSomeList(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, []map[ID]*Event) error) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeSomeList", l.middleware)
	return func(transport thrift.TTransport) error {
		iprot := pf.GetProtocol(transport)
		ctx, err := iprot.ReadRequestHeader()
		if err != nil {
			return err
		}

		name, _, _, err := iprot.ReadMessageBegin()
		if err != nil {
			return err
		}

		if name != op {
			iprot.Skip(thrift.STRUCT)
			iprot.ReadMessageEnd()
			return thrift.NewTApplicationException(frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, "Unknown function"+name)
		}
		_, size, err := iprot.ReadListBegin()
		if err != nil {
			return thrift.PrependError("error reading list begin: ", err)
		}
		req := make([]map[ID]*Event, 0, size)
		for i := 0; i < size; i++ {
			_, _, size, err := iprot.ReadMapBegin()
			if err != nil {
				return thrift.PrependError("error reading map begin: ", err)
			}
			elem21 := make(map[ID]*Event, size)
			for i := 0; i < size; i++ {
				var elem22 ID
				if v, err := iprot.ReadI64(); err != nil {
					return thrift.PrependError("error reading field 0: ", err)
				} else {
					temp := ID(v)
					elem22 = temp
				}
				elem23 := NewEvent()
				if err := elem23.Read(iprot); err != nil {
					return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", elem23), err)
				}
				(elem21)[elem22] = elem23
			}
			if err := iprot.ReadMapEnd(); err != nil {
				return thrift.PrependError("error reading map end: ", err)
			}
			req = append(req, elem21)
		}
		if err := iprot.ReadListEnd(); err != nil {
			return thrift.PrependError("error reading list end: ", err)
		}
		iprot.ReadMessageEnd()

		return method.Invoke([]interface{}{ctx, req}).Error()
	}
}
