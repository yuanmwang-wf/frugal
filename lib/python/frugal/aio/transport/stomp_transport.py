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

from aiostomp import AioStomp
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

    def __init__(self, stomp_client: AioStomp,
                 topic_prefix: str='',
                 max_message_size: int=32 * 1024 * 1024):
        self._stomp_client = stomp_client
        self._topic_prefix = topic_prefix
        self._max_message_size = max_message_size

    def get_transport(self) -> FPublisherTransport:
        """
        Get a new FStompPublisherTransport.
        """
        return FStompPublisherTransport(
            self._stomp_client, self._topic_prefix, self._max_message_size)


class FStompPublisherTransport(FPublisherTransport):
    """
    FStompPublisherTransport is used exclusively for pub/sub scopes.
    Publishers use it to publish to a topic. Messaging brokers that support
    stomp protocol can be used as the underlying bus.
    """

    def __init__(self, stomp_client: AioStomp,
                 topic_prefix: str='',
                 max_message_size: int=32 * 1024 * 1024):
        super().__init__(max_message_size)
        self._stomp_client = stomp_client
        self._topic_prefix = topic_prefix

    async def open(self):
        """
        No-op. Client connection should be established outside of frugal.
        """
        pass

    async def close(self):
        """
        No-op. Client close should be handled outside of frugal.
        :return:
        """
        pass

    def is_open(self) -> bool:
        """
        No-op.
        """
        return True

    async def publish(self, topic: str, data):
        """
        Publish a message to stomp broker on a given topic.

        Args:
            topic: string
            data: bytearray
        """
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
        print('Stomp publish destination: {}'.format(destination))
        self._stomp_client.send(
            destination,
            data,
            headers={'persistent': 'true',
                     'content-type': 'application/octet-stream'}
        )


class FStompSubscriberTransportFactory(FSubscriberTransportFactory):
    """
    FStompSubscriberTransportFactory is used to create
    FStompSubscriberTransports.
    """

    def __init__(self, stomp_client: AioStomp, consumer_prefix: str='',
                 use_queue: bool=False):
        self._stomp_client = stomp_client
        self._consumer_prefix = consumer_prefix
        self._use_queue = use_queue

    def get_transport(self) -> FSubscriberTransport:
        """
        Get a new FStompSubscriberTransport.
        """
        return FStompSubscriberTransport(
            self._stomp_client, self._consumer_prefix, self._use_queue)


class FStompSubscriberTransport(FSubscriberTransport):
    """
    FStompSubscriberTransport is used exclusively for pub/sub scopes.
    Subscribers use it to subscribe to a pub/sub topic. Messaging brokers that
    support stomp protocol can be used as the underlying bus.
    """

    def __init__(self, stomp_client: AioStomp, consumer_prefix: str,
                 use_queue: bool=True):
        self._stomp_client = stomp_client
        self._consumer_prefix = consumer_prefix
        self._use_queue = use_queue
        self._sub = None

    async def subscribe(self, topic: str, callback):
        """
        Subscribe to the given topic and register a callback to
        invoke when a message is received.

        Args:
            topic: str
            callback: func
        """
        if self.is_subscribed():
            raise TTransportException(
                TTransportExceptionType.ALREADY_OPEN,
                'stomp connection already subscribed to topic')

        if not topic:
            raise TTransportException(
                TTransportExceptionType.UNKNOWN,
                'stomp transport cannot subscribe to empty topic')

        async def msg_handler(frame, _):
            ret = callback(TMemoryBuffer(frame.body[4:]))
            if inspect.iscoroutine(ret):
                await ret
            # aiostomp acks message automatically in client-individual mode
            # as long as handler function returns non-falsy value
            return True

        destination = '/{type}/{consumer_prefix}{frugal_prefix}{topic}'.format(
            type='queue' if self._use_queue else 'topic',
            consumer_prefix=self._consumer_prefix,
            frugal_prefix=FRUGAL_PREFIX,
            topic=topic
        )
        print('Stomp subscribe destination: {}'.format(destination))
        self._sub = self._stomp_client.subscribe(
            destination, ack='client-individual', handler=msg_handler)

    async def unsubscribe(self):
        """
        Unsubscribe from the currently subscribed topic.
        """
        if not self.is_subscribed():
            return

        self._stomp_client.unsubscribe(self._sub)
        self._sub = None

    def is_subscribed(self) -> bool:
        """
        Check whether the client is subscribed or not.
        """
        return bool(self._sub)
