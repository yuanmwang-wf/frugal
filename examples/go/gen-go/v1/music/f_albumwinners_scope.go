// Autogenerated by Frugal Compiler (2.4.0-RC1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package music

import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/go"
)

const delimiter = "."

// Scopes are a Frugal extension to the IDL for declaring PubSub
// semantics. Subscribers to this scope will be notified if they win a contest.
// Scopes must have a prefix.
type AlbumWinnersPublisher interface {
	Open() error
	Close() error
	PublishContestStart(ctx frugal.FContext, req []*Album) error
	PublishTimeLeft(ctx frugal.FContext, req Minutes) error
	PublishWinner(ctx frugal.FContext, req *Album) error
}

type albumWinnersPublisher struct {
	transport       frugal.FPublisherTransport
	protocolFactory *frugal.FProtocolFactory
	methods         map[string]*frugal.Method
}

func NewAlbumWinnersPublisher(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) AlbumWinnersPublisher {
	transport, protocolFactory := provider.NewPublisher()
	methods := make(map[string]*frugal.Method)
	publisher := &albumWinnersPublisher{
		transport:       transport,
		protocolFactory: protocolFactory,
		methods:         methods,
	}
	middleware = append(middleware, provider.GetMiddleware()...)
	methods["publishContestStart"] = frugal.NewMethod(publisher, publisher.publishContestStart, "publishContestStart", middleware)
	methods["publishTimeLeft"] = frugal.NewMethod(publisher, publisher.publishTimeLeft, "publishTimeLeft", middleware)
	methods["publishWinner"] = frugal.NewMethod(publisher, publisher.publishWinner, "publishWinner", middleware)
	return publisher
}

func (p *albumWinnersPublisher) Open() error {
	return p.transport.Open()
}

func (p *albumWinnersPublisher) Close() error {
	return p.transport.Close()
}

