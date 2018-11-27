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

import asyncio
import mock

from thrift.transport.TTransport import TTransportException

from frugal.aio.transport import FStompConnectionListener
from frugal.aio.transport import FStompPublisherTransportFactory
from frugal.aio.transport import FStompSubscriberTransportFactory
from frugal.exceptions import TTransportExceptionType
from frugal.tests.aio import utils


class TestFStompPublisherTransport(utils.AsyncIOTestCase):
    def setUp(self, gen=None):
        super().setUp()
        self.mock_stomp_conn = mock.Mock()
        pub_factory = FStompPublisherTransportFactory(
            self.mock_stomp_conn, 'VirtualTopic.')
        self.pub_trans = pub_factory.get_transport()

    @utils.async_runner
    async def test_stomp_publisher_default_max_msg_size(self):
        self.assertEqual(
            self.pub_trans.get_publish_size_limit(), 32 * 1024 * 1024)

    @utils.async_runner
    async def test_stomp_publisher_open_stomp_not_connected(self):
        self.mock_stomp_conn.is_connected.return_value = False
        with self.assertRaises(TTransportException) as cm:
            await self.pub_trans.open()
        self.assertEqual(TTransportExceptionType.NOT_OPEN, cm.exception.type)

    @utils.async_runner
    async def test_stomp_publisher_close(self):
        self.mock_stomp_conn.is_connected.return_value = True
        future = asyncio.Future()
        future.set_result(None)
        disconnect_mock = mock.Mock()
        disconnect_mock.return_value = future
        self.mock_stomp_conn.disconnect = disconnect_mock

        await self.pub_trans.close()
        self.assertTrue(disconnect_mock.called)

    @utils.async_runner
    async def test_stomp_publisher_publish_successfully(self):
        self.mock_stomp_conn.is_connected.return_value = True
        data = bytearray([0, 0, 5, 2, 3, 4, 5, 6])
        future = asyncio.Future()
        future.set_result(None)
        self.mock_stomp_conn.send.return_value = future
        await self.pub_trans.publish('foo', data)

        self.mock_stomp_conn.send.assert_called_once_with(
            '/topic/VirtualTopic.frugal.foo',
            data,
            content_type='application/octet-stream',
            headers={'persistent': 'true'}
        )

    @utils.async_runner
    async def test_stomp_publisher_fails_publishing_conn_is_closed(self):
        self.mock_stomp_conn.is_connected.return_value = False
        with self.assertRaises(TTransportException) as cm:
            await self.pub_trans.publish('foo', bytearray([0, 0, 5, 2, 3]))
        self.assertEqual(TTransportExceptionType.NOT_OPEN, cm.exception.type)

    @utils.async_runner
    async def test_stomp_publisher_fails_publishing_payload_too_large(self):
        mock_check_publish_size = mock.Mock(return_value=True)
        self.pub_trans._check_publish_size = mock_check_publish_size
        with self.assertRaises(TTransportException) as cm:
            await self.pub_trans.publish('foo', bytearray([0, 0, 5, 2, 3]))
        self.assertEqual(
            TTransportExceptionType.REQUEST_TOO_LARGE, cm.exception.type)


class TestFStompSubscriberTransport(utils.AsyncIOTestCase):
    def setUp(self, gen=None):
        super().setUp()
        self.mock_stomp_conn = mock.Mock()
        sub_factory = FStompSubscriberTransportFactory(
            stomp_conn=self.mock_stomp_conn,
            consumer_prefix='Consumer.foo.',
            use_queue=True)
        self.sub_trans = sub_factory.get_transport()

    @utils.async_runner
    async def test_stomp_subscriber_stomp_not_connected(self):
        self.mock_stomp_conn.is_connected.return_value = False
        with self.assertRaises(TTransportException) as cm:
            await self.sub_trans.subscribe('bar', None)
        self.assertEqual(TTransportException.NOT_OPEN, cm.exception.type)

    @utils.async_runner
    async def test_stomp_subscriber_already_subscribed(self):
        self.sub_trans._sub_id = 'sub_id'
        with self.assertRaises(TTransportException) as cm:
            await self.sub_trans.subscribe('bar', None)
        self.assertEqual(TTransportException.ALREADY_OPEN, cm.exception.type)

    @utils.async_runner
    async def test_stomp_subscriber_empty_topic(self):
        with self.assertRaises(TTransportException) as cm:
            await self.sub_trans.subscribe('', None)
        self.assertEqual(TTransportException.UNKNOWN, cm.exception.type)

    @utils.async_runner
    async def test_stomp_subscriber_subscribe(self):
        mock_cb = mock.Mock()
        await self.sub_trans.subscribe('topic', mock_cb)
        self.mock_stomp_conn.subscribe.assert_called_once_with(
            '/queue/Consumer.foo.frugal.topic',
            self.sub_trans._sub_id,
            ack='client-individual'
        )
        self.assertTrue(self.sub_trans.is_subscribed())

    @utils.async_runner
    async def test_stomp_subscriber_unsubscribe(self):
        self.sub_trans._sub_id = 'sub_id'
        future = asyncio.Future()
        future.set_result(None)
        self.mock_stomp_conn.unsubscribe.return_value = future
        await self.sub_trans.unsubscribe()

        self.mock_stomp_conn.unsubscribe.assert_called_once_with('sub_id')
        self.assertFalse(self.sub_trans.is_subscribed())

    @utils.async_runner
    async def test_stomp_subscriber_unsubscribe_already_unsubscribed(self):
        await self.sub_trans.unsubscribe()
        self.assertFalse(self.mock_stomp_conn.called)


class TestFStompConnectionListener(utils.AsyncIOTestCase):
    def setUp(self, gen=None):
        super().setUp()
        self.mock_stomp_conn = mock.Mock()
        self.mock_cb = mock.Mock()
        self.listener = FStompConnectionListener(
            self.mock_cb, self.mock_stomp_conn)

    @utils.async_runner
    async def test_connection_listener(self):
        future = asyncio.Future()
        future.set_result(None)
        mock_ack = mock.Mock(return_value=future)
        self.mock_stomp_conn.ack = mock_ack
        ret = await self.listener.on_message(
            headers={'message-id': '123', 'subscription': 'sub_id'},
            body=mock.Mock(data=bytearray([0, 1, 2, 3, 4, 5])))
        self.assertTrue(self.mock_cb.called)
        self.assertEqual(ret, self.mock_cb.return_value)
        self.mock_stomp_conn.ack.assert_called_once_with('123', 'sub_id')
