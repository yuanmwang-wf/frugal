package com.workiva.frugal.transport;

import com.workiva.frugal.protocol.FAsyncCallback;
import org.apache.thrift.TException;
import org.apache.thrift.transport.TMemoryInputTransport;
import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.BytesMessage;
import javax.jms.Connection;
import javax.jms.Destination;
import javax.jms.JMSException;
import javax.jms.MessageConsumer;
import javax.jms.Session;

import static com.workiva.frugal.transport.FNatsTransport.FRUGAL_PREFIX;

/**
 * FJmsSubscriberTransport implements FSubscriberTransport by utilizing a JSM
 * connection.
 */
public class FJmsSubscriberTransport implements FSubscriberTransport {
    private static final Logger LOGGER = LoggerFactory.getLogger(FJmsSubscriberTransport.class);

    // TODO sessions aren't threadsafe, do we need to do more than sync on sub/unsub?
    private final Connection connection;
    private final String topicPrefix;
    private final boolean useQueues;
    Session session;
    MessageConsumer consumer;

    FJmsSubscriberTransport(Connection connection, String topicPrefix, boolean useQueues) {
        this.connection = connection;
        this.topicPrefix = topicPrefix;
        this.useQueues = useQueues;
    }

    /**
     * An FSubscriberTransportFacory implementation which creates
     * FSubscriberTransports backed by a JMS connection.
     */
    public static class Factory implements FSubscriberTransportFactory {
        private final Connection connection;
        private final String topicPrefix;
        private final boolean useQueues;

        // TODO should we make a builder for this?
        public Factory(Connection connection, String topicPrefix, boolean useQueues) {
            this.connection = connection;
            this.topicPrefix = topicPrefix;
            this.useQueues = useQueues;
        }

        @Override
        public FSubscriberTransport getTransport() {
            return new FJmsSubscriberTransport(connection, topicPrefix, useQueues);
        }
    }

    @Override
    public synchronized boolean isSubscribed() {
        return session != null && consumer != null;
    }

    /**
     * @inheritDoc
     *
     * If an exception is raised by the provided callback, the message will
     * not be acked with the broker. This behaviour allows the message to be
     * redelivered and processing to be attempted again. If an exception is
     * not raised by the provided callback, the message will be acked. This is
     * used if processing succeeded, or if it's apparent processing will never
     * succeed, as the message won't continue to be redelivered.
     */
    @Override
    public synchronized void subscribe(String topic, FAsyncCallback callback) throws TException {
        // TODO test
        if (isSubscribed()) {
            throw new TTransportException(TTransportException.ALREADY_OPEN, "jms client already subscribed");
        }

        if (topic == null || "".equals(topic)) {
            throw new TTransportException("subscribe subject cannot be empty");
        }
        String formattedTopic = getFormattedTopic(topic);

        try {
            // TODO I think these are the defaults we want
            session = connection.createSession(false, Session.CLIENT_ACKNOWLEDGE);

            Destination destination;
            if (useQueues) {
                destination = session.createQueue(formattedTopic);
            } else {
                destination = session.createTopic(formattedTopic);
            }

            consumer = session.createConsumer(destination);
            consumer.setMessageListener(message -> {
                LOGGER.debug("received a message on topic '%s'", formattedTopic);

                // TODO do we need to handle other message types?
                byte[] payload;
                if (message instanceof BytesMessage) {
                    BytesMessage bytesMessage = (BytesMessage) message;
                    try {
                        payload = new byte[(int) bytesMessage.getBodyLength()];
                        bytesMessage.readBytes(payload);
                    } catch (JMSException e) {
                        LOGGER.error("failed to get bytes from message", e);
                        return;
                    }
                } else {
                    LOGGER.error("unhandled message type '%s'", message.getClass().getName());
                    return;
                }

                if (payload.length < 4) {
                    LOGGER.warn("discarding invalid message frame, length less than four");
                    return;
                }
                try {
                    callback.onMessage(
                            new TMemoryInputTransport(payload, 4, payload.length - 4)
                    );
                } catch (TException e) {
                    LOGGER.error("error executing user provided callback", e);
                    return;
                }

                try {
                    message.acknowledge();
                } catch (JMSException e) {
                    LOGGER.error("unable to ack message", e);
                }
                LOGGER.debug("finished processing message from topic '%s'", formattedTopic);
            });
        } catch (JMSException e) {
            throw new TException(e);
        }
    }

    @Override
    public synchronized void unsubscribe() {
        if (!isSubscribed()) {
            LOGGER.info("jms transport already unsubscribed, returning");
            return;
        }

        try (Session closeSession = session) {
            consumer.close();
        } catch (JMSException e) {
            LOGGER.error("failed to close jms consumer", e);
        } finally {
            consumer = null;
            session = null;
        }
    }

    private String getFormattedTopic(String subject) {
        return topicPrefix + FRUGAL_PREFIX + subject;
    }
}