func (p *albumWinnersPublisher) PublishContestStart(ctx frugal.FContext, req []*Album) error {
	ret := p.methods["publishContestStart"].Invoke([]interface{}{ctx, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersPublisher) publishContestStart(ctx frugal.FContext, req []*Album) error {
	op := "ContestStart"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(req)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range req {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
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

func (p *albumWinnersPublisher) PublishTimeLeft(ctx frugal.FContext, req Minutes) error {
	ret := p.methods["publishTimeLeft"].Invoke([]interface{}{ctx, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersPublisher) publishTimeLeft(ctx frugal.FContext, req Minutes) error {
	op := "TimeLeft"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteDouble(float64(req)); err != nil {
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

func (p *albumWinnersPublisher) PublishWinner(ctx frugal.FContext, req *Album) error {
	ret := p.methods["publishWinner"].Invoke([]interface{}{ctx, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersPublisher) publishWinner(ctx frugal.FContext, req *Album) error {
	op := "Winner"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
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

// Scopes are a Frugal extension to the IDL for declaring PubSub
// semantics. Subscribers to this scope will be notified if they win a contest.
// Scopes must have a prefix.
type AlbumWinnersSubscriber interface {
	SubscribeContestStart(handler func(frugal.FContext, []*Album)) (*frugal.FSubscription, error)
	SubscribeTimeLeft(handler func(frugal.FContext, Minutes)) (*frugal.FSubscription, error)
	SubscribeWinner(handler func(frugal.FContext, *Album)) (*frugal.FSubscription, error)
}

type albumWinnersSubscriber struct {
	provider   *frugal.FScopeProvider
	middleware []frugal.ServiceMiddleware
}

func NewAlbumWinnersSubscriber(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) AlbumWinnersSubscriber {
	middleware = append(middleware, provider.GetMiddleware()...)
	return &albumWinnersSubscriber{provider: provider, middleware: middleware}
}

func (l *albumWinnersSubscriber) SubscribeContestStart(handler func(frugal.FContext, []*Album)) (*frugal.FSubscription, error) {
	op := "ContestStart"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvContestStart(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersSubscriber) recvContestStart(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, []*Album)) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeContestStart", l.middleware)
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
		req := make([]*Album, 0, size)
		for i := 0; i < size; i++ {
			elem1 := NewAlbum()
			if err := elem1.Read(iprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", elem1), err)
			}
			req = append(req, elem1)
		}
		if err := iprot.ReadListEnd(); err != nil {
			return thrift.PrependError("error reading list end: ", err)
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, req})
		return nil
	}
}

func (l *albumWinnersSubscriber) SubscribeTimeLeft(handler func(frugal.FContext, Minutes)) (*frugal.FSubscription, error) {
	op := "TimeLeft"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvTimeLeft(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersSubscriber) recvTimeLeft(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, Minutes)) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeTimeLeft", l.middleware)
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
		var req Minutes
		if v, err := iprot.ReadDouble(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			temp := Minutes(v)
			req = temp
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, req})
		return nil
	}
}

func (l *albumWinnersSubscriber) SubscribeWinner(handler func(frugal.FContext, *Album)) (*frugal.FSubscription, error) {
	op := "Winner"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvWinner(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersSubscriber) recvWinner(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, *Album)) frugal.FAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeWinner", l.middleware)
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
		req := NewAlbum()
		if err := req.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", req), err)
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, req})
		return nil
	}
}

// Scopes are a Frugal extension to the IDL for declaring PubSub
// semantics. Subscribers to this scope will be notified if they win a contest.
// Scopes must have a prefix.
type AlbumWinnersDurablePublisher interface {
	Open() error
	Close() error
	PublishContestStart(ctx frugal.FContext, groupID *string, req []*Album) error
	PublishTimeLeft(ctx frugal.FContext, groupID *string, req Minutes) error
	PublishWinner(ctx frugal.FContext, groupID *string, req *Album) error
}

type albumWinnersDurablePublisher struct {
	transport       frugal.FDurablePublisherTransport
	protocolFactory *frugal.FProtocolFactory
	methods         map[string]*frugal.Method
}

func NewAlbumWinnersDurablePublisher(provider *frugal.FDurableScopeProvider, middleware ...frugal.ServiceMiddleware) AlbumWinnersDurablePublisher {
	transport, protocolFactory := provider.NewPublisher()
	methods := make(map[string]*frugal.Method)
	publisher := &albumWinnersDurablePublisher{
		transport:       transport,
		protocolFactory: protocolFactory,
		methods:         methods,
	}
	middleware = append(middleware, provider.GetMiddleware()...)
	methods["publishContestStart"] = frugal.NewMethod(publisher, publisher.publishContestStart, "publishContestStart", middleware)
	methods["publishTimeLeft"] = frugal.NewMethod(publisher, publisher.publishTimeLeft, "publishTimeLeft", middleware)
	methods["publishWinner"] = frugal.NewMethod(publisher, publisher.publishWinner, "publishWinner", middleware)
	return publisher
}

func (p *albumWinnersDurablePublisher) Open() error {
	return p.transport.Open()
}

func (p *albumWinnersDurablePublisher) Close() error {
	return p.transport.Close()
}

func (p *albumWinnersDurablePublisher) PublishContestStart(ctx frugal.FContext, groupID *string, req []*Album) error {
	ret := p.methods["publishContestStart"].Invoke([]interface{}{ctx, groupID, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersDurablePublisher) publishContestStart(ctx frugal.FContext, groupID *string, req []*Album) error {
	op := "ContestStart"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(req)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range req {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
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
	return p.transport.Publish(topic, groupID, buffer.Bytes())
}

func (p *albumWinnersDurablePublisher) PublishTimeLeft(ctx frugal.FContext, groupID *string, req Minutes) error {
	ret := p.methods["publishTimeLeft"].Invoke([]interface{}{ctx, groupID, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersDurablePublisher) publishTimeLeft(ctx frugal.FContext, groupID *string, req Minutes) error {
	op := "TimeLeft"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	buffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())
	oprot := p.protocolFactory.GetProtocol(buffer)
	if err := oprot.WriteRequestHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {
		return err
	}
	if err := oprot.WriteDouble(float64(req)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return p.transport.Publish(topic, groupID, buffer.Bytes())
}

func (p *albumWinnersDurablePublisher) PublishWinner(ctx frugal.FContext, groupID *string, req *Album) error {
	ret := p.methods["publishWinner"].Invoke([]interface{}{ctx, groupID, req})
	if ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (p *albumWinnersDurablePublisher) publishWinner(ctx frugal.FContext, groupID *string, req *Album) error {
	op := "Winner"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
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
	return p.transport.Publish(topic, groupID, buffer.Bytes())
}

// Scopes are a Frugal extension to the IDL for declaring PubSub
// semantics. Subscribers to this scope will be notified if they win a contest.
// Scopes must have a prefix.
type AlbumWinnersDurableSubscriber interface {
	SubscribeContestStart(handler func(frugal.FContext, *string, []*Album)) (*frugal.FSubscription, error)
	SubscribeTimeLeft(handler func(frugal.FContext, *string, Minutes)) (*frugal.FSubscription, error)
	SubscribeWinner(handler func(frugal.FContext, *string, *Album)) (*frugal.FSubscription, error)
}

type albumWinnersDurableSubscriber struct {
	provider   *frugal.FDurableScopeProvider
	middleware []frugal.ServiceMiddleware
}

func NewAlbumWinnersDurableSubscriber(provider *frugal.FDurableScopeProvider, middleware ...frugal.ServiceMiddleware) AlbumWinnersDurableSubscriber {
	middleware = append(middleware, provider.GetMiddleware()...)
	return &albumWinnersDurableSubscriber{provider: provider, middleware: middleware}
}

func (l *albumWinnersDurableSubscriber) SubscribeContestStart(handler func(frugal.FContext, *string, []*Album)) (*frugal.FSubscription, error) {
	op := "ContestStart"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvContestStart(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersDurableSubscriber) recvContestStart(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, *string, []*Album)) frugal.FDurableAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeContestStart", l.middleware)
	return func(transport thrift.TTransport, groupID *string) error {
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
		req := make([]*Album, 0, size)
		for i := 0; i < size; i++ {
			elem2 := NewAlbum()
			if err := elem2.Read(iprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", elem2), err)
			}
			req = append(req, elem2)
		}
		if err := iprot.ReadListEnd(); err != nil {
			return thrift.PrependError("error reading list end: ", err)
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, groupID, req})
		return nil
	}
}

func (l *albumWinnersDurableSubscriber) SubscribeTimeLeft(handler func(frugal.FContext, *string, Minutes)) (*frugal.FSubscription, error) {
	op := "TimeLeft"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvTimeLeft(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersDurableSubscriber) recvTimeLeft(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, *string, Minutes)) frugal.FDurableAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeTimeLeft", l.middleware)
	return func(transport thrift.TTransport, groupID *string) error {
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
		var req Minutes
		if v, err := iprot.ReadDouble(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			temp := Minutes(v)
			req = temp
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, groupID, req})
		return nil
	}
}

func (l *albumWinnersDurableSubscriber) SubscribeWinner(handler func(frugal.FContext, *string, *Album)) (*frugal.FSubscription, error) {
	op := "Winner"
	prefix := "v1.music."
	topic := fmt.Sprintf("%sAlbumWinners%s%s", prefix, delimiter, op)
	transport, protocolFactory := l.provider.NewSubscriber()
	cb := l.recvWinner(op, protocolFactory, handler)
	if err := transport.Subscribe(topic, cb); err != nil {
		return nil, err
	}

	sub := frugal.NewFSubscription(topic, transport)
	return sub, nil
}

func (l *albumWinnersDurableSubscriber) recvWinner(op string, pf *frugal.FProtocolFactory, handler func(frugal.FContext, *string, *Album)) frugal.FDurableAsyncCallback {
	method := frugal.NewMethod(l, handler, "SubscribeWinner", l.middleware)
	return func(transport thrift.TTransport, groupID *string) error {
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
		req := NewAlbum()
		if err := req.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", req), err)
		}
		iprot.ReadMessageEnd()

		method.Invoke([]interface{}{ctx, groupID, req})
		return nil
	}
}
