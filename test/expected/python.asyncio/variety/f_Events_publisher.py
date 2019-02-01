#
# Autogenerated by Frugal Compiler (2.27.0)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#



import inspect
import sys
import traceback

from thrift.Thrift import TApplicationException
from thrift.Thrift import TMessageType
from thrift.Thrift import TType
from frugal.exceptions import TApplicationExceptionType
from frugal.middleware import Method
from frugal.subscription import FSubscription
from frugal.transport import TMemoryOutputBuffer

from .ttypes import *




class EventsPublisher(object):
    """
    This docstring gets added to the generated code because it has
    the @ sign. Prefix specifies topic prefix tokens, which can be static or
    variable.
    """

    _DELIMITER = '.'

    def __init__(self, provider, middleware=None):
        """
        Create a new EventsPublisher.

        Args:
            provider: FScopeProvider
            middleware: ServiceMiddleware or list of ServiceMiddleware
        """

        middleware = middleware or []
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]
        middleware += provider.get_middleware()
        self._transport, self._protocol_factory = provider.new_publisher()
        self._methods = {
            'publish_EventCreated': Method(self._publish_EventCreated, middleware),
            'publish_SomeInt': Method(self._publish_SomeInt, middleware),
            'publish_SomeStr': Method(self._publish_SomeStr, middleware),
            'publish_SomeList': Method(self._publish_SomeList, middleware),
        }

    async def open(self):
        await self._transport.open()

    async def close(self):
        await self._transport.close()

    async def publish_EventCreated(self, ctx, user, req):
        """
        This is a docstring.
        
        Args:
            ctx: FContext
            user: string
            req: Event
        """
        await self._methods['publish_EventCreated']([ctx, user, req])

    async def _publish_EventCreated(self, ctx, user, req):
        ctx.set_request_header('_topic_user', user)
        op = 'EventCreated'
        prefix = 'foo.{}.'.format(user)
        topic = '{}Events{}{}'.format(prefix, self._DELIMITER, op)
        buffer = TMemoryOutputBuffer(self._transport.get_publish_size_limit())
        oprot = self._protocol_factory.get_protocol(buffer)
        oprot.write_request_headers(ctx)
        oprot.writeMessageBegin(op, TMessageType.CALL, 0)
        req.write(oprot)
        oprot.writeMessageEnd()
        await self._transport.publish(topic, buffer.getvalue())


    async def publish_SomeInt(self, ctx, user, req):
        """
        Args:
            ctx: FContext
            user: string
            req: i64
        """
        await self._methods['publish_SomeInt']([ctx, user, req])

    async def _publish_SomeInt(self, ctx, user, req):
        ctx.set_request_header('_topic_user', user)
        op = 'SomeInt'
        prefix = 'foo.{}.'.format(user)
        topic = '{}Events{}{}'.format(prefix, self._DELIMITER, op)
        buffer = TMemoryOutputBuffer(self._transport.get_publish_size_limit())
        oprot = self._protocol_factory.get_protocol(buffer)
        oprot.write_request_headers(ctx)
        oprot.writeMessageBegin(op, TMessageType.CALL, 0)
        oprot.writeI64(req)
        oprot.writeMessageEnd()
        await self._transport.publish(topic, buffer.getvalue())


    async def publish_SomeStr(self, ctx, user, req):
        """
        Args:
            ctx: FContext
            user: string
            req: string
        """
        await self._methods['publish_SomeStr']([ctx, user, req])

    async def _publish_SomeStr(self, ctx, user, req):
        ctx.set_request_header('_topic_user', user)
        op = 'SomeStr'
        prefix = 'foo.{}.'.format(user)
        topic = '{}Events{}{}'.format(prefix, self._DELIMITER, op)
        buffer = TMemoryOutputBuffer(self._transport.get_publish_size_limit())
        oprot = self._protocol_factory.get_protocol(buffer)
        oprot.write_request_headers(ctx)
        oprot.writeMessageBegin(op, TMessageType.CALL, 0)
        oprot.writeString(req)
        oprot.writeMessageEnd()
        await self._transport.publish(topic, buffer.getvalue())


    async def publish_SomeList(self, ctx, user, req):
        """
        Args:
            ctx: FContext
            user: string
            req: list
        """
        await self._methods['publish_SomeList']([ctx, user, req])

    async def _publish_SomeList(self, ctx, user, req):
        ctx.set_request_header('_topic_user', user)
        op = 'SomeList'
        prefix = 'foo.{}.'.format(user)
        topic = '{}Events{}{}'.format(prefix, self._DELIMITER, op)
        buffer = TMemoryOutputBuffer(self._transport.get_publish_size_limit())
        oprot = self._protocol_factory.get_protocol(buffer)
        oprot.write_request_headers(ctx)
        oprot.writeMessageBegin(op, TMessageType.CALL, 0)
        oprot.writeListBegin(TType.MAP, len(req))
        for elem59 in req:
            oprot.writeMapBegin(TType.I64, TType.STRUCT, len(elem59))
            for elem61, elem60 in elem59.items():
                oprot.writeI64(elem61)
                elem60.write(oprot)
            oprot.writeMapEnd()
        oprot.writeListEnd()
        oprot.writeMessageEnd()
        await self._transport.publish(topic, buffer.getvalue())

