package com.workiva.frugal.transport;

import com.workiva.frugal.exception.TTransportExceptionType;
import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.BytesMessage;
import javax.jms.Connection;
import javax.jms.DeliveryMode;
import javax.jms.JMSException;
import javax.jms.MessageProducer;
import javax.jms.Session;
import javax.jms.Topic;

import static com.workiva.frugal.transport.FNatsTransport.FRUGAL_PREFIX;

/**
 * FJmsPublisherTransport implements FPublisherTransport by utilizing a JMS
 * connection.
 */
public class FJmsPublisherTransport implements FPublisherTransport {
    private static final Logger LOGGER = LoggerFactory.getLogger(FJmsPublisherTransport.class);

    // TODO should we try to batch with sessions at some point?
    // TODO sessions aren't threadsafe, do we need to do more than sync on publish?
    private final Connection connection;
    private final String topicPrefix;
    private final boolean durablePublishes;
    Session session;
    MessageProducer producer;

    FJmsPublisherTransport(Connection connection, String topicPrefix, boolean durablePublishes) {
        this.connection = connection;
        this.topicPrefix = topicPrefix;
        this.durablePublishes = durablePublishes;
    }

    /**
     * An FPublisherTransportFacory implementation which creates
     * FPublisherTransports backed by a JMS connection.
     */
    public static class Factory implements FPublisherTransportFactory {

        private final Connection connection;
        private final String topicPrefix;
        private final boolean durablePublishes;

        // TODO should we make a builder for this?
        public Factory(Connection connection, String topicPrefix, boolean durablePublishes) {
            this.connection = connection;
            this.topicPrefix = topicPrefix;
            this.durablePublishes = durablePublishes;
        }

        public FPublisherTransport getTransport() {
            return new FJmsPublisherTransport(connection, topicPrefix, durablePublishes);
        }
    }

    @Override
    public boolean isOpen() {
        return session != null && producer != null;
    }

    @Override
    public void open() throws TTransportException {
        // TODO test
        if (isOpen()) {
            LOGGER.info("jms transport already open, returning");
            return;
        }

        try {
            session = connection.createSession(false, Session.CLIENT_ACKNOWLEDGE);
            producer = session.createProducer(null);
            if (!durablePublishes) {
                producer.setDeliveryMode(DeliveryMode.NON_PERSISTENT);
            }
        } catch (JMSException e) {
            throw new TTransportException(e);
        }
    }

    @Override
    public void close() {
        if (!isOpen()) {
            LOGGER.info("jms transport already closed, returning");
            return;
        }
        // TODO test
        try {
            producer.close();
            session.close();
        } catch (JMSException e) {
            LOGGER.error("failed to close jms producer", e);
        } finally {
            producer = null;
            session = null;
        }
    }

    @Override
    public int getPublishSizeLimit() {
        // TODO check this number
        return 32 * 1024 * 1024;
    }

    @Override
    public synchronized void publish(String topic, byte[] payload) throws TTransportException {
        // TODO test
        if (!isOpen()) {
            throw new TTransportException(TTransportException.NOT_OPEN, "failed to publish, jms client not open");
        }

        if (topic == null || "".equals(topic)) {
            throw new TTransportException("publish topic cannot be empty");
        }

        if (payload.length > getPublishSizeLimit()) {
            throw new TTransportException(TTransportExceptionType.REQUEST_TOO_LARGE,
                    String.format("message exceeds %d bytes, was %d bytes",
                            getPublishSizeLimit(), payload.length));
        }

        String formattedTopic = getFormattedTopic(topic);
        LOGGER.debug("publishing message to '%s'", formattedTopic);
        try {
            Topic destination = session.createTopic(formattedTopic);
            BytesMessage message = session.createBytesMessage();
            message.writeBytes(payload);
            // TODO should set this?
            // message.setJMSCorrelationID();
            producer.send(destination, message);
        } catch (JMSException e) {
            throw new TTransportException(e);
        }
        LOGGER.debug("published message to '%s'", formattedTopic);
    }

    private String getFormattedTopic(String subject) {
        return FRUGAL_PREFIX + topicPrefix + subject;
    }
}
