package com.workiva.frugal.server;

import java.util.Map;

/**
 * Provides an interface with which to handle events from an FNatsServer.
 */
public interface FNatsServerEventHandler {
    /**
     * Called when a request is received, but before it is put onto a work
     * queue.
     *
     * @param ephemeralProperties
     */
    void onRequestReceived(Map<Object, Object> ephemeralProperties);

    /**
     * Called when a request is retrieved from a work queue to begin being
     * processed.
     *
     * @param ephemeralProperties
     */
    void onRequestStarted(Map<Object, Object> ephemeralProperties);

    /**
     * Called when a request is done being processed.
     *
     * @param ephemeralProperties
     */
    void onRequestEnded(Map<Object, Object> ephemeralProperties);
}
