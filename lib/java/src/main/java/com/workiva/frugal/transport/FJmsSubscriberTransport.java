package com.workiva.frugal.transport;

import com.workiva.frugal.protocol.FAsyncCallback;
import org.apache.thrift.TException;
import org.apache.thrift.transport.TMemoryInputTransport;
import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.BytesMessage;
import javax.jms.Destination;
import javax.jms.JMSException;
import javax.jms.MessageConsumer;
import javax.jms.Session;
import java.util.Arrays;

/**
 * TODO.
 */
public class FJmsSubscriberTransport implements FSubscriberTransport {
    private static final Logger LOGGER = LoggerFactory.getLogger(FJmsSubscriberTransport.class);

    // TODO frugal prefixing of topics

//    private String topic;
    // TODO session or connection?
    private final Session session;
    private MessageConsumer consumer;
    // TODO useTopics flag?

    protected FJmsSubscriberTransport(Session session) {
        this.session = session;
    }

    /**
     * TODO.
     */
    public static class Factory implements FSubscriberTransportFactory {
        private final Session session;

        public Factory(Session session) {
            this.session = session;
        }

        @Override
        public FSubscriberTransport getTransport() {
            return new FJmsSubscriberTransport(session);
        }
    }

    @Override
    public boolean isSubscribed() {
        return consumer != null;
    }

    @Override
    public void subscribe(String topic, FAsyncCallback callback) throws TException {
        // TODO test
        if (isSubscribed()) {
            throw new TTransportException(TTransportException.ALREADY_OPEN, "jms client already subscribed");
        }

        if (topic == null || "".equals(topic)) {
            throw new TTransportException("subscribe subject cannot be empty");
        }

        try {
            Destination destination = session.createTopic(topic);
            consumer = session.createConsumer(destination);
            consumer.setMessageListener(message -> {
                // TODO this cast is janky af
                BytesMessage bytesMessage = (BytesMessage) message;
                byte[] payload;
                try {
                    payload = new byte[(int) bytesMessage.getBodyLength()];
                    bytesMessage.readBytes(payload);
                } catch (JMSException e) {
                    LOGGER.error("failed to get bytes from message", e);
                    return;
                }

                if (payload.length < 4) {
                    LOGGER.warn("discarding invalid scope message frame");
                    return;
                }
                try {
                    // TODO better way than copying?
                    callback.onMessage(
                            new TMemoryInputTransport(Arrays.copyOfRange(payload, 4, payload.length))
                    );
                } catch (TException ignored) {
                    // TODO is this right?
                }
            });
        } catch (JMSException e) {
            throw new TException(e);
        }
    }

    @Override
    public void unsubscribe() {
        if (!isSubscribed()) {
            LOGGER.debug("jms transport already unsubscribed, returning");
            return;
        }

        try {
            consumer.close();
        } catch (JMSException e) {
            LOGGER.error("failed to close jms consumer", e);
        } finally {
            consumer = null;
        }
    }
}
