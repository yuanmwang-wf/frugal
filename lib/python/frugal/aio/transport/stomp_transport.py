# Copyright 2017 Workiva
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import inspect
from uuid import uuid4

from stomp import Connection
from stomp import ConnectionListener
from thrift.transport.TTransport import TMemoryBuffer
from thrift.transport.TTransport import TTransportException

from frugal.exceptions import TTransportExceptionType
from frugal.transport import FPublisherTransport
from frugal.transport import FPublisherTransportFactory
from frugal.transport import FSubscriberTransport
from frugal.transport import FSubscriberTransportFactory


FRUGAL_PREFIX = 'frugal.'


class FStompPublisherTransportFactory(FPublisherTransportFactory):
    """
    FStompPublisherTransportFactory is used to create
    FStompPublisherTransports.
    """

    def __init__(self, stomp_conn: Connection, topic_prefix: str,
                 max_message_size: int = 32 * 1024 * 1024):
        self._stomp_conn = stomp_conn
        self._topic_prefix = topic_prefix
        self._max_message_size = max_message_size

    def get_transport(self) -> FPublisherTransport:
        """
        Get a new FStompPublisherTransport.
        """
        return FStompPublisherTransport(
            self._stomp_conn, self._topic_prefix, self._max_message_size)


class FStompPublisherTransport(FPublisherTransport):
    """
    FStompPublisherTransport is used exclusively for pub/sub scopes.
    Publishers use it to publish to a topic. Messaging brokers that support
    stomp protocol can be used as the underlying bus.
    """

    def __init__(self, stomp_conn: Connection, topic_prefix: str,
                 max_message_size: int):
        super().__init__(max_message_size)
        self._stomp_conn = stomp_conn
        self._topic_prefix = topic_prefix

    async def open(self):
        """
        Open the stomp publisher in preparation for publishing.
        """
        if not self._stomp_conn.is_connected():
            raise TTransportException(TTransportException.NOT_OPEN,
                                      'stomp connection is not connected')

    async def close(self):
        """
        Close the stomp publisher transport and disconnect from the stomp
        broker.
        :return:
        """
        if not self.is_open():
            return

        await self._stomp_conn.disconnect()

    def is_open(self) -> bool:
        """
        Check to see if the transport is open.
        """
        return self._stomp_conn.is_connected()

    async def publish(self, topic: str, data):
        """
        Publish a message to stomp broker on a given topic.

        Args:
            topic: string
            data: bytearray
        """
        if not self.is_open():
            raise TTransportException(
                type=TTransportExceptionType.NOT_OPEN,
                message='Transport is not connected')
        if self._check_publish_size(data):
            raise TTransportException(
                type=TTransportExceptionType.REQUEST_TOO_LARGE,
                message='Message exceeds max message size'
            )
        destination = '/topic/{topic_prefix}{frugal_prefix}{topic}'.format(
            topic_prefix=self._topic_prefix,
            frugal_prefix=FRUGAL_PREFIX,
            topic=topic
        )
        await self._stomp_conn.send(destination,
                                    data,
                                    content_type='application/octet-stream',
                                    headers={'persistent': 'true'})


class FStompSubscriberTransportFactory(FSubscriberTransportFactory):
    """
    FStompSubscriberTransportFactory is used to create
    FStompSubscriberTransports.
    """

    def __init__(self, stomp_conn: Connection, consumer_prefix: str,
                 use_queue: bool):
        self._stomp_conn = stomp_conn
        self._consumer_prefix = consumer_prefix
        self._use_queue = use_queue

    def get_transport(self) -> FSubscriberTransport:
        """
        Get a new FStompSubscriberTransport.
        """
        return FStompSubscriberTransport(
            self._stomp_conn, self._consumer_prefix, self._use_queue)


class FStompSubscriberTransport(FSubscriberTransport):
    """
    FStompSubscriberTransport is used exclusively for pub/sub scopes.
    Subscribers use it to subscribe to a pub/sub topic. Messaging brokers that
    support stomp protocol can be used as the underlying bus.
    """

    def __init__(self, stomp_conn: Connection, consumer_prefix: str,
                 use_queue: bool):
        self._stomp_conn = stomp_conn
        self._consumer_prefix = consumer_prefix
        self._use_queue = use_queue
        self._sub_id = None

    @property
    def _is_subscribed(self):
        return self._sub_id is not None

    async def subscribe(self, topic: str, callback):
        """
        Subscribe to the given topic and register a callback to
        invoke when a message is received.

        Args:
            topic: str
            callback: func
        """
        if not self._stomp_conn.is_connected():
            raise TTransportException(TTransportException.NOT_OPEN,
                                      'stomp connection is not connected')

        if self.is_subscribed():
            raise TTransportException(
                TTransportExceptionType.ALREADY_OPEN,
                'stomp connection already subscribed to topic')

        if not topic:
            raise TTransportException(
                TTransportExceptionType.UNKNOWN,
                'stomp transport cannot subscribe to empty topic')

        msg_listener = FStompConnectionListener(callback, self._stomp_conn)
        sub_id = str(uuid4())
        self._stomp_conn.set_listener('msg_listener', msg_listener)
        destination = '/{type}/{consumer_prefix}{frugal_prefix}{topic}'.format(
            type='queue' if self._use_queue else 'topic',
            consumer_prefix=self._consumer_prefix,
            frugal_prefix=FRUGAL_PREFIX,
            topic=topic
        )
        self._stomp_conn.subscribe(
            destination, sub_id, ack='client-individual')
        self._sub_id = sub_id

    async def unsubscribe(self):
        """
        Unsubscribe from the currently subscribed topic.
        """
        if not self.is_subscribed():
            return

        await self._stomp_conn.unsubscribe(self._sub_id)
        self._sub_id = None

    def is_subscribed(self) -> bool:
        """
        Check whether the client is subscribed or not.
        """
        return self._is_subscribed


class FStompConnectionListener(ConnectionListener):
    """
    FStompConnectionListener can be used to register callback functions to a
    connection when frames are received from the stomp broker.
    """

    def __init__(self, msg_callback, conn: Connection):
        self._msg_callback = msg_callback
        self._stomp_conn = conn

    async def on_message(self, headers: dict, body):
        """
        Called when a MESSAGE frame is received.
        Calls the registered callback function on the received body, then
        acknowledge the message if no exception is raised.

        Args:
            headers: dict
            body: frame's payload
        """
        ret = self._msg_callback(TMemoryBuffer(body.data[4:]))
        if inspect.iscoroutine(ret):
            ret = await ret

        msg_id = headers['message-id']
        sub_id = headers['subscription']
        self._stomp_conn.ack(msg_id, sub_id)
        return ret
