#
# Autogenerated by Frugal Compiler (1.20.0)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#



from threading import Lock

from frugal.middleware import Method
from frugal.processor import FBaseProcessor
from frugal.processor import FProcessorFunction
from thrift.Thrift import TApplicationException
from thrift.Thrift import TMessageType

import excepts
import validStructs
import ValidTypes
from valid.Blah import *
from valid.ttypes import *


class Iface(object):
    """
    This is a service docstring.
    """

    def ping(self, ctx):
        """
        Use this to ping the server.
        
        Args:
            ctx: FContext
        """
        pass

    def bleh(self, ctx, one, Two, custom_ints):
        """
        Use this to tell the server how you feel.
        
        Args:
            ctx: FContext
            one: Thing
            Two: Stuff
            custom_ints: list of int (signed 32 bits)
        """
        pass

    def getThing(self, ctx):
        """
        Args:
            ctx: FContext
        """
        pass

    def getMyInt(self, ctx):
        """
        Args:
            ctx: FContext
        """
        pass


class Client(Iface):

    def __init__(self, transport, protocol_factory, middleware=None):
        """
        Create a new Client with a transport and protocol factory.

        Args:
            transport: FSynchronousTransport
            protocol_factory: FProtocolFactory
            middleware: ServiceMiddleware or list of ServiceMiddleware
        """
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]
        self._transport = transport
        self._protocol_factory = protocol_factory
        self._oprot = protocol_factory.get_protocol(transport)
        self._iprot = protocol_factory.get_protocol(transport)
        self._write_lock = Lock()
        self._methods = {
            'ping': Method(self._ping, middleware),
            'bleh': Method(self._bleh, middleware),
            'getThing': Method(self._getThing, middleware),
            'getMyInt': Method(self._getMyInt, middleware),
        }

    def ping(self, ctx):
        """
        Use this to ping the server.
        
        Args:
            ctx: FContext
        """
        return self._methods['ping']([ctx])

    def _ping(self, ctx):
        self._send_ping(ctx)
        self._recv_ping(ctx)

    def _send_ping(self, ctx):
        oprot = self._oprot
        with self._write_lock:
            oprot.get_transport().set_timeout(ctx.get_timeout())
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('ping', TMessageType.CALL, 0)
            args = ping_args()
            args.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()

    def _recv_ping(self, ctx):
        self._iprot.read_response_headers(ctx)
        _, mtype, _ = self._iprot.readMessageBegin()
        if mtype == TMessageType.EXCEPTION:
            x = TApplicationException()
            x.read(self._iprot)
            self._iprot.readMessageEnd()
            raise x
        result = ping_result()
        result.read(self._iprot)
        self._iprot.readMessageEnd()
        return

    def bleh(self, ctx, one, Two, custom_ints):
        """
        Use this to tell the server how you feel.
        
        Args:
            ctx: FContext
            one: Thing
            Two: Stuff
            custom_ints: list of int (signed 32 bits)
        """
        return self._methods['bleh']([ctx, one, Two, custom_ints])

    def _bleh(self, ctx, one, Two, custom_ints):
        self._send_bleh(ctx, one, Two, custom_ints)
        return self._recv_bleh(ctx)

    def _send_bleh(self, ctx, one, Two, custom_ints):
        oprot = self._oprot
        with self._write_lock:
            oprot.get_transport().set_timeout(ctx.get_timeout())
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('bleh', TMessageType.CALL, 0)
            args = bleh_args()
            args.one = one
            args.Two = Two
            args.custom_ints = custom_ints
            args.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()

    def _recv_bleh(self, ctx):
        self._iprot.read_response_headers(ctx)
        _, mtype, _ = self._iprot.readMessageBegin()
        if mtype == TMessageType.EXCEPTION:
            x = TApplicationException()
            x.read(self._iprot)
            self._iprot.readMessageEnd()
            raise x
        result = bleh_result()
        result.read(self._iprot)
        self._iprot.readMessageEnd()
        if result.oops is not None:
            raise result.oops
        if result.err2 is not None:
            raise result.err2
        if result.success is not None:
            return result.success
        x = TApplicationException(TApplicationException.MISSING_RESULT, "bleh failed: unknown result")
        raise x

    def getThing(self, ctx):
        """
        Args:
            ctx: FContext
        """
        return self._methods['getThing']([ctx])

    def _getThing(self, ctx):
        self._send_getThing(ctx)
        return self._recv_getThing(ctx)

    def _send_getThing(self, ctx):
        oprot = self._oprot
        with self._write_lock:
            oprot.get_transport().set_timeout(ctx.get_timeout())
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('getThing', TMessageType.CALL, 0)
            args = getThing_args()
            args.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()

    def _recv_getThing(self, ctx):
        self._iprot.read_response_headers(ctx)
        _, mtype, _ = self._iprot.readMessageBegin()
        if mtype == TMessageType.EXCEPTION:
            x = TApplicationException()
            x.read(self._iprot)
            self._iprot.readMessageEnd()
            raise x
        result = getThing_result()
        result.read(self._iprot)
        self._iprot.readMessageEnd()
        if result.success is not None:
            return result.success
        x = TApplicationException(TApplicationException.MISSING_RESULT, "getThing failed: unknown result")
        raise x

    def getMyInt(self, ctx):
        """
        Args:
            ctx: FContext
        """
        return self._methods['getMyInt']([ctx])

    def _getMyInt(self, ctx):
        self._send_getMyInt(ctx)
        return self._recv_getMyInt(ctx)

    def _send_getMyInt(self, ctx):
        oprot = self._oprot
        with self._write_lock:
            oprot.get_transport().set_timeout(ctx.get_timeout())
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('getMyInt', TMessageType.CALL, 0)
            args = getMyInt_args()
            args.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()

    def _recv_getMyInt(self, ctx):
        self._iprot.read_response_headers(ctx)
        _, mtype, _ = self._iprot.readMessageBegin()
        if mtype == TMessageType.EXCEPTION:
            x = TApplicationException()
            x.read(self._iprot)
            self._iprot.readMessageEnd()
            raise x
        result = getMyInt_result()
        result.read(self._iprot)
        self._iprot.readMessageEnd()
        if result.success is not None:
            return result.success
        x = TApplicationException(TApplicationException.MISSING_RESULT, "getMyInt failed: unknown result")
        raise x

