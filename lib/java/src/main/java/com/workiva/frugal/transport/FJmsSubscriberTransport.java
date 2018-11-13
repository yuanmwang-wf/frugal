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
import java.util.Arrays;

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
    Session session;
    MessageConsumer consumer;

    FJmsSubscriberTransport(Connection connection, String topicPrefix) {
        this.connection = connection;
        this.topicPrefix = topicPrefix;
    }

    /**
     * An FSubscriberTransportFacory implementation which creates
     * FSubscriberTransports backed by a JMS connection.
     */
    public static class Factory implements FSubscriberTransportFactory {
        private final Connection connection;

        public Factory(Connection connection) {
            this.connection = connection;
        }

        @Override
        public FSubscriberTransport getTransport() {
            return new FJmsSubscriberTransport(connection, "");
        }
    }

    @Override
    public boolean isSubscribed() {
        return session != null && consumer != null;
    }

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
            Destination destination = session.createTopic(formattedTopic);
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
                    // TODO better way than copying? what we do for other transports
                    callback.onMessage(
                            new TMemoryInputTransport(Arrays.copyOfRange(payload, 4, payload.length))
                    );
                } catch (TException ignored) {
                    // TODO is this right?
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

        try {
            consumer.close();
            session.close();
        } catch (JMSException e) {
            LOGGER.error("failed to close jms consumer", e);
        } finally {
            consumer = null;
            session = null;
        }
    }

    private String getFormattedTopic(String subject) {
        return FRUGAL_PREFIX + topicPrefix + subject;
    }
}
