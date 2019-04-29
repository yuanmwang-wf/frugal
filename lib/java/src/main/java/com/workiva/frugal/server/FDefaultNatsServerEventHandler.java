package com.workiva.frugal.server;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

/**
 * A default event handler for an FNatsServer.
 */
public class FDefaultNatsServerEventHandler implements FNatsServerEventHandler {
    private static final Logger LOGGER = LoggerFactory.getLogger(FDefaultNatsServerEventHandler.class);

    private long highWatermark;

    public FDefaultNatsServerEventHandler(long highWatermark) {
        this.highWatermark = highWatermark;
    }

    @Override
    public void onRequestReceived(Map<Object, Object> ephemeralProperties) {
        long now = System.currentTimeMillis();
        ephemeralProperties.put("_request_received_millis", now);
    }

    @Override
    public void onRequestStarted(Map<Object, Object> ephemeralProperties) {
        if (ephemeralProperties.get("_request_received_millis") != null) {
            long started = (long) ephemeralProperties.get("_request_received_millis");
            long duration = System.currentTimeMillis() - started;
            if (duration > highWatermark) {
                LOGGER.warn(String.format(
                        "request spent %d ms in the transport buffer, your consumer might be backed up", duration));
            }
        }
    }

    @Override
    public void onRequestEnded(Map<Object, Object> ephemeralProperties) {

    }
}
