package com.workiva.frugal.transport;

import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.BytesMessage;
import javax.jms.JMSException;
import javax.jms.MessageProducer;
import javax.jms.Session;
import javax.jms.Topic;

/**
 * TODO.
 */
public class FJmsPublisherTransport implements FPublisherTransport {
    private static final Logger LOGGER = LoggerFactory.getLogger(FJmsPublisherTransport.class);

    // TODO frugal prefixing of topics

    // TODO session or connection?
    private final Session session;
    private MessageProducer producer;
    // TODO useTopics flag?

    protected FJmsPublisherTransport(Session session) {
        this.session = session;
    }

    /**
     * TODO.
     */
    public static class Factory implements FPublisherTransportFactory {

        private final Session session;

        public Factory(Session session) {
            this.session = session;
        }

        public FPublisherTransport getTransport() {
            return new FJmsPublisherTransport(session);
        }
    }

    @Override
    public boolean isOpen() {
        return producer != null;
    }

    @Override
    public void open() throws TTransportException {
        // TODO test
        // TODO right defaults?
        if (isOpen()) {
            LOGGER.debug("jms transport already open, returning");
            return;
        }

        try {
            producer = session.createProducer(null);
        } catch (JMSException e) {
            throw new TTransportException(e);
        }
    }

    @Override
    public void close() {
        if (!isOpen()) {
            LOGGER.debug("jms transport already closed, returning");
            return;
        }
        // TODO preconditions
        // TODO test
        try {
            producer.close();
        } catch (JMSException e) {
            LOGGER.error("failed to close jms producer", e);
        } finally {
            producer = null;
        }
    }

    @Override
    public int getPublishSizeLimit() {
        // TODO check this number
        return Integer.MAX_VALUE;
    }

    @Override
    public void publish(String topic, byte[] payload) throws TTransportException {
        // TODO test
        if (!isOpen()) {
            throw new TTransportException(TTransportException.NOT_OPEN, "failed to publish, jms client not open");
        }

        if (topic == null || "".equals(topic)) {
            throw new TTransportException("publish topic cannot be empty");
        }

        // TODO size check?

        try {
            Topic destination = session.createTopic(topic);
            BytesMessage message = session.createBytesMessage();
            message.writeBytes(payload);
            producer.send(destination, message);
        } catch (JMSException e) {
            throw new TTransportException(e);
        }

    }
}