class Processor(FBaseProcessor):

    def __init__(self, handler, middleware=None):
        """
        Create a new Processor.

        Args:
            handler: Iface
        """
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]

        super(Processor, self).__init__()
        self.add_to_processor_map('ping', _ping(Method(handler.ping, middleware), self.get_write_lock()))
        self.add_to_processor_map('bleh', _bleh(Method(handler.bleh, middleware), self.get_write_lock()))
        self.add_to_processor_map('getThing', _getThing(Method(handler.getThing, middleware), self.get_write_lock()))
        self.add_to_processor_map('getMyInt', _getMyInt(Method(handler.getMyInt, middleware), self.get_write_lock()))


class _ping(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._lock = lock

    def process(self, ctx, iprot, oprot):
        args = ping_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = ping_result()
        try:
            self._handler([ctx])
        except Exception as e:
            with self._lock:
                _write_application_exception(ctx, oprot, TApplicationException.UNKNOWN, "ping", e.message)
            raise
        with self._lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('ping', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()


class _bleh(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._lock = lock

    def process(self, ctx, iprot, oprot):
        args = bleh_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = bleh_result()
        try:
            result.success = self._handler([ctx, args.one, args.Two, args.custom_ints])
        except InvalidOperation as oops:
            result.oops = oops
        except excepts.ttypes.InvalidData as err2:
            result.err2 = err2
        except Exception as e:
            with self._lock:
                _write_application_exception(ctx, oprot, TApplicationException.UNKNOWN, "bleh", e.message)
            raise
        with self._lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('bleh', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()


class _getThing(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._lock = lock

    def process(self, ctx, iprot, oprot):
        args = getThing_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = getThing_result()
        try:
            result.success = self._handler([ctx])
        except Exception as e:
            with self._lock:
                _write_application_exception(ctx, oprot, TApplicationException.UNKNOWN, "getThing", e.message)
            raise
        with self._lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('getThing', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()


class _getMyInt(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._lock = lock

    def process(self, ctx, iprot, oprot):
        args = getMyInt_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = getMyInt_result()
        try:
            result.success = self._handler([ctx])
        except Exception as e:
            with self._lock:
                _write_application_exception(ctx, oprot, TApplicationException.UNKNOWN, "getMyInt", e.message)
            raise
        with self._lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('getMyInt', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()


def _write_application_exception(ctx, oprot, typ, method, message):
    x = TApplicationException(type=typ, message=message)
    oprot.write_response_headers(ctx)
    oprot.writeMessageBegin(method, TMessageType.EXCEPTION, 0)
    x.write(oprot)
    oprot.writeMessageEnd()
    oprot.get_transport().flush()
